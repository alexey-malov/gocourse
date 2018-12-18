package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/gorilla/mux"
	"net/http"
)

import log "github.com/sirupsen/logrus"

func makeHandlerFunc(vr repository.VideoRepository, handler func(vr repository.VideoRepository, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(vr, w, r)
	}
}

func Router(vr repository.VideoRepository) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/list", makeHandlerFunc(vr, list)).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", makeHandlerFunc(vr, video)).Methods(http.MethodGet)
	s.HandleFunc("/video", makeHandlerFunc(vr, uploadVideo)).Methods(http.MethodPost)
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
