package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Config struct {
	Directory string
	Storage   string
	Selection bool
	AutoClip  bool
	NoColor   bool
	NoSymbols bool
}

var configFile string

func init() {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configFile = path.Join(homeDir, ".passline.json")
	}

	ensureConfigFile()
}

func ensureConfigFile() {
	_, err := os.Stat(configFile)
	if err == nil {
		return
	}

	config := new()
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(configFile, file, 0644)
}

func formatPasslineDir(dirPath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, strings.Replace(dirPath, "~", "", 1), ".passline"), nil
}

func new() Config {
	return Config{
		Directory: "~",
		Storage:   "local",
		Selection: true,
		AutoClip:  true,
		NoColor:   false,
		NoSymbols: false,
	}
}

func Get() (*Config, error) {
	config := Config{}

	file, _ := ioutil.ReadFile(configFile)
	_ = json.Unmarshal([]byte(file), &config)

	if strings.HasPrefix(config.Directory, "~") {
		var err error
		config.Directory, err = formatPasslineDir(config.Directory)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
