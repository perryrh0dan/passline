package action

import (
	"context"
	"passline/pkg/cli/input"
	"passline/pkg/cli/selection"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Edit(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	// Get all items
	names, err := s.getSiteNames(ctx)
	if err != nil {
		return ExitError(ExitUnknown, err, "Error selecting item: %s", err)
	}

	// Check if any item exists
	if len(names) <= 0 {
		return ExitError(ExitNotFound, err, "No items found")
	}

	args := c.Args()
	out.DeleteMessage()

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
	globalPassword, err := s.getMasterKey(ctx)
	if err != nil {
		return err
	}

	// Decrypt Credentials to display secrets
	err = crypt.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	// Get new username
	newUsername, err := input.Default("Please enter a new Username/Login []: (%s) ", credential.Username, "")
	if err != nil {
		return err
	}
	credential.Username = newUsername

	// Get new category
	newCategory, err := input.Default("Please enter a new Category []: (%s ", credential.Category, "")
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
	credential.Category = newCategory
	credential.RecoveryCodes = make([]string, 0) // TODO remove spaces

	// use one space to clear recovery codes
	if newRecoveryCodes != " " {
		credential.RecoveryCodes = util.StringToArray(newRecoveryCodes)
	}

	err = s.edit(ctx, item.Name, selectedUsername, credential, globalPassword)
	if err != nil {
		return err
	}

	out.SuccessfulChangedItem(item.Name, credential.Username)

	return nil
}

func (s *Action) edit(ctx context.Context, name, username string, updatedCredential storage.Credential, globalPassword []byte) error {
	item, err := s.Store.GetItemByName(ctx, name)
	if err != nil {
		return err
	}

	// find credential in item
	var credential *storage.Credential

	for i := 0; i < len(item.Credentials); i++ {
		if item.Credentials[i].Username == username {
			credential = &item.Credentials[i]
			break
		}
	}

	*credential = updatedCredential

	err = crypt.EncryptCredential(credential, globalPassword)
	if err != nil {
		return err
	}

	err = s.Store.UpdateItem(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
