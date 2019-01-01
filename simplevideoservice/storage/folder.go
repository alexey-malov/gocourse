package storage

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type Folder interface {
	UploadFile(name string, content io.Reader) (*UploadedFile, error)
	GetAbsFilePath(name string) string
	GetRelFilePath(name string) string
}

type folder struct {
	name    string
	relPath string
	absPath string
}

func safeCloseFile(f *os.File) {
	if err := f.Close(); err != nil {
		logrus.Error("Unable to close a file", err)
	}
}

func (f *folder) UploadFile(name string, content io.Reader) (*UploadedFile, error) {
	path := filepath.Join(f.absPath, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer safeCloseFile(file)

	if _, err := io.Copy(file, content); err != nil {
		logrus.Error("Failed to copy file", err)
		return nil, err
	}
	return &UploadedFile{name, fmt.Sprintf("%s/%s", f.relPath, name)}, nil
}

func (f *folder) GetAbsFilePath(name string) string {
	return filepath.Join(f.absPath, name)
}

func (f *folder) GetRelFilePath(name string) string {
	return fmt.Sprintf("%s/%s", f.relPath, name)
}
