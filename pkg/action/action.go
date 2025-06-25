package action

import (
	"context"
	"os"
	"path/filepath"

	"passline/pkg/cli/selection"
	"passline/pkg/config"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"

	"github.com/blang/semver"
	ucli "github.com/urfave/cli/v2"
)

// Action knows everything to run passline CLI actions
type Action struct {
	Name    string
	Store   storage.Storage
	cfg     *config.Config
	version semver.Version
}

func New(cfg *config.Config, sv semver.Version) (*Action, error) {
	return newAction(cfg, sv)
}

func newAction(cfg *config.Config, sv semver.Version) (*Action, error) {
	name := "passline"
	if len(os.Args) > 0 {
		name = filepath.Base(os.Args[0])
	}

	store, err := storage.New(cfg)
	if err != nil {
		return nil, ExitError(ExitUnknown, err, "Unable to initialize storage: %s", err)
	}

	act := &Action{
		Name:    name,
		cfg:     cfg,
		version: sv,
		Store:   store,
	}

	return act, nil
}

func (s *Action) selectCredential(ctx context.Context, args ucli.Args, item storage.Item) (storage.Credential, error) {
	category := ctxutil.GetCategory(ctx)

	username, err := selection.ArgOrSelect(ctx, args, 1, "Username/Login", item.GetUsernames(category))
	if err != nil {
		return storage.Credential{}, ExitError(ExitUnknown, err, "Selection Failed: %s", err)
	}

	// Check if name, username combination exists
	credential, err := item.GetCredentialByUsername(username)
	if err != nil {
		identifier := out.BuildIdentifier(item.Name, username)
		return storage.Credential{}, ExitError(ExitNotFound, err, "Username/Login not found: %s", identifier)
	}

	return credential, nil
}

func (s *Action) getSites(ctx context.Context) ([]storage.Item, error) {
	items, err := s.Store.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Action) getSiteNames(ctx context.Context) ([]selection.SelectItem, error) {
	items, err := s.getSites(ctx)
	if err != nil {
		return nil, err
	}

	var names []selection.SelectItem
	for _, item := range items {
		names = append(names, selection.SelectItem{Value: item.Name, Label: item.Name})
	}

	return names, nil
}

func (s *Action) getItemNamesByCategory(ctx context.Context) ([]selection.SelectItem, error) {
	category := ctxutil.GetCategory(ctx)

	if category == "*" {
		return s.getSiteNames(ctx)
	}

	items, err := s.getSites(ctx)
	if err != nil {
		return nil, err
	}

	var names []selection.SelectItem
	for _, item := range items {
		found := false
		for _, cred := range item.Credentials {
			if cred.Category == category {
				found = true
			}
		}
		if found {
			names = append(names, selection.SelectItem{Value: item.Name, Label: item.Name})
		}
	}

	return names, nil
}

func (s *Action) getSite(ctx context.Context, name string) (storage.Item, error) {
	item, err := s.Store.GetItemByName(ctx, name)
	if err != nil {
		return storage.Item{}, err
	}

	return item, nil
}

func (s *Action) exists(ctx context.Context, name, username string) bool {
	item, err := s.Store.GetItemByName(ctx, name)
	if err == nil {
		_, err = item.GetCredentialByUsername(username)
		if err == nil {
			return true
		}
	}

	return false
}
