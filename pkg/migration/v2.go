package migration

import (
	"passline/pkg/config"
	"passline/pkg/storage"
	"passline/pkg/util"
)

// 1.14.2 -> 2.0.0
func MigrateV2() error {
	cfg, err := config.Get(util.OSFileSystem{})
	if err != nil {
		return err
	}

	_, err = storage.New(cfg)
	if err != nil {
		return err
	}

	return nil
}
