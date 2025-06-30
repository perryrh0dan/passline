package action

import (
	"context"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v3"
)

// Sync config settings with vault e.g. reapply encryption mode
func (s *Action) Sync(c context.Context, cmd *ucli.Command) error {
	ctx := generateParseArgs(c, cmd)

	out.SyncMessage()

	items, err := s.Store.GetAllItems(ctx)
	if err != nil {
		return ExitError(ExitUnknown, err, "Failed to get all items: %s", err)
	}

	err = s.Store.SetItems(ctx, items, nil)
	if err != nil {
		return ExitError(ExitUnknown, err, "Failed to update data: %s", err)
	}

	return nil
}
