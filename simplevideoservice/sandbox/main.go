package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Task struct {
	id int
}

func GenerateTask() *Task {
	if rand.Intn(1000) < 500 {
		return nil
	}

	log.Println("Generating task")
	return &Task{rand.Intn(10)}
}

func TaskProvider(stopChan chan struct{}) <-chan *Task {
	tasksChan := make(chan *Task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}
			if task := GenerateTask(); task != nil {
				log.Printf("got the task %v\n", task)
				tasksChan <- task
			} else {
				log.Println("no task for processing, start waiting")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return tasksChan
}

func RunTaskProvider(stopChan chan struct{}) <-chan *Task {
	resultChan := make(chan *Task)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := TaskProvider(stopTaskProviderChan)
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

func Worker(tasksChan <-chan *Task, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start handle task %v on worker %v\n", task, name)
		time.Sleep(time.Duration(task.id) * time.Second)
		log.Printf("end handle task %v on worker %v\n", task, name)
	}
	log.Printf("stop worker %v\n", name)
}

func RunWorkerPool(stopChan chan struct{}) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := RunTaskProvider(stopChan)
	for i := 0; i < 3; i++ {
		go func(i int) {
			wg.Add(1)
			Worker(tasksChan, i)
			wg.Done()
		}(i)
	}
	return &wg
}

func main() {
	rand.Seed(time.Now().Unix())
	stopChan := make(chan struct{})

	killChan := getKillSignalChan()
	wg := RunWorkerPool(stopChan)

	waitForKillSignal(killChan)
	stopChan <- struct{}{}
	wg.Wait()
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Println("got SIGINT...")
	case syscall.SIGTERM:
		log.Println("got SIGTERM...")
	}
}
