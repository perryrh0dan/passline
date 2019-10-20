package renderer

import (
	"fmt"
	"strings"

	"github.com/perryrh0dan/passline/storage"
)

// DisplayWebsite single item
func DisplayWebsite(item storage.Website) {
	fmt.Printf("Website: %s, Password: %s\n", item.Domain, item.Password)
}

// InvalidWebsite error message
func InvalidWebsite(website string) {
	fmt.Printf("Unable to find password for website: %s\n", website)
}

func ClipboardError() {
	fmt.Printf("Error occured while copying to clipboard\n")
}

func MissingArgument(arguments []string) {
	fmt.Printf("error: missing required arguments %s\n", strings.Join(arguments, ", "))
}
