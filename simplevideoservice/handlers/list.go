package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"net/http"
)

type videoListItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func makeVideoListItem(v model.VideoItem) videoListItem {
	return videoListItem{
		v.Id(),
		v.Name(),
		v.Duration(),
		v.ScreenShotUrl(),
	}
}

func list(vr repository.Videos, w http.ResponseWriter, _ *http.Request) {
	var videos []videoListItem

	err := vr.Enumerate(func(v model.VideoItem) bool {
		videos = append(videos, makeVideoListItem(v))
		return true
	})
	if err != nil {
		return
	}

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
