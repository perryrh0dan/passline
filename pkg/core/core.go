package core

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/perryrh0dan/passline/pkg/config"
	"github.com/perryrh0dan/passline/pkg/crypt"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
)

type Core struct {
	config  *config.Config
	storage storage.Storage
}

func NewCore() (*Core, error) {
	c := new(Core)
	c.config, _ = config.Get()
	var err error
	switch c.config.Storage {
	case "firestore":
		c.storage, err = storage.NewFirestore()
		if err != nil {
			return nil, err
		}
	default:
		c.storage, err = storage.NewLocalStorage()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Core) CheckPassword(password []byte) (bool, error) {
	data, err := c.storage.GetAllItems()
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	item, err := c.storage.GetItemByIndex(0)
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

func (c *Core) CreateBackup(path string) error {
	items, err := c.storage.GetAllItems()
	if err != nil {
		return err
	}

	path = path + ".json"
	data := storage.Data{Items: items}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)

	return nil
}

func (c *Core) CreateItem(name, username, password string, globalPassword []byte) (storage.Credential, error) {

	// Check global password.
	valid, err := c.CheckPassword(globalPassword)
	if err != nil || !valid {
		return storage.Credential{}, err
	}

	cryptedPassword, err := crypt.AesGcmEncrypt(globalPassword, password)
	if err != nil {
		return storage.Credential{}, err
	}

	// Create Credentials
	credential := storage.Credential{Username: username, Password: cryptedPassword}

	err = c.addCredential(name, credential)
	if err != nil {
		return storage.Credential{}, err
	}

	credential.Password = password
	return credential, nil
}

func (c *Core) addCredential(name string, credential storage.Credential) error {
	// Check if item already exists
	_, err := c.storage.GetItemByName(name)
	if err != nil {
		// Generate new item entry
		item := storage.Item{Name: name, Credentials: []storage.Credential{credential}}
		err = c.storage.CreateItem(item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// TODO check if credential already exists
		// Add to existing item
		err := c.storage.AddCredential(name, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	return nil
}

func (c *Core) DecryptPassword(credential *storage.Credential, globalPassword []byte) error {
	// Decrypt passwords
	var err error
	credential.Password, err = crypt.AesGcmDecrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) GenerateItem(name, username string, globalPassword []byte) (storage.Credential, error) {
	// Check global password.
	valid, err := c.CheckPassword(globalPassword)
	if err != nil || !valid {
		return storage.Credential{}, err
	}

	// Generate password and crypt password
	password, err := crypt.GeneratePassword(20)
	if err != nil {
		return storage.Credential{}, err
	}

	cryptedPassword, err := crypt.AesGcmEncrypt(globalPassword, password)

	// Create Credentials
	credential := storage.Credential{Username: username, Password: cryptedPassword}

	err = c.addCredential(name, credential)
	if err != nil {
		return storage.Credential{}, err
	}

	credential.Password = password
	return credential, nil
}

func (c *Core) DeleteItem(name, username string) error {
	item, err := c.storage.GetItemByName(name)
	if err != nil {
		return err
	}

	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		return err
	}

	err = c.storage.DeleteCredential(item, credential)
	if err != nil {
		os.Exit(0)
	}

	return nil
}

func (c *Core) EditItem(name, username, newUsername string) error {
	item, err := c.storage.GetItemByName(name)
	if err != nil {
		return err
	}

	for i := 0; i < len(item.Credentials); i++ {
		if item.Credentials[i].Username == username {
			item.Credentials[i].Username = newUsername
		}
	}

	err = c.storage.UpdateItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) GetSites() ([]storage.Item, error) {
	items, err := c.storage.GetAllItems()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Core) GetSiteNames() ([]string, error) {
	items, err := c.GetSites()
	if err != nil {
		return nil, err
	}

	var names []string
	for _, item := range items {
		names = append(names, item.Name)
	}

	return names, nil
}

func (c *Core) GetSite(name string) (storage.Item, error) {
	item, err := c.storage.GetItemByName(name)
	if err != nil {
		return storage.Item{}, err
	}

	return item, nil
}

func (c *Core) Exists(name, username string) (bool, error) {
	item, err := c.storage.GetItemByName(name)
	if err == nil {
		_, err = item.GetCredentialByUsername(username)
		if err == nil {
			renderer.InvalidUsername(name, username)
			return true, nil
		}
	}

	return false, nil
}
