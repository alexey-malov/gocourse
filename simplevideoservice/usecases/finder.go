package usecases

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
)

type VideoFinder interface {
	Find(id string) (*domain.Video, error)
}

type videoFinder struct {
	videos repository.Videos
}

func MakeFinder(videos repository.Videos) VideoFinder {
	return &videoFinder{videos}
}

func (f *videoFinder) Find(id string) (*domain.Video, error) {
	return f.videos.Find(id)
}
