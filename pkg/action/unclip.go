package action

import (
	"context"
	"os"
	"time"

	ucli "github.com/urfave/cli/v3"

	"passline/pkg/clipboard"
)

// Unclip tries to erase the content of the clipboard
func (s *Action) Unclip(c context.Context, cmd *ucli.Command) error {
	force := cmd.Bool("force")
	timeout := cmd.Int("timeout")
	checksum := os.Getenv("PASSLINE_UNCLIP_CHECKSUM")

	time.Sleep(time.Second * time.Duration(timeout))
	if err := clipboard.Clear(c, checksum, force); err != nil {
		return ExitError(ExitIO, err, "Failed to clear clipboard: %s", err)
	}
	return nil
}
