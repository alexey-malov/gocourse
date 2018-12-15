package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"

	"net/http"
)

type videoContent struct {
	videoListItem
	Url string `json:"url"`
}

func makeVideoContent(v videoItem) videoContent {
	return videoContent{
		videoListItem{
			v.id,
			v.name,
			v.duration,
			v.screenShotUrl(),
		},
		v.videoUrl(),
	}
}

func video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]

	v := findVideo(id)
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
