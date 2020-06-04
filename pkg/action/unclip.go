package action

import (
	"os"
	"time"

	ucli "github.com/urfave/cli/v2"

	"passline/pkg/clipboard"
)

// Unclip tries to erase the content of the clipboard
func (s *Action) Unclip(c *ucli.Context) error {
	ctx := c.Context
	force := c.Bool("force")
	timeout := c.Int("timeout")
	checksum := os.Getenv("PASSLINE_UNCLIP_CHECKSUM")

	time.Sleep(time.Second * time.Duration(timeout))
	if err := clipboard.Clear(ctx, checksum, force); err != nil {
		return ExitError(ExitIO, err, "Failed to clear clipboard: %s", err)
	}
	return nil
}
