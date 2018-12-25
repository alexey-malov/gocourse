package usecases

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/alexey-malov/gocourse/simplevideoservice/storage"
	"io"
)

type Uploader interface {
	Upload(name string, content io.Reader) error
}

type uploader struct {
	videos repository.Videos
	stg    storage.Storage
}

func MakeUploader(videos repository.Videos, stg storage.Storage) Uploader {
	return &uploader{videos, stg}
}

func (u *uploader) Upload(name string, content io.Reader) error {
	id := domain.GenerateVideoId()

	folder, err := u.stg.MakeFolder(id)
	if err != nil {
		return err
	}

	file, err := folder.UploadFile("video.mp4", content)
	if err != nil {
		return err
	}

	v := domain.MakeUploadedVideo(id, name, file.Path())

	if err := u.videos.Add(*v); err != nil {
		return err
	}

	return nil
}
