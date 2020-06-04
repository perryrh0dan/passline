package action

import (
	"context"
	"os"

	"passline/pkg/cli/input"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Add(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()
	out.CreateMessage()

	// User input name
	name, err := input.ArgOrInput(args, 0, "URL", "")
	if err != nil {
		ExitError(1, err, "Failed to enter name")
	}

	// User input username
	username, err := input.ArgOrInput(args, 1, "Username/Login", "")
	if err != nil {
		ExitError(1, err, "Failed to enter username/login")
	}

	// Check if name, username combination exists
	exists, err := s.exists(ctx, name, username)
	if err != nil {
		return err
	}

	if exists {
		out.NameUsernameAlreadyExists()
		return nil
	}

	password, err := input.Default("Please enter the existing Password []: ", "")
	if err != nil {
		return err
	}

	recoveryCodesString, err := input.Default("Please enter your recovery codes if exists []: ", "")
	if err != nil {
		return err
	}

	recoveryCodes := make([]string, 0)
	if recoveryCodesString != "" {
		recoveryCodes = util.StringToArray(recoveryCodesString)
	}

	globalPassword := s.getGlobalPassword(ctx)
	println()

	credential, err := s.AddItem(ctx, name, username, password, recoveryCodes, globalPassword)
	if err != nil {
		return err
	}

	out.DisplayCredential(credential)
	return nil
}

func (s *Action) AddItem(ctx context.Context, name, username, password string, recoveryCodes []string, globalPassword []byte) (storage.Credential, error) {

	// Check global password.
	valid, err := s.checkPassword(ctx, globalPassword)
	if err != nil || !valid {
		return storage.Credential{}, err
	}

	// Create Credentials
	credential := storage.Credential{Username: username, Password: password, RecoveryCodes: recoveryCodes}

	err = crypt.EncryptCredential(&credential, globalPassword)
	if err != nil {
		return storage.Credential{}, err
	}

	// Check if item already exists
	_, err = s.Store.GetItemByName(ctx, name)
	if err != nil {
		// Generate new item entry
		item := storage.Item{Name: name, Credentials: []storage.Credential{credential}}
		err = s.Store.CreateItem(ctx, item)
		if err != nil {
			os.Exit(0)
		}
	} else {
		// TODO check if credential already exists
		// Add to existing item
		err := s.Store.AddCredential(ctx, name, credential)
		if err != nil {
			os.Exit(0)
		}
	}

	credential.Password = password
	return credential, nil
}
