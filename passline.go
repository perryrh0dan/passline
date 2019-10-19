package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/atotto/clipboard"

	"github.com/perryrh0dan/passline/renderer"
	"github.com/perryrh0dan/passline/storage"
)

func checkPassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")

	// ask for global password
	text, _ := reader.ReadString('\n')
	return text
}

func displayBySite(args []string) {
	if len(args) == 1 {
		item, err := storage.Get(args[0])
		if err != nil {
			renderer.InvalidWebsite(args[0])
			return
		}
		renderer.DisplayItem(item)
		clipboard.WriteAll(item.Password)
	}
}

// Generate a random password
func generate() string {
	return "12345678"
}
