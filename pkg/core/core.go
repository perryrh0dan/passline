package core

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/perryrh0dan/passline/pkg/crypt"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
	"github.com/perryrh0dan/passline/pkg/structs"
)

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
	data, err := storage.GetAll()
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	item, err := storage.GetByindex(0)
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

// DisplayBySite the password
func DisplayBySite(c *cli.Context) error {
	args := c.Args()

	if len(args) == 1 {
		// Check if entry for website exists
		item, err := storage.GetByName(args[0])
		if err != nil {
			renderer.InvalidName(args[0])
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

		// Display item and copy password to clipboard
		renderer.DisplayItem(item)
		err = clipboard.WriteAll(item.Password)
		if err != nil {
			renderer.ClipboardError()
			return err
		}

		return nil
	} else {
		if len(args) < 1 {
			renderer.MissingArgument([]string{"name"})
			return errors.New("Not enough arguments")
		} else {
			return errors.New("Too many arguments")
		}
	}
}

// Generate a random password for a item
func GenerateForSite(c *cli.Context) error {
	args := c.Args()

	if len(args) == 2 {
		// Check if entry exists
		_, err := storage.GetByName(args[0])
		if err == nil {
			renderer.NameAlreadyExists(args[0])
			return err
		}

		// Check for password. No needed before generation but the flow is much better before
		globalPassword, err := getPassword(c)
		if err != nil {
			return err
		}

		// Generate new item entry
		item := structs.Item{Name: args[0], Username: args[1], Password: generatePassword(20)}
		renderer.DisplayItem(item)

		item.Password, err = crypt.AesGcmEncrypt(globalPassword, item.Password)
		if err != nil {
			return err
		}

		err = storage.Add(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteItem(c *cli.Context) error {
	args := c.Args()

	if len(args) == 1 {
		item, err := storage.GetByName(args[0])
		if err != nil {
			renderer.InvalidName(args[0])
			return err
		}

		err = storage.Delete(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func ListSites() error {
	websites, err := storage.GetAll()
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
