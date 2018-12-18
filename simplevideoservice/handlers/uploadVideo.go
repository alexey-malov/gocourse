package handlers

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const dirPath string = `C:\teaching\go\src\github.com\alexey-malov\gocourse\wwwroot\content`

func uploadVideo(vr repository.VideoRepository, _ http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	// Обрабатываем ошибки

	contentType := header.Header.Get("Content-Type")
	// Убеждаемся, что пришел файл допустимого формата
	if contentType != "video/mp4" {
		// TODO
		return
	}

	fileName := header.Filename
	file, fileId, err := createFile("index.mp4")
	if err != nil {
		// TODO
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Failed to close uploaded file. Error is ", err)
		}
	}()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		// TODO
		return
	}

	if vr.AddVideo(model.MakeVideoItem(fileId, fileName, 42)) != nil {
		// TODO
		return
	}
}

func createFile(name string) (*os.File, string, error) {
	fileId := uuid.New().String()
	videoDir := filepath.Join(dirPath, fileId)
	if err := os.Mkdir(videoDir, os.ModeDir); err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(videoDir, name)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, fileId, err
}
