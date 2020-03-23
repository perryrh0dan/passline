package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	ucli "github.com/urfave/cli/v2"

	"passline/pkg/core"
	"passline/pkg/renderer"
	"passline/pkg/storage"
	"passline/pkg/util"
)

const repo = "perryrh0dan/passline"

var passline *core.Core

func Init(ctx context.Context) {
	var err error
	passline, err = core.NewCore(ctx)
	if err != nil {
		renderer.CoreInstanceError()
		os.Exit(1)
	}
}

func CreateBackup(ctx context.Context, c *ucli.Context) error {
	args := c.Args()
	renderer.BackupMessage()

	path := passline.Config.Directory

	now := time.Now().Format("2006-01-02:15:04:05.backup")
	path = filepath.Join(path, now)

	path, err := argOrInput(args, 0, "Path", path)
	if err != nil {
		return err
	}

	err = passline.CreateBackup(ctx, path)
	if err != nil {
		return err
	}

	renderer.SuccessfulCreatedBackup(path)
	return nil
}

func AddItem(ctx context.Context, c *ucli.Context) error {
	args := c.Args()
	renderer.CreateMessage()

	// User input name
	name, err := argOrInput(args, 0, "URL", "")
	handle(err)

	// User input username
	username, err := argOrInput(args, 1, "Username/Login", "")
	handle(err)

	// Check if name, username combination exists
	exists, err := passline.Exists(ctx, name, username)
	if err != nil {
		return err
	}

	if exists {
		renderer.NameUsernameAlreadyExists()
		return nil
	}

	password, err := Input("Please enter the existing Password []: ", "")
	if err != nil {
		return err
	}

	recoveryCodesString, err := Input("Please enter your recovery codes if exists []: ", "")
	if err != nil {
		return err
	}

	recoveryCodes := make([]string, 0)
	if recoveryCodesString != "" {
		recoveryCodes = util.StringToArray(recoveryCodesString)
	}

	globalPassword := getPassword("Enter Global Password: ")
	println()

	credential, err := passline.AddItem(ctx, name, username, password, recoveryCodes, globalPassword)
	if err != nil {
		return err
	}

	renderer.DisplayCredential(credential)
	return nil
}

func DeleteItem(ctx context.Context, c *ucli.Context) error {
	// Get all Sites
	names, err := passline.GetSiteNames(ctx)
	if err != nil {
		return err
	}

	// Check if site exists
	if len(names) <= 0 {
		renderer.NoItemsMessage()
		return nil
	}

	args := c.Args()
	renderer.DeleteMessage()

	item, err := selectItem(ctx, args, names)
	if err != nil {
		return err
	}

	credential, err := selectCredential(args, item)
	if err != nil {
		return err
	}

	err = passline.DeleteItem(ctx, item.Name, credential.Username)
	if err != nil {
		return err
	}

	renderer.SuccessfulDeletedItem(item.Name, credential.Username)
	return nil
}

func DisplayItem(ctx context.Context, c *ucli.Context) error {
	// Get all Sites
	names, err := passline.GetSiteNames(ctx)
	if err != nil {
		return err
	}

	// Check if sites exists
	if len(names) <= 0 {
		renderer.NoItemsMessage()
		os.Exit(0)
	}

	args := c.Args()
	renderer.DisplayMessage()

	item, err := selectItem(ctx, args, names)
	handle(err)

	credential, err := selectCredential(args, item)
	handle(err)

	// Get global password.
	globalPassword := getPassword("Enter Global Password: ")
	println()

	// Check global password.
	valid, err := passline.CheckPassword(ctx, globalPassword)
	if err != nil || !valid {
		handle(err)
	}

	// globalPassword := []byte("test")

	err = passline.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	renderer.DisplayCredential(credential)

	err = clipboard.WriteAll(credential.Password)
	if err != nil {
		renderer.ClipboardError()
		os.Exit(0)
	}

	renderer.SuccessfulCopiedToClipboard(item.Name, credential.Username)

	return nil
}

func EditItem(ctx context.Context, c *ucli.Context) error {
	// Get all Sites
	names, err := passline.GetSiteNames(ctx)
	if err != nil {
		return err
	}

	// Check if site exists
	if len(names) <= 0 {
		renderer.NoItemsMessage()
		os.Exit(0)
	}

	args := c.Args()
	renderer.DeleteMessage()

	item, err := selectItem(ctx, args, names)
	handle(err)

	credential, err := selectCredential(args, item)
	handle(err)

	selectedUsername := credential.Username

	// Get global password.
	globalPassword := getPassword("Enter Global Password: ")
	println()

	// Check global password.
	valid, err := passline.CheckPassword(ctx, globalPassword)
	if err != nil || !valid {
		handle(err)
	}

	// Decrypt Credentials to display secrets
	passline.DecryptCredential(&credential, globalPassword)
	if err != nil {
		handle(err)
	}

	// Get new username
	newUsername, err := Input("Please enter a new Username/Login []: (%s) ", credential.Username)
	handle(err)

	credential.Username = newUsername

	// Get new recoveryCodes
	recoveryCodes := util.ArrayToString(credential.RecoveryCodes)
	newRecoveryCodes, err := Input("Please enter your recovery codes []: (%s) ", recoveryCodes)
	handle(err)

	// TODO remove spaces
	credential.RecoveryCodes = make([]string, 0)

	// use one space to clear recovery codes
	if newRecoveryCodes != " " {
		credential.RecoveryCodes = util.StringToArray(newRecoveryCodes)
	}

	err = passline.EditItem(ctx, item.Name, selectedUsername, credential, globalPassword)
	handle(err)

	renderer.SuccessfulChangedItem(item.Name, credential.Username)

	return nil
}

