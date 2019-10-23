package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type LocalStorage struct {
	storageDir  string
	storageFile string
}

func (ls *LocalStorage) Init() error {
	mainDir, _ := getMainDir()

	ls.storageDir = path.Join(mainDir, "storage")
	ls.storageFile = path.Join(ls.storageDir, "storage.json")

	ls.ensureDirectories()
	return nil
}

// Get item by name
func (ls LocalStorage) GetByName(name string) (Item, error) {
	data := ls.getData()
	for i := 0; i < len(data.Items); i++ {
		if data.Items[i].Name == name {
			return data.Items[i], nil
		}
	}

	return Item{}, errors.New("Item not found")
}

func (ls LocalStorage) GetByIndex(index int) (Item, error) {
	data := ls.getData()
	if index < 0 && index > len(data.Items) {
		return Item{}, errors.New("Out of index")
	}

	return data.Items[index], nil
}

func (ls LocalStorage) GetAll() ([]Item, error) {
	data := ls.getData()
	return data.Items, nil
}

// Add data
func (ls LocalStorage) AddItem(website Item) error {
	data := ls.getData()
	data.Items = append(data.Items, website)
	ls.setData(data)
	return nil
}

func (ls LocalStorage) AddCredential(name string, credential Credential) error {
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

func (ls LocalStorage) DeleteItem(item Item) error {
	data := ls.getData()
	index := getIndexOfItem(data.Items, item)
	data.Items = removeFromItems(data.Items, index)
	ls.setData(data)
	return nil
}

func (ls LocalStorage) DeleteCredential(item Item, credential Credential) error {
	data := ls.getData()
	indexItem := getIndexOfItem(data.Items, item)
	if indexItem == -1 {
		return errors.New("Item not found")
	}

	indexCredential := getIndexOfCredential(data.Items[indexItem].Credentials, credential)
	if indexCredential == -1 {
		return errors.New("Item not found")
	}

	data.Items[indexItem].Credentials = removeFromCredentials(data.Items[indexItem].Credentials, indexCredential)
	ls.setData(data)
	return nil
}

func (ls LocalStorage) ensureDirectories() {
	ls.ensureMainDir()
	ls.ensureStorageDir()
}

func (ls LocalStorage) ensureMainDir() error {
	mainDir, err := getMainDir()
	if err != nil {
		return err
	}

	_, err = os.Stat(mainDir)
	if err != nil {
		err := os.MkdirAll(mainDir, os.ModePerm)
		if err != nil {
			println("Cant create directory")
		}
	}

	return nil
}

func (ls LocalStorage) ensureStorageDir() {
	_, err := os.Stat(ls.storageDir)
	if err != nil {
		err := os.Mkdir(ls.storageDir, os.ModePerm)
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
