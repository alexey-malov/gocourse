package domain

import "fmt"

type Video struct {
	id       string
	name     string
	duration int
}

func MakeVideo(id string, name string, duration int) Video {
	return Video{id, name, duration}
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
	return fmt.Sprintf("/content/%s/screen.jpg", v.id)
}

func (v Video) VideoUrl() string {
	return fmt.Sprintf("/content/%s/index.mp4", v.id)
}
