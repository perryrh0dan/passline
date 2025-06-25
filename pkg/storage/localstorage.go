package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"

	"passline/pkg/cli/input"
	"passline/pkg/config"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
)

type LocalStorage struct {
	storageDir   string
	storageFile  string
	keyFile      string
	items        []Item
	decryptedKey []byte
}

func NewLocalStorage() (*LocalStorage, error) {
	mainDir := config.Directory()

	storageDir := path.Join(mainDir, "storage")
	storageFile := path.Join(storageDir, "storage")
	keyFile := path.Join(storageDir, "key")

	ensureDirectories(storageDir)

	ls := LocalStorage{storageDir: storageDir, storageFile: storageFile, keyFile: keyFile}

	return &ls, nil
}

func (ls *LocalStorage) GetItemByName(ctx context.Context, name string) (Item, error) {
	items, err := ls.getItems(ctx)
	if err != nil {
		return Item{}, err
	}

	for i := 0; i < len(items); i++ {
		if items[i].Name == name {
			return items[i], nil
		}
	}

	return Item{}, errors.New("Item not found")
}

func (ls *LocalStorage) GetItemByIndex(ctx context.Context, index int) (Item, error) {
	items, err := ls.getItems(ctx)
	if err != nil {
		return Item{}, err
	}

	if index < 0 && index > len(items) {
		return Item{}, errors.New("Out of index")
	}

	return items[index], nil
}

func (ls *LocalStorage) GetAllItems(ctx context.Context) ([]Item, error) {
	items, err := ls.getItems(ctx)
	if err != nil {
		return nil, err
	}

	sort.Sort(ByName(items))
	return items, nil
}

