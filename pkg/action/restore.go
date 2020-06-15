package action

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"passline/pkg/cli/input"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Restore(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()
	out.RestoreMessage()

	// User input path
	path, err := input.ArgOrInput(args, 0, "Path", "")
	if err != nil {
		return err
	}

	if !ctxutil.HasAlwaysYes(ctx) {
		message := fmt.Sprintf("Are you sure you want to restore this  backup: %s (y/n): ", path)
		confirm := input.Confirmation(message)

		if !confirm {
			return nil
		}
	}

	err = s.restore(ctx, path)
	if err != nil {
		return err
	}

	out.SuccessfulRestoredBackup(path)
	return nil
}

func (s *Action) restore(ctx context.Context, path string) error {
	data := storage.Data{}

	_, err := os.Stat(path)
	if err != nil {
		out.InvalidFilePath()
		return err
	}

	file, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal([]byte(file), &data)

	err = s.Store.SetData(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
