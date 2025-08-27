package action

import (
	"context"
	"passline/pkg/cli/selection"
	"passline/pkg/clipboard"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) Default(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	args := cmd.Args()
	out.DisplayMessage()

	names, err := s.getItemNamesByCategory(ctx)
	if err != nil {
		return err
	}

	// Check if sites exists
	if len(names) <= 0 {
		return ExitError(ExitNotFound, err, "No items found")
	}

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

	if ctxutil.IsQuickSelect(ctx) && !ctxutil.IsNoClip(ctx) {
		// disable notifications for quick select
		if err = clipboard.CopyTo(ctxutil.WithNotifications(ctx, false), "", []byte(credential.Username)); err != nil {
			out.ClipboardError()
		}
	}

	// get and check global password
	globalPassword, err := s.Store.GetDecryptedKey(ctx, "to decrypt the password")
	if err != nil {
		return err
	}

	err = storage.DecryptCredential(&credential, globalPassword)
	if err != nil {
		return err
	}

	if ctxutil.IsAutoClip(ctx) && !ctxutil.IsNoClip(ctx) {
		identifier := out.BuildIdentifier(name, credential.Username)
		if err = clipboard.CopyTo(ctx, identifier, []byte(credential.Password)); err != nil {
			out.FailedCopyToClipboard()
		} else {
			out.SuccessfulCopiedToClipboard(name, credential.Username)
		}
	}

	if cmd.Bool("print") {
		out.DisplayCredential(credential)
	}
	return nil
}
