package app

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger(name string) (*os.File, error) {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	return file, err
}
