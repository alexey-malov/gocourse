package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

import log "github.com/sirupsen/logrus"

type UseCases struct {
	uploader usecases.Uploader
	videos   repository.Videos
	finder   usecases.VideoFinder
}

func MakeUseCases(finder usecases.VideoFinder, uploader usecases.Uploader, videos repository.Videos) UseCases {
	return UseCases{uploader, videos, finder}
}

type handler struct {
	router   *mux.Router
	useCases UseCases
}

func MakeHandler(useCases UseCases) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	h := &handler{r,
		useCases}

	s.HandleFunc("/list", h.useCases.list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", h.useCases.video).Methods(http.MethodGet)
	s.HandleFunc("/video", h.useCases.upload).Methods(http.MethodPost)
	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method":     r.Method,
		"url":        r.URL,
		"remoteAddr": r.RemoteAddr,
		"userAgent":  r.UserAgent(),
	}).Info("got a new request")
	h.router.ServeHTTP(w, r)

}
