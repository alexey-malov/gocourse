package handlers

import (
	"encoding/json"
	"net/http"
)

type videoListItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func makeVideoListItem(v videoItem) videoListItem {
	return videoListItem{
		v.id,
		v.name,
		v.duration,
		v.screenShotUrl(),
	}
}

func list(w http.ResponseWriter, _ *http.Request) {

	var videos []videoListItem
	enumVideos(func(v videoItem) bool {
		videos = append(videos, makeVideoListItem(v))
		return true
	})

	b, err := json.Marshal(videos)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err == nil {
		return
	}
}
