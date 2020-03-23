package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"passline/pkg/config"
	"passline/pkg/crypt"
	"passline/pkg/renderer"
	"passline/pkg/storage"
)

type Core struct {
	Config  *config.Config
	storage storage.Storage
}

func NewCore(ctx context.Context) (*Core, error) {
	c := new(Core)
	c.Config, _ = config.Get()
	var err error
	switch c.Config.Storage {
	case "firestore":
		c.storage, err = storage.NewFirestore(ctx)
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

func (c *Core) CheckPassword(ctx context.Context, password []byte) (bool, error) {
	data, err := c.storage.GetAllItems(ctx)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	item, err := c.storage.GetItemByIndex(ctx, 0)
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

func (c *Core) CreateBackup(ctx context.Context, path string) error {
	items, err := c.storage.GetAllItems(ctx)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}

	time := time.Now()
	data := storage.Backup{Date: time, Items: items}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)

	return nil
}

func (c *Core) RestoreBackup(ctx context.Context, path string) error {
	data := storage.Data{}

	_, err := os.Stat(path)
	if err != nil {
		renderer.InvalidFilePath()
		return err
	}

	file, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal([]byte(file), &data)

	err = c.storage.SetData(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) AddItem(ctx context.Context, name, username, password string, recoveryCodes []string, globalPassword []byte) (storage.Credential, error) {

	// Check global password.
	valid, err := c.CheckPassword(ctx, globalPassword)
	if err != nil || !valid {
		return storage.Credential{}, err
	}

	// Create Credentials
	credential := storage.Credential{Username: username, Password: password, RecoveryCodes: recoveryCodes}

	err = c.EncryptCredential(&credential, globalPassword)
	if err != nil {
		return storage.Credential{}, err
	}

	err = c.addCredential(ctx, name, credential)
	if err != nil {
		return storage.Credential{}, err
	}

	credential.Password = password
	return credential, nil
}

func (c *Core) addCredential(ctx context.Context, name string, credential storage.Credential) error {
	// Check if item already exists
	_, err := c.storage.GetItemByName(ctx, name)
	if err != nil {
		// Generate new item entry
		item := storage.Item{Name: name, Credentials: []storage.Credential{credential}}
		err = c.storage.CreateItem(ctx, item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// TODO check if credential already exists
		// Add to existing item
		err := c.storage.AddCredential(ctx, name, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	return nil
}

func (c *Core) DecryptCredential(credential *storage.Credential, globalPassword []byte) error {
	err := c.DecryptPassword(credential, globalPassword)
	if err != nil {
		return err
	}

	err = c.DecryptRecoveryCodes(credential, globalPassword)
	if err != nil {
		return err
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

func (c *Core) DecryptRecoveryCodes(credential *storage.Credential, globalPassword []byte) error {
	var decryptedRecoveryCodes = make([]string, 0)
	for _, c := range credential.RecoveryCodes {
		decryptedRecoveryCode, err := crypt.AesGcmDecrypt(globalPassword, c)
		if err != nil {
			return err
		}
		decryptedRecoveryCodes = append(decryptedRecoveryCodes, decryptedRecoveryCode)
	}

	credential.RecoveryCodes = decryptedRecoveryCodes
	return nil
}

func (c *Core) EncryptCredential(credential *storage.Credential, globalPassword []byte) error {
	err := c.EncryptPassword(credential, globalPassword)
	if err != nil {
		return err
	}

	err = c.EncryptRecoveryCodes(credential, globalPassword)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) EncryptPassword(credential *storage.Credential, globalPassword []byte) error {
	var err error
	credential.Password, err = crypt.AesGcmEncrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) EncryptRecoveryCodes(credential *storage.Credential, globalPassword []byte) error {
	var encryptedRecoveryCodes = make([]string, 0)

	for _, c := range credential.RecoveryCodes {
		encryptedRecoveryCode, err := crypt.AesGcmEncrypt(globalPassword, c)
		if err != nil {
			return err
		}
		encryptedRecoveryCodes = append(encryptedRecoveryCodes, encryptedRecoveryCode)
	}

	credential.RecoveryCodes = encryptedRecoveryCodes
	return nil
}

func (c *Core) GenerateItem(ctx context.Context, name, username string, recoveryCodes []string, globalPassword []byte) (storage.Credential, error) {
	// Generate password and crypt password
	password, err := crypt.GeneratePassword(20)
	if err != nil {
		return storage.Credential{}, err
	}

	return c.AddItem(ctx, name, username, password, recoveryCodes, globalPassword)
}

func (c *Core) DeleteItem(ctx context.Context, name, username string) error {
	item, err := c.storage.GetItemByName(ctx, name)
	if err != nil {
		return err
	}

	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		return err
	}

	err = c.storage.DeleteCredential(ctx, item, credential)
	if err != nil {
		os.Exit(0)
	}

	return nil
}

func (c *Core) EditItem(ctx context.Context, name, username string, updatedCredential storage.Credential, globalPassword []byte) error {
	item, err := c.storage.GetItemByName(ctx, name)
	if err != nil {
		return err
	}

	// find credential in item
	var credential *storage.Credential

	for i := 0; i < len(item.Credentials); i++ {
		if item.Credentials[i].Username == username {
			credential = &item.Credentials[i]
			break
		}
	}

	*credential = updatedCredential

	err = c.EncryptCredential(credential, globalPassword)
	if err != nil {
		return err
	}

	err = c.storage.UpdateItem(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) GetSites(ctx context.Context) ([]storage.Item, error) {
	items, err := c.storage.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Core) GetSiteNames(ctx context.Context) ([]string, error) {
	items, err := c.GetSites(ctx)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, item := range items {
		names = append(names, item.Name)
	}

	return names, nil
}

func (c *Core) GetSite(ctx context.Context, name string) (storage.Item, error) {
	item, err := c.storage.GetItemByName(ctx, name)
	if err != nil {
		return storage.Item{}, err
	}

	return item, nil
}

func (c *Core) Exists(ctx context.Context, name, username string) (bool, error) {
	item, err := c.storage.GetItemByName(ctx, name)
	if err == nil {
		_, err = item.GetCredentialByUsername(username)
		if err == nil {
			renderer.InvalidUsername(name, username)
			return true, nil
		}
	}

	return false, nil
}
