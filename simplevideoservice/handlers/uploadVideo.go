package handlers

import (
	"database/sql"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const dirPath string = `C:\teaching\go\src\github.com\alexey-malov\gocourse\wwwroot\content`

func uploadVideo(db *sql.DB, _ http.ResponseWriter, r *http.Request) {
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

	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		// TODO
		return
	}

	vr := makeVideoRepository(db)
	if vr.addVideo(videoItem{fileId, fileName, 42}) != nil {
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
