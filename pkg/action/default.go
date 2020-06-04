package action

import (
	"os"

	"passline/pkg/cli/selection"
	"passline/pkg/clipboard"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/renderer"

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
		renderer.NoItemsMessage()
		os.Exit(0)
	}

	args := c.Args()
	renderer.DisplayMessage()

	name, err := selection.ArgOrSelect(ctx, args, 0, "URL", names)
	if err != nil {
		return err
	}

	item, err := s.getSite(ctx, name)
	if err != nil {
		renderer.InvalidName(name)
		os.Exit(0)
	}

	credential, err := s.selectCredential(ctx, args, item)
	if err != nil {
		return err
	}

	// Get global password.
	globalPassword := s.getGlobalPassword(ctx)
	println()

	// globalPassword := []byte("test")

	err = crypt.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	renderer.DisplayCredential(credential)

	err = clipboard.CopyTo(ctx, credential.Username, []byte(credential.Password))
	if err != nil {
		renderer.ClipboardError()
		os.Exit(0)
	}

	renderer.SuccessfulCopiedToClipboard(item.Name, credential.Username)

	return nil
}
