package app

import "github.com/alexey-malov/gocourse/simplevideoservice/storage"

const dirPath string = `C:\teaching\go\src\github.com\alexey-malov\gocourse\wwwroot`

func MakeStorage() storage.Storage {
	return storage.MakeStorage(dirPath, "content")
}
