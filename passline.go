package main

import (
	"bufio"
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

func displayBySite(args []string) error {
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
		renderer.DisplayItem(item)
		clipboard.WriteAll(item.Password)
	}

	return nil
}

// Generate a random password
func generateForSite(args []string) error {
	if len(args) == 2 {
		domain := args[0]
		username := args[1]

		globalPassword := getPassword()

		password, err := crypt.Encrypt(globalPassword, generatePassword(20))
		if err != nil {
			return err
		}
		storage.Add(storage.Website{Domain: domain, Username: username, Password: password})
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
