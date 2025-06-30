package action

import (
	"context"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) Me(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	username := ctxutil.GetDefaultUsername(ctx)
	phoneNumber := ctxutil.GetPhoneNumber(ctx)

	out.InfoCardMessage(username, phoneNumber)
	return nil
}
