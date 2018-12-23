package storage

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type UploadedFile struct {
	name string
	path string
}

func (f UploadedFile) Name() string {
	return f.name
}

func (f UploadedFile) Path() string {
	return f.path
}

type Folder interface {
	UploadFile(name string, content io.Reader) (*UploadedFile, error)
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

type Storage interface {
	MakeFolder(name string) (Folder, error)
}

type storage struct {
	baseDir    string
	contentDir string
}

func MakeStorage(baseDir string, contentDir string) Storage {
	s := storage{baseDir, contentDir}
	return &s
}

func (s *storage) MakeFolder(name string) (Folder, error) {
	p := filepath.Join(s.baseDir, s.contentDir, name)
	if err := os.Mkdir(p, os.ModeDir); err != nil {
		return nil, err
	}
	f := folder{name, fmt.Sprintf("/%s/%s", s.contentDir, name), p}
	return &f, nil
}

/*
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

func (f *files) createFile(id, name string) (file *os.File, err error) {
	videoDir := filepath.Join(f.baseDir, id)

	if err := os.Mkdir(videoDir, os.ModeDir); err != nil {
		return nil, "", err
	}

	filePath := filepath.Join(videoDir, name)
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, fileId, err
}
*/
