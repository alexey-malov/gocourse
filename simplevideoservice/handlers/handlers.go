package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

import log "github.com/sirupsen/logrus"

type handlerBase struct {
	uploader usecases.Uploader
	videos   repository.Videos
}

type handler struct {
	router *mux.Router
	handlerBase
}

func MakeHandler(uploader usecases.Uploader, videos repository.Videos) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	h := &handler{r,
		handlerBase{
			uploader,
			videos}}

	s.HandleFunc("/list", h.list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", h.video).Methods(http.MethodGet)
	s.HandleFunc("/video", h.upload).Methods(http.MethodPost)
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
