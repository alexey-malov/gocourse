package storage

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
