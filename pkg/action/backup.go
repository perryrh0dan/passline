package action

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"passline/pkg/cli/input"
	"passline/pkg/config"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Backup(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()
	out.BackupMessage()

	//TODO this should happen in config
	path := config.Directory() + "/backup"

	now := time.Now().Format("2006-01-02-15-04-05")
	path = filepath.Join(path, now)

	path, err := input.ArgOrInput(args, 0, "Path", path)
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
	items, err := s.Store.GetAllItems(ctx)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}

	time := time.Now()
	data := storage.Backup{Date: time, Items: items}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)

	return nil
}
