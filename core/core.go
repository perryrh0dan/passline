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

	"github.com/perryrh0dan/passline/crypt"
	"github.com/perryrh0dan/passline/renderer"
	"github.com/perryrh0dan/passline/storage"
)

func getPassword(c *cli.Context) ([]byte, error) {
	password := c.String("password")

	if password == "" {
		// Get global password
		fmt.Print("Enter Global Password: ")

		// Ask for global password
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}

		fmt.Println()
		return bytePassword, nil
	} else {
		return []byte(password), nil
	}
}

// DisplayBySite the password
func DisplayBySite(c *cli.Context) error {
	args := c.Args()

	if len(args) == 1 {
		// Check if entry for website exists
		item, err := storage.Get(args[0])
		if err != nil {
			renderer.InvalidWebsite(args[0])
			return err
		}

		globalPassword, err := getPassword(c)
		if err != nil {
			return err
		}
		// Generate key from password with kdf
		key := crypt.GenerateKey(globalPassword)

		// Decrypt password
		item.Password, err = crypt.AesGcmDecrypt(key, item.Password, item.Nonce)
		if err != nil {
			return err
		}

		// Display item and copy password to clipboard
		renderer.DisplayWebsite(item)
		err = clipboard.WriteAll(item.Password)
		if err != nil {
			renderer.ClipboardError()
			return err
		}

		return nil
	} else {
		if len(args) < 1 {
			renderer.MissingArgument([]string{"domain"})
			return errors.New("Not enough arguments")
		} else {
			return errors.New("Too many arguments")
		}
	}
}

// Generate a random password for a website
func GenerateForSite(c *cli.Context) error {
	args := c.Args()

	if len(args) == 2 {
		// Check if entry exists
		_, err := storage.Get(args[0])
		if err == nil {
			return nil
		}

		// Generate new website entry
		website := storage.Website{Domain: args[0], Username: args[1], Password: generatePassword(20)}
		renderer.DisplayWebsite(website)

		globalPassword, err := getPassword(c)
		if err != nil {
			return err
		}

		key := crypt.GenerateKey(globalPassword)
		website.Password, website.Nonce, err = crypt.AesGcmEncrypt(key, website.Password)
		if err != nil {
			return err
		}

		err = storage.Add(website)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func ListSites() error {
	websites, err := storage.GetAll()
	if err != nil {
		return nil
	}

	renderer.DisplayWebsites(websites)
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