func GenerateItem(ctx context.Context, c *ucli.Context) error {
	args := c.Args()
	renderer.GenerateMessage()

	// User input name
	name, err := argOrInput(args, 0, "URL", "")
	handle(err)

	// User input username
	username, err := argOrInput(args, 1, "Username/Login", "")
	handle(err)

	// Check if name, username combination exists
	exists, err := passline.Exists(ctx, name, username)
	if err != nil {
		return err
	}

	if exists {
		os.Exit(0)
	}

	// Check if advanced mode is active
	if c.String("mode") == "advanced" {
		err := getAdvancedParamters(ctx)
		if err != nil {
			return err
		}
	}

	recoveryCodesString, err := Input("Please enter your recovery codes if exists []: ", "")
	if err != nil {
		return err
	}

	recoveryCodes := make([]string, 0)
	if recoveryCodesString != "" {
		recoveryCodes = util.StringToArray(recoveryCodesString)
	}

	globalPassword := getPassword("Enter Global Password: ")
	println()

	credential, err := passline.GenerateItem(ctx, name, username, recoveryCodes, globalPassword)
	handle(err)

	err = clipboard.WriteAll(credential.Password)
	if err != nil {
		renderer.ClipboardError()
		os.Exit(0)
	}

	renderer.SuccessfulCopiedToClipboard(name, credential.Username)
	return nil
}

func ListItems(ctx context.Context, c *ucli.Context) error {
	args := c.Args()

	if args.Len() >= 1 {
		item, err := passline.GetSite(ctx, args.Get(0))
		if err != nil {
			renderer.InvalidName(args.Get(0))
			os.Exit(0)
		}

		renderer.DisplayItem(item)
	} else {
		items, err := passline.GetSites(ctx)
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

func RestoreBackup(ctx context.Context, c *ucli.Context) error {
	args := c.Args()
	renderer.RestoreMessage()

	// User input path
	path, err := argOrInput(args, 0, "Path", "")
	if err != nil {
		return err
	}

	// TODO Are you sure

	err = passline.RestoreBackup(ctx, path)
	if err != nil {
		return err
	}

	renderer.SuccessfulRestoredBackup(path)
	return nil
}

func Update(ctx context.Context, c *ucli.Context, version string) error {
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		renderer.DetectVersionError(err)
		return err
	}

	v := semver.MustParse(version)
	if !found || latest.Version.LTE(v) {
		renderer.NoUpdatesFound()
		return nil
	}

	message := "Do you want to update to: " + latest.Version.String() + "? (y/n): "
	input, err := Input(message, "n")
	if err != nil || (input != "y" && input != "n") {
		renderer.InvalidInput()
		return err
	}

	if input == "n" {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return err
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		renderer.UpdateError(err)
		return err
	}

	renderer.SuccessfulUpdated(latest.Version.String())
	renderer.DisplayReleaseNotes(latest.ReleaseNotes)

	return nil
}

func getAdvancedParamters(ctx context.Context) error {
	length, err := Input("Please enter the length of the password []: (%s)", "20")
	if err != nil {
		return err
	}

	fmt.Println(length)
	return nil
}

func argOrInput(args ucli.Args, index int, message string, defaultValue string) (string, error) {
	input := ""
	if args.Len()-1 >= index {
		input = args.Get(index)
	}
	if input == "" {
		message := fmt.Sprintf("Please enter a %s []: ", message)
		if defaultValue != "" {
			message += "(%s)"
		}

		var err error
		input, err = Input(message, defaultValue)
		if err != nil {
			return "", err
		}
	}

	return input, nil
}

func argOrSelect(args ucli.Args, index int, message string, items []string) (string, error) {
	input := ""
	if args.Len()-1 >= index {
		input = args.Get(index)

		// if input is no item name use as filter
		if !util.ArrayContains(items, input) {
			items = util.FilterArray(items, input)
			if len(items) == 0 {
				fmt.Printf("No items with filter: %v found\n", input)
				return "", errors.New("No items found")
			}
			input = ""
		}
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

func selectItem(ctx context.Context, args ucli.Args, names []string) (storage.Item, error) {
	name, err := argOrSelect(args, 0, "URL", names)
	handle(err)

	// Get item
	item, err := passline.GetSite(ctx, name)
	if err != nil {
		renderer.InvalidName(name)
		os.Exit(0)
	}

	return item, nil
}

func selectCredential(args ucli.Args, item storage.Item) (storage.Credential, error) {
	username, err := argOrSelect(args, 1, "Username/Login", item.GetUsernameArray())
	handle(err)

	// Check if name, username combination exists
	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		renderer.InvalidUsername(item.Name, username)
		os.Exit(0)
	}

	return credential, nil
}
