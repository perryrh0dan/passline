package renderer

import (
	"fmt"

	"github.com/perryrh0dan/passline/storage"
)

// DisplayItem single item
func DisplayItem(item storage.Website) {
	fmt.Printf("Website: %s, Password: %s\n", item.Domain, item.Password)
}

// InvalidWebsite error message
func InvalidWebsite(website string) {
	fmt.Printf("Unable to find password for website: %s\n", website)
}
