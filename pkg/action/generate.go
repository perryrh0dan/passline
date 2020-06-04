package action

import (
	"context"
	"os"
	"strconv"

	"passline/pkg/cli/input"
	"passline/pkg/clipboard"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func generateParseArgs(c *ucli.Context) context.Context {
	ctx := ctxutil.WithGlobalFlags(c)
	if c.IsSet("advanced") {
		ctx = ctxutil.WithAdvanced(ctx, c.Bool("clip"))
	}
	if c.IsSet("force") {
		ctx = ctxutil.WithForce(ctx, c.Bool("force"))
	}

	return ctx
}

func (s *Action) Generate(c *ucli.Context) error {
	ctx := generateParseArgs(c)
	// force := c.Bool("force")

	args := c.Args()
	out.GenerateMessage()

	options := crypt.GeneratorOptions{
		Length: 20,
	}

	// User input name
	name, err := input.ArgOrInput(args, 0, "URL", "")
	if err != nil {
		return err
	}

	// User input username
	username, err := input.ArgOrInput(args, 1, "Username/Login", "")
	if err != nil {
		return ExitError(1, err, "")
	}

	// Check if name, username combination exists
	exists, err := s.exists(ctx, name, username)
	if err != nil {
		return ExitError(1, err, "")
	}

	if exists {
		return ExitError(1, err, "Item/Username already exists")
	}

	// Get advanced parameters
	recoveryCodes := make([]string, 0)
	if c.Bool("advanced") {
		length, err := input.Default("Please enter the length of the password []: (%s) ", "20")
		if err != nil {
			return err
		}
		options.Length, _ = strconv.Atoi(length)

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

	globalPassword := s.getGlobalPassword(ctx)
	println()

	credential, err := s.AddItem(ctx, name, username, password, recoveryCodes, globalPassword)
	if err != nil {
		os.Exit(1)
	}

	if ctxutil.IsAutoClip(ctx) || IsClip(ctx) {
		if err = clipboard.CopyTo(ctx, credential.Username, []byte(credential.Password)); err != nil {
			return ExitError(ExitIO, err, "failed to copy to clipboard: %s", err)
		}
		if ctxutil.IsAutoClip(ctx) && !c.Bool("print") {
			return nil
		}
	}

	out.SuccessfulCopiedToClipboard(name, credential.Username)
	return nil
}
