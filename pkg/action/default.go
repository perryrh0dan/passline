package action

import (
	"context"
	"fmt"
	"passline/pkg/cli/selection"
	"passline/pkg/clipboard"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v2"
)

func defaultParseArgs(c *ucli.Context) context.Context {
	ctx := ctxutil.WithGlobalFlags(c)
	if c.IsSet("category") {
		ctx = ctxutil.WithCategory(ctx, c.String("category"))
	}

	return ctx
}

func (s *Action) Default(c *ucli.Context) error {
	ctx := defaultParseArgs(c)

	// Get all Sites
	names, err := s.getItemNamesByCategory(ctx)
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

	if ctxutil.IsQuickSelect(ctx) {
		// disable notifications for quick select
		if err = clipboard.CopyTo(ctxutil.WithNotifications(ctx, false), "", []byte(credential.Username)); err != nil {
			return ExitError(ExitIO, err, "failed to copy to clipboard: %s", err)
		}
	}

	// get and check global password
	globalPassword, err := s.getMasterKey(ctx)
	if err != nil {
		return err
	}

	err = crypt.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	fmt.Printf("Category: %s\n", credential.Category)

	if ctxutil.IsAutoClip(ctx) || IsClip(ctx) {
		identifier := out.BuildIdentifier(name, credential.Username)
		if err = clipboard.CopyTo(ctx, identifier, []byte(credential.Password)); err != nil {
			return ExitError(ExitIO, err, "failed to copy to clipboard: %s", err)
		}
		if ctxutil.IsAutoClip(ctx) && !c.Bool("print") {
			out.SuccessfulCopiedToClipboard(name, credential.Username)
			return nil
		}
	}

	out.DisplayCredential(credential)
	out.SuccessfulCopiedToClipboard(name, credential.Username)
	return nil
}
