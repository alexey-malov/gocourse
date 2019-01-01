package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

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
