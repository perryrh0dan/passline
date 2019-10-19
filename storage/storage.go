package storage

// ItemStorage for all websites
type ItemStorage struct {
	Items []Item `json:"items"`
}

// Item in storage
type Item struct {
	// Website name
	Website string `json:"website"`
	// Password (encoded)
	Password string `json:"password"`
}

// Storage interface
type Storage interface {
	get(website string) (Item, error)
	set(item Item)
}
