package action

import (
	"os"

	"passline/pkg/ctxutil"
	"passline/pkg/renderer"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) List(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()

	if args.Len() >= 1 {
		item, err := s.getSite(ctx, args.Get(0))
		if err != nil {
			renderer.InvalidName(args.Get(0))
			os.Exit(0)
		}

		renderer.DisplayItem(item)
	} else {
		items, err := s.getSites(ctx)
		if err != nil {
			return nil
		}

		if len(items) == 0 {
			renderer.NoItemsMessage()
		}
		renderer.DisplayItems(items)
	}

	return nil
}
