package structs

// Data for passline
type Data struct {
	Items []Item
}

// Item structure
type Item struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
}
