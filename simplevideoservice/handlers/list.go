package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
)

type videoListItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func makeVideoListItem(v domain.Video) videoListItem {
	return videoListItem{
		v.Id(),
		v.Name(),
		v.Duration(),
		v.ThumbnailURL(),
	}
}

func list(lister usecases.VideoLister, w http.ResponseWriter, _ *http.Request) {
	var videos []videoListItem

	err := lister.List(func(v *domain.Video) (bool, error) {
		videos = append(videos, makeVideoListItem(*v))
		return true, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(videos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		logrus.Error(err)
	}
}
