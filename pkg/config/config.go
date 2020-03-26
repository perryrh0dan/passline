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
	config, _ := Get()
	_ = ensureMainDir(config)
	_ = ensureBackupDir(config)
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

func ensureMainDir(config *Config) error {
	mainDir := config.Directory
	_, err := os.Stat(mainDir)
	if err != nil {
		err := os.MkdirAll(mainDir, os.ModePerm)
		if err != nil {
			println("Cant create directory")
		}
	}

	return nil
}

func ensureBackupDir(config *Config) error {
	backupDir := config.Directory + "/backup"
	_, err := os.Stat(backupDir)
	if err != nil {
		err := os.Mkdir(backupDir, os.ModePerm)
		if err != nil {
			print("Cant create backup directory")
		}
	}

	return nil
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
