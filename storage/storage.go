package storage

// Data for passline
type Data struct {
	Websites []Website
}

// Website structure
type Website struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
}

// Storage interface
type Storage interface {
	Get(website string) (Website, error)
	GetAll() ([]Website, error)
	Add(item Website)
}
