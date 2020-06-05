package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrConfigNotFound is returned on load if the config was not found
	ErrConfigNotFound = errors.Errorf("config not found")
	// ErrConfigNotParsed is returned on load if the config could not be decoded
	ErrConfigNotParsed = errors.Errorf("config not parseable")
)

type Config struct {
	Directory     string `yaml:"directory"`
	Storage       string `yaml:"storage"`
	AutoClip      bool   `yaml:"autoclip"`
	Notifications bool   `yaml:"notifications"`
	NoColor       bool   `yaml:"nocolor"`
}

var configFile string

func init() {
	configFile = configLocation()

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

// configLocation returns the location of the config file
// (a YAML file that contains values such as the path to the password store)
func configLocation() string {
	// First, check for the "PASSLINE_CONFIG" environment variable
	if cf := os.Getenv("PASSLINE_CONFIG"); cf != "" {
		return cf
	}

	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, ".passline", "config.json")
}

// Directory returns the configuration directory for the gopass config file
func Directory() string {
	return filepath.Dir(configLocation())
}
