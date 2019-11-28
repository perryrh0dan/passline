package core

import (
	"os"

	"github.com/perryrh0dan/passline/pkg/config"
	"github.com/perryrh0dan/passline/pkg/crypt"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
)

type Passline struct {
	config *config.Config
	store  storage.Storage
}

func NewPassline() *Passline {
	pl := new(Passline)
	pl.config, _ = config.Get()
	var err error
	switch pl.config.Storage {
	case "firestore":
		pl.store, err = storage.NewFirestore()
		handle(err)
	default:
		pl.store, err = storage.NewLocalStorage()
		handle(err)
	}
	return pl
}

func (pl *Passline) CheckPassword(password []byte) (bool, error) {
	data, err := pl.store.GetAllItems()
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	item, err := pl.store.GetItemByIndex(0)
	if err != nil {
		return false, err
	}

	_, err = crypt.AesGcmDecrypt(password, item.Credentials[0].Password)
	if err != nil {
		renderer.InvalidPassword()
		return false, err
	}

	return true, nil
}

func handle(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func (pl *Passline) AddItem(name, username, password string, globalPassword []byte) (storage.Credential, error) {

	// Check global password.
	valid, err := pl.CheckPassword(globalPassword)
	if err != nil || !valid {
		handle(err)
	}

	cryptedPassword, err := crypt.AesGcmEncrypt(globalPassword, password)
	handle(err)

	// Create Credentials
	credential := storage.Credential{Username: username, Password: cryptedPassword}

	// Check if item already exists
	_, err = pl.store.GetItemByName(name)
	if err != nil {
		// Generate new item entry
		item := storage.Item{Name: name, Credentials: []storage.Credential{credential}}
		err = pl.store.AddItem(item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// TODO check if credential already exists
		// Add to existing item
		err := pl.store.AddCredential(name, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	credential.Password = password
	return credential, nil
}

func (pl *Passline) DecryptPassword(credential *storage.Credential, globalPassword []byte) error {
	// Decrypt passwords
	var err error
	credential.Password, err = crypt.AesGcmDecrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func (pl *Passline) GenerateItem(name, username string, globalPassword []byte) (storage.Credential, error) {
	// Check global password.
	valid, err := pl.CheckPassword(globalPassword)
	if err != nil || !valid {
		handle(err)
	}

	// Generate password and crypt password
	password, err := crypt.GeneratePassword(20)
	handle(err)

	cryptedPassword, err := crypt.AesGcmEncrypt(globalPassword, password)

	// Create Credentials
	credential := storage.Credential{Username: username, Password: cryptedPassword}

	// Check if item already exists
	_, err = pl.store.GetItemByName(name)
	if err != nil {
		// Generate new item entry
		item := storage.Item{Name: name, Credentials: []storage.Credential{credential}}
		err = pl.store.AddItem(item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// Add to existing item
		err := pl.store.AddCredential(name, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	credential.Password = password
	return credential, nil
}

func (pl *Passline) DeleteItem(name, username string) error {
	item, err := pl.store.GetItemByName(name)
	if err != nil {
		return err
	}

	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		return err
	}

	err = pl.store.DeleteCredential(item, credential)
	if err != nil {
		os.Exit(0)
	}

	return nil
}

func (pl *Passline) EditItem(name, username, newUsername string) error {
	item, err := pl.store.GetItemByName(name)
	if err != nil {
		return err
	}

	for i := 0; i < len(item.Credentials); i++ {
		if item.Credentials[i].Username == username {
			item.Credentials[i].Username = newUsername
		}
	}

	err = pl.store.UpdateItem(item)
	handle(err)

	return nil
}

func (pl *Passline) GetSites() ([]storage.Item, error) {
	items, err := pl.store.GetAllItems()
	if err != nil {
		handle(err)
	}

	return items, nil
}

func (pl *Passline) GetSiteNames() ([]string, error) {
	items, err := pl.GetSites()
	if err != nil {
		return nil, err
	}

	var names []string
	for _, item := range items {
		names = append(names, item.Name)
	}

	return names, nil
}

func (pl *Passline) GetSite(name string) (storage.Item, error) {
	item, err := pl.store.GetItemByName(name)
	if err != nil {
		handle(err)
	}

	return item, nil
}

func (pl *Passline) Exists(name, username string) (bool, error) {
	item, err := pl.store.GetItemByName(name)
	if err == nil {
		_, err = item.GetCredentialByUsername(username)
		if err == nil {
			renderer.InvalidUsername(name, username)
			return true, nil
		}
	}

	return false, nil
}
