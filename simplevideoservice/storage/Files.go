package storage

import (
	"errors"
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

type Storage interface {
	MakeFolder(name string) (Folder, error)
	GetFolder(name string) (Folder, error)
	GetAbsPath(path string) string
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
	return s.makeFolder(name), nil
}

func (s *storage) GetAbsPath(path string) string {
	p := filepath.Join(s.baseDir, path)
	return p
}

func (s *storage) makeFolder(name string) Folder {
	p := filepath.Join(s.baseDir, s.contentDir, name)
	return &folder{name, fmt.Sprintf("/%s/%s", s.contentDir, name), p}
}

func (s *storage) GetFolder(name string) (Folder, error) {
	p := filepath.Join(s.baseDir, s.contentDir, name)

	info, err := os.Stat(p)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Folder does not exist")
	}

	return s.makeFolder(name), nil
}

func (f *folder) GetAbsFilePath(name string) string {
	return filepath.Join(f.absPath, name)
}

func (f *folder) GetRelFilePath(name string) string {
	return fmt.Sprintf("%s/%s", f.relPath, name)
}
