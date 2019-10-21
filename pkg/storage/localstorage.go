package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/perryrh0dan/passline/pkg/config"
)

var storageDir string
var storageFile string

func init() {
	mainDir, _ := getMainDir()

	storageDir = path.Join(mainDir, "storage")
	storageFile = path.Join(storageDir, "storage.json")

	ensureDirectories()
}

// Get data
func GetByName(name string) (Item, error) {
	data := getData()
	for i := 0; i < len(data.Items); i++ {
		if data.Items[i].Name == name {
			return data.Items[i], nil
		}
	}

	return Item{}, fmt.Errorf("No entry for website %s", name)
}

func GetByindex(index int) (Item, error) {
	data := getData()
	if 0 <= index && index < len(data.Items) {
		return data.Items[index], nil
	} else {
		return Item{}, errors.New("Out of index")
	}
}

func GetAll() ([]Item, error) {
	data := getData()
	return data.Items, nil
}

// Add data
func Add(website Item) error {
	data := getData()
	data.Items = append(data.Items, website)
	setData(data)
	return nil
}

func getMainDir() (string, error) {
	config, err := config.Get()
	if err != nil {
		return "", err
	}

	return config.Directory, nil
}

func ensureDirectories() {
	ensureMainDir()
	ensureStorageDir()
}

func ensureMainDir() error {
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

func ensureStorageDir() {
	_, err := os.Stat(storageDir)
	if err != nil {
		err := os.Mkdir(storageDir, os.ModePerm)
		if err != nil {
			println("Cant create directory")
		}
	}
}

func getData() Data {
	data := Data{}

	_, err := os.Stat(storageFile)
	if err == nil {
		file, _ := ioutil.ReadFile(storageFile)
		_ = json.Unmarshal([]byte(file), &data)
	}

	return data
}

func setData(data Data) {
	_, err := os.Stat(storageDir)
	if err == nil {
		file, _ := json.MarshalIndent(data, "", " ")
		_ = ioutil.WriteFile(storageFile, file, 0644)
	}
}
