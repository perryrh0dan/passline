package action

import (
	"context"
	"strconv"

	"passline/pkg/cli/input"
	"passline/pkg/clipboard"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func generateParseArgs(c *ucli.Context) context.Context {
	ctx := ctxutil.WithGlobalFlags(c)
	if c.IsSet("advanced") {
		ctx = ctxutil.WithAdvanced(ctx, c.Bool("advanced"))
	}
	if c.IsSet("force") {
		ctx = ctxutil.WithForce(ctx, c.Bool("force"))
	}

	return ctx
}

func (s *Action) Generate(c *ucli.Context) error {
	ctx := generateParseArgs(c)

	args := c.Args()
	out.GenerateMessage()

	options := crypt.DefaultOptions()

	// User input name
	name, err := input.ArgOrInput(args, 0, "URL", "")
	if err != nil {
		return ExitError(ExitUnknown, err, "Failed to read input: %s", err)
	}

	// Load default user
	defaultUsername := ctxutil.GetDefaultUsername(ctx)

	// User input username
	username, err := input.ArgOrInput(args, 1, "Username/Login", defaultUsername)
	if err != nil {
		return ExitError(ExitUnknown, err, "Failed to read input: %s", err)
	}

	// Check if name, userrname combination exists
	exists := s.exists(ctx, name, username)
	if exists {
		identifier := out.BuildIdentifier(name, username)
		return ExitError(ExitDuplicated, err, "Item/Username already exists: %s", identifier)
	}

	// Get advanced parameters
	recoveryCodes := make([]string, 0)

	if ctxutil.IsAdvanced(ctx) {
		length, err := input.Default("Please enter the length of the password []: (%s) ", strconv.Itoa(options.Length))
		if err != nil {
			return err
		}
		options.Length, err = strconv.Atoi(length)
		if err != nil {
			return err
		}

		// check length
		if options.Length < 8 {
			out.PasswordTooShort()
			return nil
		}

		options.IncludeCharacters = input.Confirmation("Should the password include Character (y/n): ")
		options.IncludeNumbers = input.Confirmation("Should the password include Numbers (y/n): ")
		options.IncludeSymbols = input.Confirmation("Should the password include Symbols (y/n): ")

		recoveryCodesString, err := input.Default("Please enter your recovery codes if exists []: ", "")
		if err != nil {
			return err
		}

		if recoveryCodesString != "" {
			recoveryCodes = util.StringToArray(recoveryCodesString)
		}
	}

	password, err := crypt.GeneratePassword(&options)
	if err != nil {
		return err
	}

	// get and check global password
	globalPassword, err := s.getMasterKey(ctx)
	if err != nil {
		return err
	}

	// Create credentials
	credential := storage.Credential{Username: username, Password: password, RecoveryCodes: recoveryCodes}

	// Encrypt credentials
	err = crypt.EncryptCredential(&credential, globalPassword)
	if err != nil {
		return ExitError(ExitEncrypt, err, "Error Encrypting credentials")
	}

	// Save credentials
	err = s.Store.AddCredential(ctx, name, credential)
	if err != nil {
		return ExitError(ExitUnknown, err, "Error occured: %s", err)
	}

	// set unencrypted password to copy to clipboard and to show in terminal
	credential.Password = password

	if ctxutil.IsAutoClip(ctx) || IsClip(ctx) {
		identifier := out.BuildIdentifier(name, credential.Username)
		if err = clipboard.CopyTo(ctx, identifier, []byte(credential.Password)); err != nil {
			return ExitError(ExitIO, err, "failed to copy to clipboard: %s", err)
		}
		if ctxutil.IsAutoClip(ctx) && !c.Bool("print") {
			out.SuccessfulCopiedToClipboard(name, credential.Username)
			return nil
		}
	}

	out.DisplayCredential(credential)
	out.SuccessfulCopiedToClipboard(name, credential.Username)
	return nil
}
