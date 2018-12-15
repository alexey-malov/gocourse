package handlers

import (
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const dirPath string = "uploads"

func uploadVideo(_ http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	// Обрабатываем ошибки

	contentType := header.Header.Get("Content-Type")
	// Убеждаемся, что пришел файл допустимого формата
	if contentType != "video/mp4" {
		// TODO
		return
	}

	fileName := header.Filename
	file, err := createFile(fileName)
	if err != nil {
		// TODO
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		// TODO
		return
	}
}

func createFile(name string) (*os.File, error) {
	videoDir := filepath.Join(dirPath, uuid.New().String())
	if err := os.Mkdir(videoDir, os.ModeDir); err != nil {
		return nil, err
	}

	filePath := filepath.Join(videoDir, name)
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}
