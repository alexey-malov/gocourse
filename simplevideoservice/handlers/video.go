package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"net/http"
)

type videoContent struct {
	videoListItem
	Url string `json:"url"`
}

func makeVideoContent(v domain.Video) videoContent {
	return videoContent{
		videoListItem{
			v.Id(),
			v.Name(),
			v.Duration(),
			v.ThumbnailURL(),
		},
		v.VideoUrl(),
	}
}

func video(finder usecases.VideoFinder, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["ID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v, err := finder.Find(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error(err)
		return
	}

	if v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(makeVideoContent(*v))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		logrus.Error(err)
		return
	}
}
