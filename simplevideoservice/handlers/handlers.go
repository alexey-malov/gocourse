package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

import log "github.com/sirupsen/logrus"

func makeHandlerFunc(db *sql.DB, handler func(db *sql.DB, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}
}

func Router(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/list", makeHandlerFunc(db, list)).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", makeHandlerFunc(db, video)).Methods(http.MethodGet)
	s.HandleFunc("/video", makeHandlerFunc(db, uploadVideo)).Methods(http.MethodPost)
	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}
