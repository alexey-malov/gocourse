package main

import (
	"context"
	"github.com/alexey-malov/gocourse/simplevideoservice/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

import log "github.com/sirupsen/logrus"

func setupLogger() (*os.File, error) {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	return file, err
}

func main() {
	if file, _ := setupLogger(); file != nil {
		defer func() {
			_ = file.Close()
		}()
	} else {
		log.Fatal("Failed to create log")
	}

	killSignalChan := getKillSignalChan()
	srv := startServer(":8000")

	waitForKillSignal(killSignalChan)
	_ = srv.Shutdown(context.Background())
}

func startServer(serverUrl string) *http.Server {
	router := handlers.Router()
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
