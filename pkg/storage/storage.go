package storage

import (
	"errors"
	"time"

	"passline/pkg/config"

	"golang.org/x/net/context"
)

// Storage interface
type Storage interface {
	GetItemByName(context.Context, string) (Item, error)
	GetItemByIndex(context.Context, int) (Item, error)
	GetAllItems(context.Context) ([]Item, error)
	CreateItem(context.Context, Item) error
	AddCredential(context.Context, string, Credential) error
	DeleteCredential(context.Context, Item, Credential) error
	UpdateItem(context.Context, Item) error
	SetData(context.Context, Data) error
}

// Data structure
type Data struct {
	Items []Item `json:"items"`
}

// Backup structure
type Backup struct {
	Date  time.Time `json:"date"`
	Items []Item    `json:"items"`
}

// Item structure
type Item struct {
	Name        string       `json:"name"`
	Credentials []Credential `json:"credentials"`
}

// For sorting items
type ByName []Item

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (item *Item) GetCredentialByUsername(username string) (Credential, error) {
	for i := 0; i < len(item.Credentials); i++ {
		if item.Credentials[i].Username == username {
			return item.Credentials[i], nil
		}
	}

	return Credential{}, errors.New("Not found")
}

func (item *Item) GetCredentialsName() ([]string, error) {
	var list []string

	for i := 0; i < len(item.Credentials); i++ {
		list = append(list, item.Credentials[i].Username)
	}

	return list, nil
}

func (item *Item) GetUsernameArray() []string {
	var creds []string
	for _, cred := range item.Credentials {
		creds = append(creds, cred.Username)
	}

	return creds
}

type Credential struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

func getMainDir() (string, error) {
	config, err := config.Get()
	if err != nil {
		return "", err
	}

	return config.Directory, nil
}

func removeFromItems(s []Item, i int) []Item {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func removeFromCredentials(s []Credential, i int) []Credential {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func getIndexOfItem(slice []Item, item Item) int {
	for p, v := range slice {
		if v.Name == item.Name {
			return p
		}
	}
	return -1
}

func getIndexOfCredential(slice []Credential, credential Credential) int {
	for p, v := range slice {
		if v.Username == credential.Username {
			return p
		}
	}
	return -1
}
