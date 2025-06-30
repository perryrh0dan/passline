package action

import (
	"context"
	"fmt"

	"passline/pkg/cli/input"
	"passline/pkg/cli/selection"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) Delete(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	// Get all Sites
	names, err := s.getItemNamesByCategory(ctx)
	if err != nil {
		return err
	}

	// Check if any item exists
	if len(names) <= 0 {
		return ExitError(ExitNotFound, err, "No items found")
	}

	args := cmd.Args()
	out.DeleteMessage()

	name, err := selection.ArgOrSelect(ctx, args, 0, "URL", names)
	if err != nil {
		return ExitError(ExitUnknown, err, "Error selecting item: %s", err)
	}

	item, err := s.getSite(ctx, name)
	if err != nil {
		return ExitError(ExitNotFound, err, "Item not found: %s", name)
	}

	credential, err := s.selectCredential(ctx, args, item)
	if err != nil {
		return err
	}

	identifier := out.BuildIdentifier(item.Name, credential.Username)
	message := fmt.Sprintf("Are you sure you want to delete this item: %s (y/n): ", identifier)
	confirm := input.Confirmation(message)
	if !confirm {
		return nil
	}

	err = s.delete(ctx, item.Name, credential.Username)
	if err != nil {
		return ExitError(ExitUnknown, err, "Unable to delete item: %s", err)
	}

	out.SuccessfulDeletedItem(item.Name, credential.Username)
	return nil
}

func (s *Action) delete(ctx context.Context, name, username string) error {
	item, err := s.Store.GetItemByName(ctx, name)
	if err != nil {
		return err
	}

	_, err = item.GetCredentialByUsername(username)
	if err != nil {
		return err
	}

	err = s.Store.DeleteCredential(ctx, item, username)
	if err != nil {
		return err
	}

	return nil
}
