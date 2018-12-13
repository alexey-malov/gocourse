package handlers

import "fmt"

type videoItem struct {
	id       string
	name     string
	duration int
}

var videoItems = []videoItem{
	{"d290f1ee-6c54-4b01-90e6-d701748f0851", "Black Retrospetive Woman", 15},
	{"sldjfl34-dfgj-523k-jk34-5jk3j45klj34", "Go Rally TEASER-HD", 41},
	{"hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345", "Танцор", 92},
}

type videoItemCallback func(v videoItem) bool

func enumVideos(cb videoItemCallback) {
	for _, v := range videoItems {
		if !cb(v) {
			return
		}
	}
}

func (v videoItem) screenShotUrl() string {
	return fmt.Sprintf("/content/%s/screen.jpg", v.id)
}

func (v videoItem) videoUrl() string {
	return fmt.Sprintf("/content/%s/index.mp4", v.id)
}

func findVideo(id string) *videoItem {
	var result *videoItem
	enumVideos(func(v videoItem) bool {
		if v.id == id {
			result = &v
			return false
		}
		return true
	})
	return result
}
