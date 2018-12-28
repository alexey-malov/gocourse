package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func upload(uploader usecases.Uploader, w http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer closeFile(fileReader)

	contentType := header.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = uploader.Upload(header.Filename, fileReader); err != nil {
		log.Error(err)
		return
	}
}

func closeFile(closer io.Closer) {
	func() {
		if err := closer.Close(); err != nil {
			log.Error(err)
		}
	}()
}
