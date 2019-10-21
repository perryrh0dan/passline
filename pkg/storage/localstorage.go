package storage

import (
	"encoding/json"
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
func Get(website string) (Website, error) {
	data := getData()
	for i := 0; i < len(data.Websites); i++ {
		if data.Websites[i].Domain == website {
			return data.Websites[i], nil
		}
	}

	return Website{}, fmt.Errorf("No entry for website %s", website)
}

func GetAll() ([]Website, error) {
	data := getData()
	return data.Websites, nil
}

// Add data
func Add(website Website) error {
	data := getData()
	data.Websites = append(data.Websites, website)
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
