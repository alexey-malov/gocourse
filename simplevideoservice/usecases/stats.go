package usecases

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/storage"
	"os"
)

type Stats interface {
	VideoPath(v domain.Video) (string, error)
}

type stats struct {
	stg storage.Storage
}

func MakeStats(stg storage.Storage) Stats {
	return &stats{stg}
}

func (s *stats) VideoPath(v domain.Video) (string, error) {
	p := s.stg.GetAbsPath(v.VideoUrl())
	_, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	return p, nil
}
