package action

import (
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v2"
)

// Sync config settings with vault e.g. reapply encryption mode
func (s *Action) Sync(c *ucli.Context) error {
	ctx := generateParseArgs(c)

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
