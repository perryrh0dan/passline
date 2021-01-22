package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"passline/pkg/config"
)

type LocalStorage struct {
	storageDir  string
	storageFile string
}

func NewLocalStorage() (*LocalStorage, error) {
	mainDir := config.Directory()

	storageDir := path.Join(mainDir, "storage")
	storageFile := path.Join(storageDir, "storage.json")

	ensureDirectories(storageDir, storageFile)
	return &LocalStorage{storageDir: storageDir, storageFile: storageFile}, nil
}

// Get item by name
func (ls *LocalStorage) GetItemByName(ctx context.Context, name string) (Item, error) {
	data := ls.getData()
	for i := 0; i < len(data.Items); i++ {
		if data.Items[i].Name == name {
			return *data.Items[i], nil
		}
	}

	return Item{}, errors.New("Item not found")
}

func (ls *LocalStorage) GetItemByIndex(ctx context.Context, index int) (Item, error) {
	data := ls.getData()
	if index < 0 && index > len(data.Items) {
		return Item{}, errors.New("Out of index")
	}

	return data.Items[index], nil
}

func (ls *LocalStorage) GetAllItems(ctx context.Context) ([]Item, error) {
	data := ls.getData()
	sort.Sort(ByName(data.Items))
	return data.Items, nil
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
	data := ls.getData()
	for i := 0; i < len(data.Items); i++ {
		if data.Items[i].Name == name {
			for y := 0; y < len(data.Items[i].Credentials); y++ {
				if data.Items[i].Credentials[y].Username == credential.Username {
					return errors.New("Username already exists")
				}
			}
			data.Items[i].Credentials = append(data.Items[i].Credentials, credential)
			break
		}
	}

	ls.setData(data)
	return nil
}

func (ls *LocalStorage) DeleteCredential(ctx context.Context, item Item, credential Credential) error {
	data := ls.getData()
	indexItem := getIndexOfItem(data.Items, item)
	if indexItem == -1 {
		return errors.New("Item not found")
	}

	if len(data.Items[indexItem].Credentials) > 1 {
		indexCredential := getIndexOfCredential(data.Items[indexItem].Credentials, credential)
		if indexCredential == -1 {
			return errors.New("Item not found")
		}

		data.Items[indexItem].Credentials = removeFromCredentials(data.Items[indexItem].Credentials, indexCredential)
		ls.setData(data)
	} else {
		ls.deleteItem(data.Items[indexItem])
	}
	return nil
}

func (ls *LocalStorage) UpdateItem(ctx context.Context, item Item) error {
	// TODO check if username is valid

	ls.deleteItem(item)
	ls.createItem(ctx, item)

	return nil
}

func (ls *LocalStorage) SetData(ctx context.Context, data Data) error {
	ls.setData(data)
	return nil
}

func (ls *LocalStorage) GetKey(ctx context.Context) (string, error) {
	data := ls.getData()
	return data.Key, nil
}

func (ls *LocalStorage) SetKey(ctx context.Context, key string) error {
	data := ls.getData()
	data.Key = key
	ls.setData(data)
	return nil
}

func (ls *LocalStorage) createItem(ctx context.Context, item Item) {
	data := ls.getData()
	data.Items = append(data.Items, item)
	ls.setData(data)
}

func (ls *LocalStorage) deleteItem(item Item) {
	data := ls.getData()
	index := getIndexOfItem(data.Items, item)
	data.Items = removeFromItems(data.Items, index)
	ls.setData(data)
}

func ensureDirectories(storageDir, storageFile string) {
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

func (ls LocalStorage) getData() Data {
	data := Data{}

	_, err := os.Stat(ls.storageFile)
	if err == nil {
		file, _ := ioutil.ReadFile(ls.storageFile)
		_ = json.Unmarshal([]byte(file), &data)
	}

	return data
}

func (ls LocalStorage) setData(data Data) {
	_, err := os.Stat(ls.storageDir)
	if err == nil {
		file, _ := json.MarshalIndent(data, "", " ")
		_ = ioutil.WriteFile(ls.storageFile, file, 0644)
	}
}
