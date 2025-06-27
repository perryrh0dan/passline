package util

import "os"

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	IsNotExist(error error) bool
}

type OSFileSystem struct{}

func (OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
func (OSFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}
func (OSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
