package action

import (
	"context"
	"errors"

	"passline/pkg/cli/input"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) Password(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	// User input
	key, err := s.Store.GetDecryptedKey(ctx, "")
	if err != nil {
		return err
	}
	newPassword := input.Password("Enter new password: ")
	println()
	newPasswordTwo := input.Password("Enter new password again: ")
	println()

	if string(newPassword) != string(newPasswordTwo) {
		return ExitError(ExitIO, errors.New("Passwords do not match"), "Passwords do not match")
	}

	encryptedKey, err := crypt.EncryptKey(newPassword, string(key))
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant encrypt key: %s", err)
	}

	err = s.Store.SetKey(ctx, encryptedKey)
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant save to storage: %s", err)
	}

	out.SuccessfulChangedPassword()
	return nil
}
