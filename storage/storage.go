package storage

// Data for passline
type Data struct {
	Websites []Website
}

// Website structure
type Website struct {
	Domain   string `json:"domain`
	Username string `json:"username`
	Password string `json:"password`
}

// Storage interface
type Storage interface {
	get(website string) (Website, error)
	set(item Website)
}
