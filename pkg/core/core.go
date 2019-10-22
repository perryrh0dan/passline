package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/perryrh0dan/passline/pkg/config"
	"github.com/perryrh0dan/passline/pkg/crypt"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
	"github.com/perryrh0dan/passline/pkg/structs"
	"github.com/perryrh0dan/passline/pkg/utils"
)

var store storage.Storage

func init() {
	conf, _ := config.Get()
	switch conf.Storage {
	case "local":
		store = &storage.LocalStorage{}
	case "firestore":
		store = &storage.FireStore{}
	}
	err := store.Init()
	if err != nil {
		renderer.StorageError()
	}
}

func getPassword(c *cli.Context) ([]byte, error) {
	password := []byte(c.String("password"))

	if len(password) <= 0 {
		// Get global password
		fmt.Print("Enter Global Password: ")

		// Ask for global password
		var err error
		password, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}

		fmt.Println()
	}

	valid, err := checkPassword(password)
	if err != nil || valid == false {
		return nil, errors.New("Invalid password")
	}

	return password, nil
}

func checkPassword(password []byte) (bool, error) {
	data, err := store.GetAll()
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	item, err := store.GetByIndex(0)
	if err != nil {
		return false, err
	}

	_, err = crypt.AesGcmDecrypt(password, item.Password)
	if err != nil {
		renderer.InvalidPassword()
		return false, err
	}

	return true, nil
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	return text
}

// DisplayBySite the password
func DisplayByName(c *cli.Context) error {

	// Check for name flag or ask for input
	name := c.String("name")
	if name == "" {
		fmt.Printf("Please enter the URL []: ")
		name = getInput()
	}

	// Check if entry/name exists
	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		return err
	}

	globalPassword, err := getPassword(c)
	if err != nil {
		return err
	}

	// Decrypt password
	item.Password, err = crypt.AesGcmDecrypt(globalPassword, item.Password)
	if err != nil {
		return err
	}

	err = clipboard.WriteAll(item.Password)
	if err != nil {
		renderer.ClipboardError()
		return err
	}

	return nil
}

// Generate a random password for a item
func GenerateForSite(c *cli.Context) error {
	renderer.CreatingMessage()

	// Check for name flag or ask for input
	name := c.String("name")
	if name == "" {
		fmt.Printf("Please enter the URL []: ")
		name = getInput()
	}

	// Check if entry/name already exists
	_, err := store.GetByName(name)
	if err == nil {
		renderer.NameAlreadyExists(name)
		return err
	}

	// Check for username flag or aks for input
	username := c.String("username")
	if username == "" {
		fmt.Printf("Please enter the Username/Login []: ")
		username = getInput()
	}

	// Get global password.
	globalPassword, err := getPassword(c)
	if err != nil {
		return err
	}

	// Generate new item
	item := structs.Item{Name: name, Username: username, Password: utils.GeneratePassword(20)}

	item.Password, err = crypt.AesGcmEncrypt(globalPassword, item.Password)
	if err != nil {
		return err
	}

	err = store.Add(item)
	if err != nil {
		return err
	}

	coloredName := color.YellowString(name)
	fmt.Fprintf(color.Output, "Copied Password for %s to clipboard\n", coloredName)

	return nil
}

func DeleteByName(c *cli.Context) error {
	// Check for name flag or ask for input
	name := c.String("name")
	if name == "" {
		fmt.Printf("Please enter the URL []: ")
		name = getInput()
	}

	// Check if entry/name already exists
	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		return err
	}

	// Delete item
	err = store.Delete(item)
	if err != nil {
		return err
	}

	return nil
}

func ListAllItems() error {
	websites, err := store.GetAll()
	if err != nil {
		return nil
	}

	renderer.DisplayItems(websites)
	return nil
}
