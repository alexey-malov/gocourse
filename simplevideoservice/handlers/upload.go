package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (router *MyRouter) upload(w http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	defer fileReader.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = router.uploader.Upload(header.Filename, fileReader); err != nil {
		log.Error(err)
		return
	}
}
