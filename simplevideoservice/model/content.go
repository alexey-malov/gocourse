package model

import "fmt"

type VideoItem struct {
	id       string
	name     string
	duration int
}

func MakeVideoItem(id string, name string, duration int) VideoItem {
	return VideoItem{id, name, duration}
}

func (v VideoItem) Id() string {
	return v.id
}

func (v VideoItem) Name() string {
	return v.name
}

func (v VideoItem) Duration() int {
	return v.duration
}

func (v VideoItem) ScreenShotUrl() string {
	return fmt.Sprintf("/content/%s/screen.jpg", v.id)
}

func (v VideoItem) VideoUrl() string {
	return fmt.Sprintf("/content/%s/index.mp4", v.id)
}
