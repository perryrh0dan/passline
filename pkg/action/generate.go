package action

import (
	"context"
	"fmt"
	"os"

	"passline/pkg/cli/input"
	"passline/pkg/clipboard"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/renderer"
	"passline/pkg/storage"
	"passline/pkg/util"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Generate(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()
	renderer.GenerateMessage()

	// User input name
	name, err := input.ArgOrInput(args, 0, "URL", "")
	if err != nil {
		return err
	}

	// User input username
	username, err := input.ArgOrInput(args, 1, "Username/Login", "")
	if err != nil {
		return err
	}
	// Check if name, username combination exists
	exists, err := s.exists(ctx, name, username)
	if err != nil {
		return err
	}

	if exists {
		os.Exit(0)
	}

	// Check if advanced mode is active
	if c.String("mode") == "advanced" {
		err := getAdvancedParamters(ctx)
		if err != nil {
			return err
		}
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

	credential, err := s.generate(ctx, name, username, recoveryCodes, globalPassword)
	if err != nil {
		return err
	}

	if ctxutil.IsAutoClip(ctx) || IsClip(ctx) {
		if err = clipboard.CopyTo(ctx, credential.Username, []byte(credential.Password)); err != nil {
			return ExitError(ExitIO, err, "failed to copy to clipboard: %s", err)
		}
		if ctxutil.IsAutoClip(ctx) && !c.Bool("print") {
			return nil
		}
	}

	renderer.SuccessfulCopiedToClipboard(name, credential.Username)
	return nil
}

func getAdvancedParamters(ctx context.Context) error {
	length, err := input.Default("Please enter the length of the password []: (%s)", "20")
	if err != nil {
		return err
	}

	fmt.Println(length)
	return nil
}

func (s *Action) generate(ctx context.Context, name, username string, recoveryCodes []string, globalPassword []byte) (storage.Credential, error) {
	// Generate password and crypt password
	password, err := crypt.GeneratePassword(20)
	if err != nil {
		return storage.Credential{}, err
	}

	return s.AddItem(ctx, name, username, password, recoveryCodes, globalPassword)
}
