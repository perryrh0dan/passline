package action

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"passline/pkg/cli/input"
	"passline/pkg/cli/selection"
	"passline/pkg/config"
	"passline/pkg/crypt"
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

func (s *Action) getMasterKey(ctx context.Context) ([]byte, error) {
	// Get encrypted content encryption key from store
	encryptedEncryptionKey, err := s.Store.GetKey(ctx)
	if err != nil {
		return []byte{}, ExitError(ExitUnknown, err, "Unable to load key: %s", err)
	}

	if encryptedEncryptionKey != "" {
		// If encrypted encryption key exists decrypt it
		envKey := []byte(os.Getenv("PASSLINE_MASTER_KEY"))
		if len(envKey) > 0 {
			encryptionKey, err := crypt.DecryptKey(envKey, encryptedEncryptionKey)
			if err == nil {
				return []byte(encryptionKey), nil
			}
		}

		// try password three times
		counter := 0
		for counter < 3 {
			password := input.Password("Enter master password: ")
			fmt.Println()

			encryptionKey, err := crypt.DecryptKey(password, encryptedEncryptionKey)
			if err == nil {
				return []byte(encryptionKey), nil
			} else if counter != 2 {
				fmt.Println("Wrong password! Please try again")
			}

			counter++
		}

		return []byte{}, ExitError(ExitPassword, err, "Wrong Password")
	}

	// initiate new encryption key
	encryptionKey, err := s.initMasterKey(ctx)
	if err != nil {
		return nil, err
	}

	return encryptionKey, nil
}

func (s *Action) initMasterKey(ctx context.Context) ([]byte, error) {
	decryptedEncryptionKey, err := crypt.GenerateKey()
	if err != nil {
		return []byte{}, ExitError(ExitUnknown, err, "Unable to generate key: %s", err)
	}

	password := input.Password("Enter master password: ")
	fmt.Println()
	passwordTwo := input.Password("Enter master password again: ")
	fmt.Println()

	if string(password) != string(passwordTwo) {
		return []byte{}, ExitError(ExitPassword, err, "Password do not match")
	}

	encryptedEncryptionKey, err := crypt.EncryptKey(password, decryptedEncryptionKey)
	if err != nil {
		return []byte{}, ExitError(ExitUnknown, err, "Unable to store key: %s", err)
	}
	s.Store.SetKey(ctx, encryptedEncryptionKey)

	return []byte(decryptedEncryptionKey), nil
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
