package storage

import "github.com/perryrh0dan/passline/pkg/structs"

// Storage interface
type Storage interface {
	GetByName(string) (structs.Item, error)
	GetByIndex(int) (structs.Item, error)
	GetAll() ([]structs.Item, error)
	Add(structs.Item) error
	Delete(structs.Item) error
}
