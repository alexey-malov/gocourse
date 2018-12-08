package main

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/handlers"
	"net/http"
	"os"
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

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")

	router := handlers.Router()
	log.Fatal(http.ListenAndServe(":8000", router))

}
