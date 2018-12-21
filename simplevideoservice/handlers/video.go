package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/gorilla/mux"

	"net/http"
)

type videoContent struct {
	videoListItem
	Url string `json:"url"`
}

func makeVideoContent(v model.VideoItem) videoContent {
	return videoContent{
		videoListItem{
			v.Id(),
			v.Name(),
			v.Duration(),
			v.ScreenShotUrl(),
		},
		v.VideoUrl(),
	}
}

func video(vr repository.Videos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]

	v, err := vr.Find(id)
	if err != nil {
		return
	}

	if v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(makeVideoContent(*v))
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		return
	}
}
