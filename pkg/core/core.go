package core

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	ucli "github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/perryrh0dan/passline/pkg/config"
	"github.com/perryrh0dan/passline/pkg/crypt"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
	"github.com/perryrh0dan/passline/pkg/ui"
)

var store storage.Storage

func init() {
	conf, _ := config.Get()
	switch conf.Storage {
	case "firestore":
		store = &storage.FireStore{}
	default:
		store = &storage.LocalStorage{}
	}
	err := store.Init()
	if err != nil {
		renderer.StorageError()
		os.Exit(1)
	}
}

func getPassword(c *ucli.Context) ([]byte, error) {
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
	if err != nil || !valid {
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

	_, err = crypt.AesGcmDecrypt(password, item.Credentials[0].Password)
	if err != nil {
		renderer.InvalidPassword()
		return false, err
	}

	return true, nil
}

// DisplayByName the password
func DisplayByName(c *ucli.Context) error {
	args := c.Args()
	renderer.DisplayMessage()

	name, err := ui.GetArgOrInput(args, 0, "Please enter the URL []: ", nil)
	if err != nil {
		os.Exit(1)
	}

	// Check if item for name exists
	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		return nil
	}

	var credential storage.Credential
	// Only need username for multiple credentials
	if len(item.Credentials) > 1 {
		// User input username

		conf, err := config.Get()
		if err != nil {
			os.Exit(1)
		}

		var username string

		if conf.Selection {
			username, err = ui.GetArgOrSelect(args, 1, "Please select a Username/Login:", item.GetUsernameArray())
			if err != nil {
				os.Exit(1)
			}
		} else {
			username, err = ui.GetArgOrInput(args, 1, "Please enter the Username/Login []: ", nil)
			if err != nil {
				os.Exit(1)
			}
		}

		// Check if name, username combination exists
		item, err := store.GetByName(name)
		if err == nil {
			credential, err = item.GetCredentialByUsername(username)
			if err != nil {
				renderer.InvalidUsername(name, username)
				os.Exit(0)
			}
		}
	} else {
		credential = item.Credentials[0]
	}

	// Get and Check for global password.
	globalPassword, err := getPassword(c)
	if err != nil {
		return nil
	}

	// Decrypt passwords
	credential.Password, err = crypt.AesGcmDecrypt(globalPassword, credential.Password)
	if err != nil {
		os.Exit(0)
	}

	// Display item and copy password to clipboard
	renderer.DisplayCredential(credential)
	err = clipboard.WriteAll(credential.Password)
	if err != nil {
		renderer.ClipboardError()
		return nil
	}

	renderer.SuccessfulCopiedToClipboard(item.Name, credential.Username)
	return nil
}

// Generate a random password for a item
func GenerateForSite(c *ucli.Context) error {
	args := c.Args()
	renderer.CreateMessage()

	// User input name
	name, err := ui.GetArgOrInput(args, 0, "Please enter the URL []: ", nil)
	if err != nil {
		os.Exit(1)
	}

	// User input username
	username, err := ui.GetArgOrInput(args, 1, "Please enter the Username/Login []: ", nil)
	if err != nil {
		os.Exit(1)
	}

	// Check if name, username combination exists
	item, err := store.GetByName(name)
	if err == nil {
		_, err = item.GetCredentialByUsername(username)
		if err == nil {
			return nil
		}
	}

	// Get and Check for global password.
	globalPassword, err := getPassword(c)
	if err != nil {
		return nil
	}

	// Generate password and crypt password
	password := generatePassword(20)
	cryptedPassword, err := crypt.AesGcmEncrypt(globalPassword, password)

	// Create Credentials
	credential := storage.Credential{Username: username, Password: cryptedPassword}

	// Check if item already exists
	item, err = store.GetByName(name)
	if err != nil {
		// Generate new item entry
		item := storage.Item{Name: name, Credentials: []storage.Credential{credential}}
		err = store.AddItem(item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// Add to existing item
		err := store.AddCredential(name, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	err = clipboard.WriteAll(password)
	if err != nil {
		renderer.ClipboardError()
		os.Exit(0)
	}

	renderer.SuccessfulCopiedToClipboard(name, username)
	return nil
}

func DeleteItem(c *ucli.Context) error {
	args := c.Args()

	name, err := ui.GetArgOrInput(args, 0, "Please enter the URL []: ", nil)
	if err != nil {
		os.Exit(1)
	}

	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		os.Exit(0)
	}

	if len(item.Credentials) <= 1 {
		// If Item has only one Credential delete item
		err = store.DeleteItem(item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// If Item has multiple Credentials ask for username
		// User input username
		username, err := ui.GetArgOrInput(args, 1, "Please enter the Username/Login []: ", nil)
		if err != nil {
			os.Exit(1)
		}

		// Check if name, username combination exists
		var credential storage.Credential
		credential, err = item.GetCredentialByUsername(username)
		if err != nil {
			renderer.InvalidUsername(name, username)
			os.Exit(0)
		}

		err = store.DeleteCredential(item, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	return nil
}

func EditItem(c *ucli.Context) error {
	args := c.Args()
	renderer.ChangeMessage()

	// User input name
	name, err := ui.GetArgOrInput(args, 0, "Please enter the URL []: ", nil)
	if err != nil {
		os.Exit(1)
	}

	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		os.Exit(0)
	}

	username := ""
	if len(item.Credentials) == 1 {
		username = item.Credentials[0].Username
	} else {
		// User input username
		username, err := ui.GetArgOrInput(args, 1, "Please enter the Username/Login []: ", nil)
		if err != nil {
			os.Exit(1)
		}

		// Check if name, username combination exists
		_, err = item.GetCredentialByUsername(username)
		if err != nil {
			renderer.InvalidUsername(name, username)
			os.Exit(0)
		}
	}

	// Get new username
	newUsername := ui.GetInput("Please enter a new Username/Login []: (%s) ", []string{username})
	if newUsername == "" {
		newUsername = username
	}

	// Get new password
	// newPassword := getInput("Please enter a new Password []: (********)")

	for i := 0; i < len(item.Credentials); i++ {
		if item.Credentials[i].Username == username {
			item.Credentials[i].Username = newUsername
		}
	}

	err = store.UpdateItem(item)
	if err != nil {
		os.Exit(1)
	}

	renderer.SuccessfulChangedItem()

	return nil
}

func ListSites(c *ucli.Context) error {
	args := c.Args()

	if len(args) >= 1 {
		item, err := store.GetByName(args[0])
		if err != nil {
			renderer.InvalidName(args[0])
			os.Exit(0)
		}

		renderer.DisplayItem(item)
	} else {
		items, err := store.GetAll()
		if err != nil {
			return nil
		}

		if len(items) == 0 {
			renderer.NoItemsMessage()
		}
		renderer.DisplayItems(items)
	}

	return nil
}

func generatePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!$%&()/?")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	password := b.String() // E.g. "ExcbsVQs"
	return password
}
