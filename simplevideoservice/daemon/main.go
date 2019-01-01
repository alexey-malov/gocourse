package main

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/app"
	"github.com/alexey-malov/gocourse/simplevideoservice/daemon/processor"
	"github.com/alexey-malov/gocourse/simplevideoservice/daemon/task"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if logFile, err := app.SetupLogger("daemon.log"); err != nil {
		log.Fatal("Failed to create log")
	} else {
		defer func() {
			_ = logFile.Close()
		}()
	}

	persister, err := app.MakeVideoPersister()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := persister.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	stg, err := app.MakeStorage()
	if err != nil {
		logrus.Fatal(err)
	}

	stopChan := make(chan struct{})

	videoProcessor := processor.MakeVideoProcessor(persister.Videos(), stg)

	wg := task.RunWorkerPool(stopChan, task.MakeDefaultTaskProvider(videoProcessor))
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
