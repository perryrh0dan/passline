package core

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"

	"github.com/perryrh0dan/passline/crypt"
	"github.com/perryrh0dan/passline/renderer"
	"github.com/perryrh0dan/passline/storage"

	"golang.org/x/crypto/ssh/terminal"
)

func getPassword() ([]byte, error) {
	fmt.Print("Enter Global Password: ")

	// ask for global password
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println()
	return bytePassword, nil
}

// DisplayBySite the password
func DisplayBySite(args []string) error {
	if len(args) == 1 {
		item, err := storage.Get(args[0])
		if err != nil {
			renderer.InvalidWebsite(args[0])
			return err
		}

		globalPassword, err := getPassword()
		if err != nil {
			return err
		}

		key := crypt.GenerateKey(globalPassword)
		item.Password, err = crypt.AesGcmDecrypt(key, item.Password, item.Nonce)
		if err != nil {
			return err
		}
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
func GenerateForSite(args []string) error {
	if len(args) == 2 {
		domain := args[0]
		username := args[1]
		password := generatePassword(20)

		website := storage.Website{Domain: domain, Username: username, Password: password}
		renderer.DisplayWebsite(website)

		globalPassword, err := getPassword()
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
