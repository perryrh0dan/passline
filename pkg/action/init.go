package action

import (
	"passline/pkg/ctxutil"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Init(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	_, err := s.Store.GetKey(ctx)
	if err == nil {
		println("Key found")
	} else {
		_, err = s.initMasterKey(ctx)
		if err != nil {

		}
	}

	return nil
}
