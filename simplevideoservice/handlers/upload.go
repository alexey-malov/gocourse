package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *handlerBase) upload(w http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	defer func() {
		if err := fileReader.Close(); err != nil {
			log.Error(err)
		}
	}()

	contentType := header.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.uploader.Upload(header.Filename, fileReader); err != nil {
		log.Error(err)
		return
	}
}
