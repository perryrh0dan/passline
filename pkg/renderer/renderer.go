package renderer

import (
	"fmt"
	"strings"

	"github.com/perryrh0dan/passline/pkg/storage"
)

// DisplayItem single item
func DisplayItem(item storage.Item) {
	fmt.Printf("Name: %s\nUsername: %s\nPassword: %s\n", item.Name, item.Username, item.Password)
}

func DisplayItems(websites []storage.Item) {
	for _, website := range websites {
		fmt.Printf("%s\n", website.Name)
	}
}

// InvalidName error message
func InvalidName(name string) {
	fmt.Printf("Unable to find item with name: %s\n", name)
}

func ClipboardError() {
	fmt.Printf("Error occured while copying to clipboard\n")
}

func MissingArgument(arguments []string) {
	fmt.Printf("error: missing required arguments %s\n", strings.Join(arguments, ", "))
}
