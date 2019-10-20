package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var storageDir string
var storageFile string

func init() {
	storageDir = path.Join(getMainDir(), "storage")
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

// Add data
func Add(website Website) error {
	data := getData()
	data.Websites = append(data.Websites, website)
	setData(data)
	return nil
}

func getMainDir() string {
	dir := path.Join("~", ".passline")
	return dir
}

func ensureDirectories() {
	ensureMainDir()
	ensureStorageDir()
}

func ensureMainDir() {
	_, err := os.Stat(getMainDir())
	if err != nil {
		err := os.MkdirAll(getMainDir(), os.ModePerm)
		if err != nil {
			println("Cant create directory")
		}
	}
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
