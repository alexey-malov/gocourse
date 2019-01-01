package main

import (
	"context"
	"github.com/alexey-malov/gocourse/simplevideoservice/app"
	"github.com/alexey-malov/gocourse/simplevideoservice/handlers"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

import log "github.com/sirupsen/logrus"

func main() {
	if _, err := app.SetupLogger("server.log"); err != nil {
		log.Fatal("Failed to create log")
	}

	persister, err := app.MakeVideoPersister()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := persister.Close(); err != nil {
			log.Error(err)
		}
	}()

	stg, err := app.MakeStorage()
	if err != nil {
		log.Fatal(err)
	}

	finder := usecases.MakeFinder(persister.Videos())
	uploader := usecases.MakeUploader(persister.Videos(), stg)
	lister := usecases.MakeVideoLister(persister.Videos())

	killSignalChan := getKillSignalChan()
	uc := handlers.MakeUseCases(finder, uploader, lister)
	srv := startServer(":8000", uc)

	waitForKillSignal(killSignalChan)
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}

func startServer(serverUrl string, useCases handlers.UseCases) *http.Server {
	router := handlers.MakeHandler(useCases)
	srv := &http.Server{Addr: serverUrl, Handler: router}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan /*, os.Interrupt, syscall.SIGTERM*/)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}
