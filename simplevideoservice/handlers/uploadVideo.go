package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (router *MyRouter) uploadVideo(_ http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	// Обрабатываем ошибки

	contentType := header.Header.Get("Content-Type")
	// Убеждаемся, что пришел файл допустимого формата
	if contentType != "video/mp4" {
		// TODO
		return
	}

	if err = router.uploader.Upload(header.Filename, fileReader); err != nil {
		log.Error(err)
		return
	}
}
