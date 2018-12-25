package domain

import (
	"github.com/google/uuid"
)

type Status int

const (
	StatusUploaded Status = iota + 1
	StatusProcessing
	StatusReady
	StatusDeleted
	StatusError
)

type Video struct {
	id        string
	name      string
	path      string
	thumbnail string
	duration  int
	status    Status
}

func GenerateVideoId() string {
	return uuid.New().String()
}

func MakeVideo(id, name, path, screenshot string, duration int, status Status) *Video {
	return &Video{id, name, path, screenshot, duration, status}
}

func MakeUploadedVideo(id, name, path string) *Video {
	return &Video{id, name, path, "", 0, StatusUploaded}
}

func (v *Video) Id() string {
	return v.id
}

func (v *Video) Status() Status {
	return v.status
}

func (v *Video) SetStatus(s Status) {
	v.status = s
}

func (v *Video) Name() string {
	return v.name
}

func (v *Video) Duration() int {
	return v.duration
}

func (v *Video) ThumbnailURL() string {
	return v.thumbnail
}

func (v *Video) VideoUrl() string {
	return v.path
}

func (v *Video) SetThumbnailURL(url string) {
	v.thumbnail = url
}

func (v *Video) SetDuration(duration int) {
	v.duration = duration
}
