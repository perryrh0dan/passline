package action

import (
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Me(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	username := ctxutil.GetDefaultUsername(ctx)
	phoneNumber := ctxutil.GetPhoneNumber(ctx)

	out.InfoCardMessage(username, phoneNumber)
	return nil
}
