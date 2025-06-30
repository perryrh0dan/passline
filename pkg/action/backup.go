package action

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"passline/pkg/cli/input"
	"passline/pkg/config"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) Backup(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	args := cmd.Args()
	out.BackupMessage()

	path := config.BackupDirectory()

	now := time.Now().Format("2006-01-02-15-04-05")
	path = filepath.Join(path, now) + ".json"

	path, err := input.ArgOrInput(args, 0, "Path", path, "required")
	if err != nil {
		return err
	}

	err = s.backup(ctx, path)
	if err != nil {
		return err
	}

	out.SuccessfulCreatedBackup(path)
	return nil
}

func (s *Action) backup(ctx context.Context, path string) error {
	items, err := s.Store.GetRawItems(ctx)
	if err != nil {
		return err
	}

	key, err := s.Store.GetKey(ctx)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}

	t := time.Now()
	type Alias storage.Backup
	aux := struct {
		Items json.RawMessage `json:"items"`
		*Alias
	}{
		Items: items,
		Alias: (*Alias)(&storage.Backup{
			Date: t,
			Key:  key,
		}),
	}

	file, err := json.MarshalIndent(aux, "", " ")
	if err != nil {
		return err
	}

	_ = os.WriteFile(path, file, 0644)

	return nil
}
