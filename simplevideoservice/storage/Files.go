package storage

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type Files interface {
	Add(r io.Reader) (path string, id string, err error)
}

type files struct {
	baseDir string
}

func MakeFiles(baseDir string) Files {
	return &files{baseDir}
}

func (f *files) Add(r io.Reader) (path string, id string, err error) {
	file, id, err := f.createFile("index.mp4")
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Failed to close file", err)
		}
	}()

	if _, err = io.Copy(file, r); err != nil {
		log.Error("Failed to write to file. Err: ", err)
		return "", "", err
	}
	return file.Name(), id, nil
}

func (f *files) createFile(name string) (file *os.File, fileId string, err error) {
	fileId = uuid.New().String()
	videoDir := filepath.Join(f.baseDir, fileId)

	if err := os.Mkdir(videoDir, os.ModeDir); err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(videoDir, name)
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, fileId, err
}
