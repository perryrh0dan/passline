package action

import (
	"context"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) CategoryList(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	items, err := s.getSites(ctx)
	if err != nil {
		return ExitError(ExitUnknown, err, "Unable to load items")
	}

	category := ctxutil.GetCategory(ctx)

	var categories []string
	for _, item := range items {
		for _, cred := range item.Credentials {
			if !contains(categories, cred.Category) {
				categories = append(categories, cred.Category)
			}
		}
	}

	out.DisplayCategories(categories, category)

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
