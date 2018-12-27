package handlers

import (
	"bytes"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"io"
)

type mockUploader struct {
	fileName string
	content  string
	err      error
}

type mockVideos struct {
	videos     []*domain.Video
	findResult *domain.Video
	savedVideo *domain.Video
}

func (u *mockUploader) Upload(name string, content io.Reader) error {
	u.fileName = name
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(content); err != nil {
		return err
	}
	u.content = buf.String()
	return u.err
}

func (r *mockVideos) Enumerate(handler func(v *domain.Video) bool) error {
	for _, v := range r.videos {
		if !handler(v) {
			return nil
		}
	}
	return nil
}

func (r *mockVideos) Find(id string) (*domain.Video, error) {
	return r.findResult, nil
}

func (r *mockVideos) Add(v *domain.Video) error {
	r.videos = append(r.videos, v)
	return nil
}

func (r *mockVideos) EnumerateWithStatus(status domain.Status, handler func(v *domain.Video) bool) error {
	for _, v := range r.videos {
		if v.Status() == status && !handler(v) {
			return nil
		}
	}
	return nil
}

func (r *mockVideos) SaveVideo(v *domain.Video) error {
	r.savedVideo = v
	return nil
}
