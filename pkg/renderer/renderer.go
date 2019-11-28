package renderer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/perryrh0dan/passline/pkg/storage"
)

// DisplayItem single item
func DisplayItem(item storage.Item) {
	for i := 0; i < len(item.Credentials); i++ {
		fmt.Printf("%s\n", item.Credentials[i].Username)
	}
}

func DisplayCredential(credential storage.Credential) {
	fmt.Printf("Username: %s\n", credential.Username)
	fmt.Printf("Password: %s\n", credential.Password)
}

func DisplayItems(websites []storage.Item) {
	for _, website := range websites {
		fmt.Printf("%s\n", website.Name)
	}
}

func SuccessfulCopiedToClipboard(name string, username string) {
	name = color.YellowString(name + "/" + username)
	fmt.Fprintf(color.Output, "Copied Password for %s to clipboard\n", name)
}

func SuccessfulChangedItem() {
	d := color.New(color.FgGreen)
	d.Printf("Successful changed item")
}

// InvalidName error message
func InvalidName(name string) {
	fmt.Printf("Unable to find item with name: %s\n", name)
}

func InvalidUsername(name string, username string) {
	fmt.Printf("Unable to find username: %s in item: %s\n", username, name)
}

func InvalidPassword() {
	fmt.Printf("Invalid Password\n")
}

func ClipboardError() {
	fmt.Printf("Error occured while copying to clipboard\n")
}

func StorageError() {
	d := color.New(color.FgRed)
	d.Printf("error: unable to initialice storage\n")
}

func NoItemsMessage() {
	d := color.New(color.FgYellow)
	d.Printf("No items yet\n")
}

func DisplayMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Display item...\n")
}

func CreateMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Creating item...\n")
}

func GenerateMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Generating item...\n")
}

func ChangeMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Changing item...\n")
}

func DeleteMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Deleting item...\n")
}

func MissingArgument(arguments []string) {
	d := color.New(color.FgRed)
	d.Printf("error: missing required arguments %s\n", strings.Join(arguments, ", "))
}

func NameAlreadyExists(name string) {
	fmt.Printf("error: name already exists %s\n", name)
}
