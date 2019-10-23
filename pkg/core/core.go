package core

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/perryrh0dan/passline/pkg/config"
	"github.com/perryrh0dan/passline/pkg/crypt"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
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
	// TODO Test if working for Linux
	text = strings.TrimSuffix(text, "\r")
	return text
}

// DisplayByName the password
func DisplayByName(c *cli.Context) error {
	args := c.Args()

	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if name == "" {
		fmt.Printf("Please enter the URL []: ")
		name = getInput()
	}

	// Check if item for name exists
	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		return err
	}

	// Get and Check for global password.
	globalPassword, err := getPassword(c)
	if err != nil {
		return err
	}

	// Decrypt password
	item.Password, err = crypt.AesGcmDecrypt(globalPassword, item.Password)
	if err != nil {
		return err
	}

	// Display item and copy password to clipboard
	renderer.DisplayItem(item)
	err = clipboard.WriteAll(item.Password)
	if err != nil {
		renderer.ClipboardError()
		return err
	}

	renderer.SuccessfullCopiedToClipboard(item.Name)
	return nil
}

// Generate a random password for a item
func GenerateForSite(c *cli.Context) error {
	args := c.Args()
	renderer.CreatingMessage()

	// User input name
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
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

	// User input username
	username := ""
	if len(args) >= 2 {
		username = args[1]
	}
	if username == "" {
		fmt.Printf("Please enter the Username/Login []: ")
		username = getInput()
	}

	// Get and Check for global password.
	globalPassword, err := getPassword(c)
	if err != nil {
		return err
	}

	// Generate new item entry
	password := generatePassword(20)
	item := storage.Item{Name: name, Username: username, Password: password}

	item.Password, err = crypt.AesGcmEncrypt(globalPassword, item.Password)
	if err != nil {
		return err
	}

	err = store.Add(item)
	if err != nil {
		return err
	}
	err = clipboard.WriteAll(password)
	if err != nil {
		renderer.ClipboardError()
		return err
	}

	renderer.SuccessfullCopiedToClipboard(item.Name)
	return nil
}

func DeleteItem(c *cli.Context) error {
	args := c.Args()

	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if name == "" {
		fmt.Printf("Please enter the URL []: ")
		name = getInput()
	}

	item, err := store.GetByName(name)
	if err != nil {
		renderer.InvalidName(name)
		return err
	}

	err = store.Delete(item)
	if err != nil {
		return err
	}

	return nil
}

func ListSites() error {
	websites, err := store.GetAll()
	if err != nil {
		return nil
	}

	renderer.DisplayItems(websites)
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
