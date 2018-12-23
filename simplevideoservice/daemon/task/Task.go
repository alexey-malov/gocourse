package task

import (
	"log"
	"sync"
	"time"
)

type Task func()

type TaskProvider func(stopChan chan struct{}) <-chan Task

type TaskGenerator func() Task

func DefaultTaskProvider(stopChan chan struct{}, generator TaskGenerator) <-chan Task {
	tasksChan := make(chan Task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}
			if task := generator(); task != nil {
				log.Printf("got the task %v\n", task)
				select {
				case <-stopChan:
					close(tasksChan)
					return
				case tasksChan <- task:
				}
			} else {
				log.Println("no task for processing, start waiting")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return tasksChan
}

func MakeDefaultTaskProvider(generator TaskGenerator) TaskProvider {
	return func(stopChan chan struct{}) <-chan Task {
		return DefaultTaskProvider(stopChan, generator)
	}
}

func RunWorkerPool(stopChan chan struct{}, provider TaskProvider) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := runTaskProvider(stopChan, provider)
	for i := 0; i < 3; i++ {
		go func(i int) {
			wg.Add(1)
			worker(tasksChan, i)
			wg.Done()
		}(i)
	}
	return &wg
}

func runTaskProvider(stopChan chan struct{}, provider TaskProvider) <-chan Task {
	resultChan := make(chan Task)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := provider(stopTaskProviderChan)
	onStop := func() {
		stopTaskProviderChan <- struct{}{}
		close(resultChan)
	}
	go func() {
		for {
			select {
			case <-stopChan:
				onStop()
				return
			case task := <-taskProviderChan:
				select {
				case <-stopChan:
					onStop()
					return
				case resultChan <- task:
				}
			}
		}
	}()
	return resultChan
}

func worker(tasksChan <-chan Task, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start handle task %v on worker %v\n", task, name)
		task()
		log.Printf("end handle task %v on worker %v\n", task, name)
	}
	log.Printf("stop worker %v\n", name)
}
