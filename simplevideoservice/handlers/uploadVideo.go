package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/alexey-malov/gocourse/simplevideoservice/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const dirPath string = `C:\teaching\go\src\github.com\alexey-malov\gocourse\wwwroot\content`

func uploadVideo(vr repository.Videos, _ http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	// Обрабатываем ошибки

	contentType := header.Header.Get("Content-Type")
	// Убеждаемся, что пришел файл допустимого формата
	if contentType != "video/mp4" {
		// TODO
		return
	}

	fileName := header.Filename

	files := storage.MakeFiles(dirPath)
	_, id, err := files.Add(fileReader)
	if err != nil {
		log.Error(err)
		return
	}

	if err = vr.Add(model.MakeVideoItem(id, fileName, 42)); err != nil {
		log.Error(err)
		return
	}
}
