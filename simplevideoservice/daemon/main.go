package main

import (
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/daemon/task"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func taskGenerator() task.Task {
	return func() {
		d := rand.Intn(3)
		fmt.Printf("Task: Sleeping %d seconds\n", d)
		time.Sleep(time.Duration(d) * time.Second)
	}
}

func main() {
	stopChan := make(chan struct{})

	wg := task.RunWorkerPool(stopChan, task.MakeDefaultTaskProvider(taskGenerator))
	defer wg.Wait()

	waitForKillSignal(getKillSignalChan())

	stopChan <- struct{}{}
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
