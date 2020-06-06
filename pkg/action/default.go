package action

import (
	"os"

	"passline/pkg/cli/selection"
	"passline/pkg/clipboard"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Default(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	// Get all Sites
	names, err := s.getSiteNames(ctx)
	if err != nil {
		return err
	}

	// Check if sites exists
	if len(names) <= 0 {
		return ExitError(ExitNotFound, err, "No items found")
	}

	args := c.Args()
	out.DisplayMessage()

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

	// get and check global password
	globalPassword, err := s.getGlobalPassword(ctx)
	if err != nil {
		return err
	}

	err = crypt.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	out.DisplayCredential(credential)

	identifier := out.BuildIdentifier(name, credential.Username)
	err = clipboard.CopyTo(ctx, identifier, []byte(credential.Password))
	if err != nil {
		out.ClipboardError()
		os.Exit(0)
	}

	out.SuccessfulCopiedToClipboard(item.Name, credential.Username)

	return nil
}
