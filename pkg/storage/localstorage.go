package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/perryrh0dan/passline/pkg/config"
	"github.com/perryrh0dan/passline/pkg/structs"
	"github.com/perryrh0dan/passline/pkg/util"
)

type LocalStorage struct {
	storageDir  string
	storageFile string
}

func (ls *LocalStorage) Init() error {
	mainDir, _ := ls.getMainDir()

	ls.storageDir = path.Join(mainDir, "storage")
	ls.storageFile = path.Join(ls.storageDir, "storage.json")

	ls.ensureDirectories()
	return nil
}

// Get data
func (ls LocalStorage) GetByName(name string) (structs.Item, error) {
	data := ls.getData()
	for i := 0; i < len(data.Items); i++ {
		if data.Items[i].Name == name {
			return data.Items[i], nil
		}
	}

	return structs.Item{}, fmt.Errorf("No entry for website %s", name)
}

func (ls LocalStorage) GetByIndex(index int) (structs.Item, error) {
	data := ls.getData()
	if 0 <= index && index < len(data.Items) {
		return data.Items[index], nil
	} else {
		return structs.Item{}, errors.New("Out of index")
	}
}

func (ls LocalStorage) GetAll() ([]structs.Item, error) {
	data := ls.getData()
	return data.Items, nil
}

// Add data
func (ls LocalStorage) Add(website structs.Item) error {
	data := ls.getData()
	data.Items = append(data.Items, website)
	ls.setData(data)
	return nil
}

func (ls LocalStorage) Delete(item structs.Item) error {
	data := ls.getData()
	index := util.GetIndexOfItem(data.Items, item)
	data.Items = util.RemoveFromArray(data.Items, index)
	ls.setData(data)
	return nil
}

func (s LocalStorage) getMainDir() (string, error) {
	config, err := config.Get()
	if err != nil {
		return "", err
	}

	return config.Directory, nil
}

func (ls LocalStorage) ensureDirectories() {
	ls.ensureMainDir()
	ls.ensureStorageDir()
}

func (ls LocalStorage) ensureMainDir() error {
	mainDir, err := ls.getMainDir()
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

func (ls LocalStorage) getData() structs.Data {
	data := structs.Data{}

	_, err := os.Stat(ls.storageFile)
	if err == nil {
		file, _ := ioutil.ReadFile(ls.storageFile)
		_ = json.Unmarshal([]byte(file), &data)
	}

	return data
}

func (ls LocalStorage) setData(data structs.Data) {
	_, err := os.Stat(ls.storageDir)
	if err == nil {
		file, _ := json.MarshalIndent(data, "", " ")
		_ = ioutil.WriteFile(ls.storageFile, file, 0644)
	}
}
