package action

import (
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
		return ExitError(1, err, "Failed to enter name")
	}

	// User input username
	username, err := input.ArgOrInput(args, 1, "Username/Login", "")
	if err != nil {
		return ExitError(1, err, "Failed to enter username/login")
	}

	// Check if name, username combination exists
	exists := s.exists(ctx, name, username)
	if exists {
		identifier := out.BuildIdentifier(name, username)
		return ExitError(ExitDuplicated, err, "Item/Username already exists: %s", identifier)
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

	// get and check global password
	globalPassword, err := s.getGlobalPassword(ctx)
	if err != nil {
		return err
	}

	// Create Credentials
	credential := storage.Credential{Username: username, Password: password, RecoveryCodes: recoveryCodes}

	err = crypt.EncryptCredential(&credential, globalPassword)
	if err != nil {
		return ExitError(ExitEncrypt, err, "Error Encrypting credentials")
	}

	err = s.Store.AddCredential(ctx, name, credential)
	if err != nil {
		return ExitError(ExitUnknown, err, "Error occured: %s", err)
	}

	credential.Password = password

	out.DisplayCredential(credential)
	return nil
}
