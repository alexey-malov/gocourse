package app

import (
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/storage"
	"os"
)

func MakeStorage() (storage.Storage, error) {
	const storageRootEnvVar = "SIMPLE_VIDEO_SERVICE_STORAGE"
	storagePath := os.Getenv(storageRootEnvVar)
	if storagePath == "" {
		return nil, fmt.Errorf("No %s environment variable", storageRootEnvVar)
	}

	return storage.MakeStorage(storagePath, "content"), nil
}
