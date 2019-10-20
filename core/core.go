package core

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"

	"github.com/perryrh0dan/passline/crypt"
	"github.com/perryrh0dan/passline/renderer"
	"github.com/perryrh0dan/passline/storage"
)

func getPassword() []byte {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Global Password: ")

	// ask for global password
	text, _ := reader.ReadString('\n')
	return []byte(text)
}

// DisplayBySite the password
func DisplayBySite(args []string) error {
	if len(args) == 1 {
		item, err := storage.Get(args[0])
		if err != nil {
			renderer.InvalidWebsite(args[0])
			return err
		}

		globalPassword := getPassword()

		item.Password, err = crypt.Decrypt(globalPassword, item.Password)
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

		globalPassword := getPassword()

		var err error
		website.Password, err = crypt.Encrypt(globalPassword, website.Password)
		if err != nil {
			return err
		}

		storage.Add(website)
		return nil
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
