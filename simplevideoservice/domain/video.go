package domain

import (
	"github.com/google/uuid"
)

type Video struct {
	id         string
	name       string
	path       string
	screenshot string
	duration   int
}

func GenerateVideoId() string {
	return uuid.New().String()
}

func MakeVideo(id, name, path, screenshot string, duration int) Video {
	return Video{id, name, path, screenshot, duration}
}

func (v Video) Id() string {
	return v.id
}

func (v Video) Name() string {
	return v.name
}

func (v Video) Duration() int {
	return v.duration
}

func (v Video) ScreenShotUrl() string {
	return v.screenshot
}

func (v Video) VideoUrl() string {
	return v.path
}
