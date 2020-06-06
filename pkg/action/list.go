package action

import (
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) List(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()

	if args.Len() >= 1 {
		item, err := s.getSite(ctx, args.Get(0))
		if err != nil {
			return ExitError(ExitNotFound, err, "Item not found: %s", args.Get(0))
		}

		out.DisplayItem(item)
	} else {
		items, err := s.getSites(ctx)
		if err != nil {
			return ExitError(ExitUnknown, err, "Unable to list items")
		}

		if len(items) == 0 {
			return ExitError(ExitNotFound, err, "No items found")
		}
		out.DisplayItems(items)
	}

	return nil
}
