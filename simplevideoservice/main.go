package main

import (
	"context"
	"database/sql"
	"github.com/alexey-malov/gocourse/simplevideoservice/handlers"
	_ "github.com/go-sql-driver/mysql"
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
	if _, err := setupLogger(); err != nil {
		log.Fatal("Failed to create log")
	}

	db, err := sql.Open("mysql", "root:Jcbdsl7625@/simplevideoservice")
	if err != nil {
		log.Fatal("Failed to open DB")
	}
	defer db.Close()

	if db.Ping() != nil {
		log.Fatal("Failed to ping db")
	}

	killSignalChan := getKillSignalChan()
	srv := startServer(":8000", db)

	waitForKillSignal(killSignalChan)
	if err := srv.Shutdown(context.Background()); err != nil {
	}
}

func startServer(serverUrl string, db *sql.DB) *http.Server {
	router := handlers.Router(db)
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
