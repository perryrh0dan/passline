package action

import (
	"context"
	"fmt"
	"os"

	"passline/pkg/cli/input"
	"passline/pkg/cli/selection"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Delete(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	// Get all Items
	names, err := s.getSiteNames(ctx)
	if err != nil {
		return err
	}

	// Check if any item exists
	if len(names) <= 0 {
		out.NoItemsMessage()
		return nil
	}

	args := c.Args()
	out.DeleteMessage()

	name, err := selection.ArgOrSelect(ctx, args, 0, "URL", names)
	if err != nil {
		return err
	}

	item, err := s.getSite(ctx, name)
	if err != nil {
		out.InvalidName(name)
		os.Exit(0)
	}

	credential, err := s.selectCredential(ctx, args, item)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Are you sure you want to delete this item: %s (y/n): ", credential.Username)
	confirm := input.Confirmation(message)
	if !confirm {
		return nil
	}

	err = s.delete(ctx, item.Name, credential.Username)
	if err != nil {
		return err
	}

	out.SuccessfulDeletedItem(item.Name, credential.Username)
	return nil
}

func (s *Action) delete(ctx context.Context, name, username string) error {
	item, err := s.Store.GetItemByName(ctx, name)
	if err != nil {
		return err
	}

	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		return err
	}

	err = s.Store.DeleteCredential(ctx, item, credential)
	if err != nil {
		os.Exit(0)
	}

	return nil
}
