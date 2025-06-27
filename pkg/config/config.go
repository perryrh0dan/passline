package config

import (
	"encoding/json"
	"os"
	"passline/pkg/util"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	// ErrConfigNotFound is returned on load if the config was not found
	ErrConfigNotFound = errors.Errorf("config not found")
	// ErrConfigNotParsed is returned on load if the config could not be decoded
	ErrConfigNotParsed = errors.Errorf("config not parseable")
)

type Config struct {
	Storage         string `yaml:"storage"`
	Encryption      int    `yaml:"encryption"`
	AutoClip        bool   `yaml:"autoclip"`
	Notifications   bool   `yaml:"notifications"`
	QuickSelect     bool   `yaml:"quickselect"`
	DefaultUsername string `yaml:"defaultUsername"`
	DefaultCategory string `yaml:"defaultCategory"`
	PhoneNumber     string `yaml:"phoneNumber"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if c.DefaultCategory == "" {
		c.DefaultCategory = "*"
	}

	return nil
}

func init() {
	ensureConfigFile()
	_ = ensureMainDir()
	_ = ensureBackupDir()
}

func ensureConfigFile() {
	_, err := os.Stat(configLocation())
	if err == nil {
		return
	}

	config := new()
	file, _ := json.MarshalIndent(config, "", " ")
	_ = os.WriteFile(configLocation(), file, 0644)
}

func ensureMainDir() error {
	mainDir := Directory()
	_, err := os.Stat(mainDir)
	if err != nil {
		err := os.MkdirAll(mainDir, os.ModePerm)
		if err != nil {
			println("Cant create directory")
		}
	}

	return nil
}

func ensureBackupDir() error {
	backupDir := Directory() + "/backup"
	_, err := os.Stat(backupDir)
	if err != nil {
		err := os.Mkdir(backupDir, os.ModePerm)
		if err != nil {
			print("Cant create backup directory")
		}
	}

	return nil
}

const (
	PartialEncryption = iota
	FullEncryption    = iota
)

func new() Config {
	return Config{
		Storage:       "local",
		AutoClip:      true,
		Notifications: true,
		Encryption:    PartialEncryption,
	}
}

func Get(fs util.FileSystem) (*Config, error) {
	config := new()

	file, _ := fs.ReadFile(configLocation())
	_ = json.Unmarshal([]byte(file), &config)

	return &config, nil
}

// configLocation returns the location of the config file
// (a JSON file that contains values such as the path to the password store)
func configLocation() string {
	// First, check for the "PASSLINE_CONFIG" environment variable
	if cf := os.Getenv("PASSLINE_CONFIG"); cf != "" {
		return cf
	}

	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, ".passline", "config.json")
}

// Directory returns the configuration directory for the passline config file
func Directory() string {
	return filepath.Dir(configLocation())
}

func BackupDirectory() string {
	return Directory() + "/backup"
}
