package usecases

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
)

type VideoLister interface {
	List(handler func(v *domain.Video) (bool, error)) error
}

type videoLister struct {
	videos repository.Videos
}

func MakeVideoLister(videos repository.Videos) VideoLister {
	return &videoLister{videos}
}

func (l *videoLister) List(handler func(v *domain.Video) (bool, error)) error {
	return l.videos.Enumerate(handler)
}
