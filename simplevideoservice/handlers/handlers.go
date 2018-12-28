package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

import log "github.com/sirupsen/logrus"

type UseCases struct {
	uploader usecases.Uploader
	finder   usecases.VideoFinder
	lister   usecases.VideoLister
}

func MakeUseCases(finder usecases.VideoFinder, uploader usecases.Uploader, lister usecases.VideoLister) UseCases {
	return UseCases{uploader, finder, lister}
}

type handler struct {
	router   *mux.Router
	useCases UseCases
}

func MakeHandler(uc UseCases) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	h := &handler{r,
		uc}

	s.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		list(uc.lister, w, r)
	}).Methods(http.MethodGet)

	s.HandleFunc("/video/{ID}", func(w http.ResponseWriter, r *http.Request) {
		video(uc.finder, w, r)
	}).Methods(http.MethodGet)

	s.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		upload(uc.uploader, w, r)
	}).Methods(http.MethodPost)

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
