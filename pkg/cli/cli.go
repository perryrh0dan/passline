package cli

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	ucli "github.com/urfave/cli"

	"github.com/perryrh0dan/passline/pkg/core"
	"github.com/perryrh0dan/passline/pkg/renderer"
	"github.com/perryrh0dan/passline/pkg/storage"
)

var passline *core.Passline

func init() {
	passline = core.NewPassline()
}

func AddItem(c *ucli.Context) error {
	args := c.Args()
	renderer.CreateMessage()

	// User input name
	name, err := ArgOrInput(args, 0, "URL", "")
	handle(err)

	// User input username
	username, err := ArgOrInput(args, 1, "Username/Login", "")
	handle(err)

	// Check if name, username combination exists
	exists, err := passline.Exists(name, username)
	if exists {
		os.Exit(0)
	}

	password, err := Input("Password", "")
	handle(err)

	globalPassword := GetPassword("Enter Global Password: ")

	credential, err := passline.AddItem(name, username, password, globalPassword)

	renderer.DisplayCredential(credential)

	return nil
}

func DeleteItem(c *ucli.Context) error {
	// Get all Sites
	names, err := passline.GetSiteNames()

	// Check if site exists
	if len(names) <= 0 {
		renderer.NoItemsMessage()
		os.Exit(0)
	}

	args := c.Args()
	renderer.DeleteMessage()

	item, err := selectItem(args, names)
	handle(err)

	credential, err := selectCredential(args, item)
	handle(err)

	err = passline.DeleteItem(item.Name, credential.Username)
	handle(err)

	return nil
}

func DisplayItem(c *ucli.Context) error {
	// Get all Sites
	names, err := passline.GetSiteNames()

	// Check if site exists
	if len(names) <= 0 {
		renderer.NoItemsMessage()
		os.Exit(0)
	}

	args := c.Args()
	renderer.DisplayMessage()

	item, err := selectItem(args, names)
	handle(err)

	credential, err := selectCredential(args, item)
	handle(err)

	// Get global password.
	globalPassword := GetPassword("Enter Global Password: ")

	// Check global password.
	valid, err := passline.CheckPassword(globalPassword)
	if err != nil || !valid {
		handle(err)
	}

	err = passline.DecryptPassword(&credential, globalPassword)
	handle(err)

	renderer.DisplayCredential(credential)

	return nil
}

func EditItem(c *ucli.Context) error {
	// Get all Sites
	names, err := passline.GetSiteNames()

	// Check if site exists
	if len(names) <= 0 {
		renderer.NoItemsMessage()
		os.Exit(0)
	}

	args := c.Args()
	renderer.DeleteMessage()

	item, err := selectItem(args, names)
	handle(err)

	credential, err := selectCredential(args, item)
	handle(err)

	// Get new username
	newUsername, err := Input("Please enter a new Username/Login []: (%s) ", credential.Username)
	handle(err)

	if newUsername == "" {
		newUsername = credential.Username
	}

	err = passline.EditItem(item.Name, credential.Username, newUsername)
	handle(err)

	renderer.SuccessfulChangedItem()

	return nil
}

func GenerateItem(c *ucli.Context) error {
	args := c.Args()
	renderer.GenerateMessage()

	// User input name
	name, err := ArgOrInput(args, 0, "URL", "")
	handle(err)

	// User input username
	username, err := ArgOrInput(args, 1, "Username/Login", "")
	handle(err)

	// Check if name, username combination exists
	exists, err := passline.Exists(name, username)
	if exists {
		os.Exit(0)
	}

	globalPassword := GetPassword("Enter Global Password: ")

	credential, err := passline.GenerateItem(name, username, globalPassword)
	handle(err)

	err = clipboard.WriteAll(credential.Password)
	if err != nil {
		renderer.ClipboardError()
		os.Exit(0)
	}

	renderer.SuccessfulCopiedToClipboard(name, credential.Username)
	return nil
}

func ListItems(c *ucli.Context) error {
	args := c.Args()

	if len(args) >= 1 {
		item, err := passline.GetSite(args[0])
		if err != nil {
			renderer.InvalidName(args[0])
			os.Exit(0)
		}

		renderer.DisplayItem(item)
	} else {
		items, err := passline.GetSites()
		if err != nil {
			return nil
		}

		if len(items) == 0 {
			renderer.NoItemsMessage()
		}
		renderer.DisplayItems(items)
	}

	return nil
}

func ArgOrInput(args []string, index int, message string, value string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		message := fmt.Sprintf("Please enter a %s []: ", message)
		var err error
		input, err = Input(message, value)
		if err != nil {
			return "", err
		}
	}

	return input, nil
}

func ArgOrSelect(args []string, index int, message string, items []string) (string, error) {
	input := ""
	if len(args)-1 >= index {
		input = args[index]
	}
	if input == "" {
		if len(items) > 1 {
			message := fmt.Sprintf("Please select a %s: ", message)
			selection, err := Select(message, items)
			if err != nil {
				return "", err
			}

			input = items[selection]
			fmt.Printf("%s%s\n", message, input)
		} else if len(items) == 1 {
			fmt.Printf("Selected %s: %s\n", message, items[0])
			return items[0], nil
		}
	}

	return input, nil
}

func handle(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func selectItem(args, names []string) (storage.Item, error) {
	name, err := ArgOrSelect(args, 0, "URL", names)
	handle(err)

	// Get item
	item, err := passline.GetSite(name)
	if err != nil {
		os.Exit(0)
	}

	return item, nil
}

func selectCredential(args []string, item storage.Item) (storage.Credential, error) {
	username, err := ArgOrSelect(args, 1, "Username/Login", item.GetUsernameArray())
	handle(err)

	// Check if name, username combination exists
	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		renderer.InvalidUsername(item.Name, username)
		os.Exit(0)
	}

	return credential, nil
}
