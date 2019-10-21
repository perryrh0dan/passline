package storage

// Data for passline
type Data struct {
	Items []Item
}

// Item structure
type Item struct {
	Name   	 string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
}

// Storage interface
type Storage interface {
	Get(website string) (Item, error)
	GetAll() ([]Item, error)
	Add(item Item)
}
