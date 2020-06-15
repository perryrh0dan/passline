package out

import (
	"fmt"
	"strings"

	"passline/pkg/storage"
	"passline/pkg/util"

	"github.com/fatih/color"
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

	// TODO check if recovery codes exist
	if len(credential.RecoveryCodes) > 0 {
		fmt.Printf("Recovery codes: %s\n", util.ArrayToString(credential.RecoveryCodes))
	}
}

func DisplayItems(websites []storage.Item) {
	for _, website := range websites {
		fmt.Printf("%s\n", website.Name)
	}
}

func DisplayReleaseNotes(releaseNotes string) {
	fmt.Println("Release note:\n", releaseNotes)
}

func SuccessfulCopiedToClipboard(name, username string) {
	identifier := color.YellowString(BuildIdentifier(name, username))
	fmt.Fprintf(color.Output, "Copied Password for %s to clipboard\n", identifier)
}

func SuccessfulChangedItem(name, username string) {
	identifier := color.YellowString(BuildIdentifier(name, username))
	d := color.New(color.FgGreen)
	d.Printf("Successful changed item: %s\n", identifier)
}

func SuccessfulDeletedItem(name, username string) {
	identifier := color.YellowString(BuildIdentifier(name, username))
	d := color.New(color.FgGreen)
	d.Printf("Successful deleted item: %s\n", identifier)
}

func SuccessfulCreatedBackup(path string) {
	d := color.New(color.FgGreen)
	d.Printf("Successful created backup: %s\n", path)
}

func SuccessfulRestoredBackup(path string) {
	d := color.New(color.FgGreen)
	d.Printf("Successful restored backup: %s\n", path)
}

func SuccessfulUpdated(version string) {
	d := color.New(color.FgGreen)
	d.Printf("Successfully updated to version: %v\n", version)
}

func SuccessfulChangedPassword() {
	d := color.New(color.FgGreen)
	d.Printf("Successfully changed password\n")
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

func InvalidFilePath() {
	fmt.Printf("Invalid file path\n")
}

func InvalidInput() {
	fmt.Printf("Invalid input\n")
}

func ClipboardError() {
	fmt.Printf("Error occured while copying to clipboard\n")
}

func CoreInstanceError() {
	fmt.Printf("Error occured while instantiating core\n")
}

func StorageError() {
	d := color.New(color.FgRed)
	d.Printf("error: unable to initialice storage\n")
}

func DetectVersionError(err error) {
	d := color.New(color.FgRed)
	d.Printf("Error occurred while detecting version: %v\n", err)
}

func UpdateError(err error) {
	d := color.New(color.FgRed)
	d.Printf("Error occurred while updating binary: %v\n", err)
}

func NoUpdatesFound() {
	d := color.New(color.FgGreen)
	d.Printf("Current version is the latest\n")
}

func NoItemsMessage() {
	d := color.New(color.FgYellow)
	d.Printf("No items yet\n")
}

func DisplayMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Display item...\n")
}

func BackupMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Creating backup...\n")
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

func RestoreMessage() {
	d := color.New(color.FgGreen)
	d.Printf("Restoring backup...\n")
}

func MissingArgument(arguments []string) {
	d := color.New(color.FgRed)
	d.Printf("error: missing required arguments %s\n", strings.Join(arguments, ", "))
}

func NameAlreadyExists(name string) {
	fmt.Printf("error: name already exists %s\n", name)
}

func NameUsernameAlreadyExists() {
	fmt.Println("error: name & username combination already exists")
}

func BuildIdentifier(name, username string) string {
	return name + "/" + username
}