func (ls *LocalStorage) AddCredential(ctx context.Context, name string, credential Credential) error {
	// Check if item already exists
	_, err := ls.GetItemByName(ctx, name)
	if err != nil {
		// Generate new item entry
		item := Item{Name: name, Credentials: []Credential{credential}}
		ls.createItem(ctx, item)
		return nil
	}

	// if item exists just append
	items, err := ls.getItems(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < len(items); i++ {
		if items[i].Name == name {
			for y := 0; y < len(items[i].Credentials); y++ {
				if items[i].Credentials[y].Username == credential.Username {
					return errors.New("Username already exists")
				}
			}
			items[i].Credentials = append(items[i].Credentials, credential)
			break
		}
	}

	ls.SetItems(ctx, items)
	return nil
}

func (ls *LocalStorage) DeleteCredential(ctx context.Context, item Item, username string) error {
	items, err := ls.getItems(ctx)
	if err != nil {
		return err
	}

	indexItem := getIndexOfItem(items, item.Name)
	if indexItem == -1 {
		return errors.New("Item not found")
	}

	if len(items[indexItem].Credentials) > 1 {
		indexCredential := getIndexOfCredential(items[indexItem].Credentials, username)
		if indexCredential == -1 {
			return errors.New("Item not found")
		}

		items[indexItem].Credentials = removeFromCredentials(items[indexItem].Credentials, indexCredential)
		ls.SetItems(ctx, items)
	} else {
		ls.deleteItem(ctx, items[indexItem])
	}
	return nil
}

func (ls *LocalStorage) UpdateItem(ctx context.Context, item Item) error {
	// TODO check if username is valid

	ls.deleteItem(ctx, item)
	ls.createItem(ctx, item)

	return nil
}

func (ls *LocalStorage) GetKey(ctx context.Context) (string, error) {
	key, err := os.ReadFile(ls.keyFile)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return string(key), nil
}

func (ls *LocalStorage) GetDecryptedKey(ctx context.Context, reason string) ([]byte, error) {
	if ls.decryptedKey != nil {
		return ls.decryptedKey, nil
	}

	// Get encrypted content encryption key from store
	encryptedEncryptionKey, err := ls.GetKey(ctx)
	if err != nil {
		return []byte{}, err
	}

	if encryptedEncryptionKey != "" {
		encryptionKey, err := input.MasterPassword(encryptedEncryptionKey, reason)
		if err != nil {
			return []byte{}, err
		}

		ls.decryptedKey = encryptionKey

		return encryptionKey, nil
	}

	decryptedEncryptionKey, err := crypt.GenerateKey()
	if err != nil {
		return []byte{}, err
	}

	password := input.Password("Enter master password: ")
	fmt.Println()
	passwordTwo := input.Password("Enter master password again: ")
	fmt.Println()

	if string(password) != string(passwordTwo) {
		return []byte{}, err
	}

	encryptedEncryptionKey, err = crypt.EncryptKey(password, decryptedEncryptionKey)
	if err != nil {
		return []byte{}, err
	}
	ls.SetKey(ctx, encryptedEncryptionKey)

	ls.decryptedKey = []byte(decryptedEncryptionKey)

	return ls.decryptedKey, nil
}

func (ls *LocalStorage) SetKey(ctx context.Context, key string) error {
	err := os.WriteFile(ls.keyFile, []byte(key), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ls *LocalStorage) GetRawItems(ctx context.Context) (json.RawMessage, error) {
	file, err := os.ReadFile(ls.storageFile)
	if err != nil {
		return nil, err
	}

	var js json.RawMessage
	if json.Unmarshal(file, &js) != nil {
		return json.RawMessage(fmt.Sprintf("\"%s\"", file)), nil
	}

	return json.RawMessage(file), nil
}

func ensureDirectories(storageDir string) {
	ensureStorageDir(storageDir)
}

func ensureStorageDir(storageDir string) {
	_, err := os.Stat(storageDir)
	if err != nil {
		err := os.Mkdir(storageDir, os.ModePerm)
		if err != nil {
			println("Cant create directory")
		}
	}
}

func (ls *LocalStorage) createItem(ctx context.Context, item Item) error {
	items, err := ls.getItems(ctx)
	if err != nil {
		return err
	}

	items = append(items, item)
	ls.SetItems(ctx, items)

	return nil
}

func (ls *LocalStorage) deleteItem(ctx context.Context, item Item) error {
	items, err := ls.getItems(ctx)
	if err != nil {
		return err
	}

	index := getIndexOfItem(items, item.Name)
	items = removeFromItems(items, index)
	ls.SetItems(ctx, items)

	return nil
}

func (ls *LocalStorage) SetItems(ctx context.Context, items []Item) error {

	file, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	encryption := ctxutil.GetEncryption(ctx)
	if encryption == config.FullEncryption {
		key, err := ls.GetDecryptedKey(ctx, "encrypt the password")
		if err != nil {
			return err
		}

		encryptedResult, err := crypt.AesGcmEncrypt(key, string(file))
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		file = []byte(fmt.Sprintf("%s", encryptedResult))
	}

	err = os.WriteFile(ls.storageFile, file, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func (ls *LocalStorage) getItems(ctx context.Context) ([]Item, error) {
	if len(ls.items) > 0 {
		return ls.items, nil
	}

	_, err := os.Stat(ls.storageFile)
	if err != nil {
		return []Item{}, nil
	}

	file, _ := os.ReadFile(ls.storageFile)
	rawItems := json.RawMessage(file)

	// Check if the file is full encrypted by checking if it is a valid json
	var js json.RawMessage
	if json.Unmarshal(file, &js) != nil {
		decryptedKey, err := ls.GetDecryptedKey(ctx, "to decrypt your vault")
		if err != nil {
			return nil, err
		}

		decryptedItems, err := crypt.AesGcmDecrypt(decryptedKey, string(file))
		if err != nil {
			return nil, err
		}

		rawItems = json.RawMessage(decryptedItems)
	}

	items := []Item{}
	err = json.Unmarshal(rawItems, &items)
	if err != nil {
		return nil, err
	}

	ls.items = items

	return ls.items, nil
}
