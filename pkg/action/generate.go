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

	ucli "github.com/urfave/cli/v3"
)

const PASSWORD_MIN_LENGTH = 8

func generateParseArgs(c context.Context, cmd *ucli.Command) context.Context {
	ctx := ctxutil.WithGlobalFlags(c, cmd)
	if cmd.IsSet("advanced") {
		ctx = ctxutil.WithAdvanced(ctx, cmd.Bool("advanced"))
	}
	if cmd.IsSet("force") {
		ctx = ctxutil.WithForce(ctx, cmd.Bool("force"))
	}

	return ctx
}

func (s *Action) Generate(c context.Context, cmd *ucli.Command) error {
	ctx := generateParseArgs(c, cmd)

	args := cmd.Args()
	out.GenerateMessage()

	options := crypt.DefaultOptions()

	// User input name
	name, err := input.ArgOrInput(args, 0, "URL", "", "required")
	if err != nil {
		return ExitError(ExitUnknown, err, "Failed to read input: %s", err)
	}

	// Load default user
	defaultUsername := ctxutil.GetDefaultUsername(ctx)

	// User input username
	username, err := input.ArgOrInput(args, 1, "Username/Login", defaultUsername, "required")
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

	// Get the default category from the context or use default
	category := ctxutil.GetCategory(ctx)
	if category == "*" {
		category = "default"
	}

	if ctxutil.IsAdvanced(ctx) {
		length, err := input.Default("Please enter the length of the password (%s): ", strconv.Itoa(options.Length), "number:8")
		if err != nil {
			return err
		}
		options.Length, err = strconv.Atoi(length)
		if err != nil {
			return err
		}

		// check length
		if options.Length < PASSWORD_MIN_LENGTH {
			out.PasswordTooShort()
			return nil
		}

		category, err = input.Default("Please enter a category (%s): ", "default", "")
		if err != nil {
			return err
		}

		if !ctxutil.IsAlwaysYes(ctx) {
			options.IncludeCharacters = input.Confirmation("Should the password include Characters (y/n): ")
			options.IncludeNumbers = input.Confirmation("Should the password include Numbers (y/n): ")
			options.IncludeSymbols = input.Confirmation("Should the password include Symbols (y/n): ")
		} else {
			options.IncludeCharacters = true
			options.IncludeNumbers = true
			options.IncludeSymbols = true
		}

		recoveryCodesString, err := input.Default("Please enter your recovery codes if exists []: ", "", "")
		if err != nil {
			return err
		}

		if recoveryCodesString != "" {
			recoveryCodes = util.StringToArray(recoveryCodesString)
		}
	}

	if !options.Validate() {
		out.InvalidGeneratorOptions()
		return nil
	}

	password, err := crypt.GeneratePassword(&options)
	if err != nil {
		return err
	}

	// get and check global password
	globalPassword, err := s.Store.GetDecryptedKey(ctx, "to encrypt the new password")
	if err != nil {
		return err
	}

	// Create credentials
	credential := storage.Credential{Username: username, Password: password, RecoveryCodes: recoveryCodes, Category: category}

	// Encrypt credentials
	err = storage.EncryptCredential(&credential, globalPassword)
	if err != nil {
		return ExitError(ExitEncrypt, err, "Error Encrypting credentials")
	}

	// Save credentials
	err = s.Store.AddCredential(ctx, name, credential)
	if err != nil {
		return ExitError(ExitUnknown, err, "Error occured: %s", err)
	}

	// Set decrypted password to copy to clipboard and to show in terminal
	credential.Password = password

	if ctxutil.IsAutoClip(ctx) {
		identifier := out.BuildIdentifier(name, credential.Username)
		if err = clipboard.CopyTo(ctx, identifier, []byte(credential.Password)); err != nil {
			out.FailedCopyToClipboard()
		} else {
			out.SuccessfulCopiedToClipboard(name, credential.Username)
		}
	}

	if cmd.Bool("print") {
		out.DisplayCredential(credential)
	}

	return nil
}
