package storage

type Data struct {
	Items []Item
}

// Item structure
type Item struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Storage interface
type Storage interface {
	Init() error
	GetByName(string) (Item, error)
	GetByIndex(int) (Item, error)
	GetAll() ([]Item, error)
	Add(Item) error
	Delete(Item) error
}
