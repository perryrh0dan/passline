package action

import (
	"fmt"

	"passline/pkg/cli/input"
	"passline/pkg/cli/selection"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Edit(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	// Get all Sites
	names, err := s.getItemNamesByCategory(ctx)
	if err != nil {
		return err
	}

	// Check if any item exists
	if len(names) <= 0 {
		return ExitError(ExitNotFound, err, "No items found")
	}

	args := c.Args()
	out.EditMessage()

	name, err := selection.ArgOrSelect(ctx, args, 0, "URL", names)
	if err != nil {
		return err
	}

	item, err := s.getSite(ctx, name)
	if err != nil {
		return ExitError(ExitNotFound, err, "Item not found: %s", name)
	}

	credential, err := s.selectCredential(ctx, args, item)
	if err != nil {
		return err
	}

	selectedUsername := credential.Username

	// get and check global password
	globalPassword, err := s.Store.GetDecryptedKey(ctx, "to decrypt the password")
	if err != nil {
		return err
	}

	// Decrypt Credentials to display secrets
	err = storage.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	// Get new URL
	newName, err := input.Default("Please enter a new URL []: (%s) ", item.Name, "")
	if err != nil {
		return err
	}

	// Get new username
	newUsername, err := input.Default("Please enter a new username/Login []: (%s) ", credential.Username, "")
	if err != nil {
		return err
	}

	// Get new category
	newCategory, err := input.Default("Please enter a new category []: (%s) ", credential.Category, "")
	if err != nil {
		return err
	}

	newComment, err := input.Default("Please enter a new comment []: (%s) ", credential.Comment, "")
	if err != nil {
		return err
	}

	// Get new recoveryCodes
	recoveryCodes := util.ArrayToString(credential.RecoveryCodes)
	newRecoveryCodes, err := input.Default("Please enter your recovery codes []: (%s) ", recoveryCodes, "")
	if err != nil {
		return err
	}

	// Edit credential
	credential.Username = newUsername
	credential.Category = newCategory
	credential.Comment = newComment
	credential.RecoveryCodes = make([]string, 0) // TODO remove spaces

	// use one space to clear recovery codes
	if newRecoveryCodes != " " {
		credential.RecoveryCodes = util.StringToArray(newRecoveryCodes)
	}

	err = storage.EncryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	existing := false
	newItem, err := s.Store.GetItemByName(ctx, newName)
	if err == nil {
		_, err = newItem.GetCredentialByUsername(credential.Username)
		if err == nil {
			existing = true
		}
	}

	if existing {
		identifier := out.BuildIdentifier(newName, credential.Username)
		message := fmt.Sprintf("Overwrite existing item %s []: (y/f) ", identifier)
		confirm := input.Confirmation(message)
		if !confirm {
			return nil
		}
	}

	err = s.Store.DeleteCredential(ctx, item, selectedUsername)
	if err != nil {
		return err
	}

	err = s.Store.AddCredential(ctx, newName, credential)
	if err != nil {
		return err
	}

	out.SuccessfulChangedItem(item.Name, credential.Username)

	return nil
}
