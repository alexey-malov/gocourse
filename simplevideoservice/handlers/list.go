package handlers

import (
	"encoding/json"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

func list(lister usecases.VideoLister, w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil || offset < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 || limit > 50 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	videos := make([]videoListItem, 0)
	err = lister.List(uint32(offset), uint32(limit), func(v *domain.Video) (bool, error) {
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
