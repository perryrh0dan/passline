package action

import (
	"os"

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
			out.InvalidName(args.Get(0))
			os.Exit(0)
		}

		out.DisplayItem(item)
	} else {
		items, err := s.getSites(ctx)
		if err != nil {
			return nil
		}

		if len(items) == 0 {
			out.NoItemsMessage()
		}
		out.DisplayItems(items)
	}

	return nil
}
