package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

import log "github.com/sirupsen/logrus"

type MyRouter struct {
	router   *mux.Router
	uploader usecases.Uploader
	videos   repository.Videos
}

func makeHandlerFunc(vr repository.Videos, handler func(vr repository.Videos, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(vr, w, r)
	}
}

func MakeRouter(uploader usecases.Uploader, videos repository.Videos) *MyRouter {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/list", makeHandlerFunc(videos, list)).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", makeHandlerFunc(videos, video)).Methods(http.MethodGet)

	myRouter := &MyRouter{r, uploader, videos}
	s.HandleFunc("/video", myRouter.uploadVideo).Methods(http.MethodPost)
	return myRouter
}

func (mr *MyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method":     r.Method,
		"url":        r.URL,
		"remoteAddr": r.RemoteAddr,
		"userAgent":  r.UserAgent(),
	}).Info("got a new request")
	mr.router.ServeHTTP(w, r)

}
