package action

import (
	"context"
	"errors"
	"fmt"

	"passline/pkg/cli/input"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"

	ucli "github.com/urfave/cli/v3"
)

func (s *Action) Password(c context.Context, cmd *ucli.Command) error {
	ctx := ctxutil.WithGlobalFlags(c, cmd)

	// Decrypt the current CEK (prompts for current master password).
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

	// Parse the current key record to detect legacy vaults for user messaging.
	rawKey, err := s.Store.GetKey(ctx)
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant read key: %s", err)
	}
	kr, err := crypt.ParseKeyString(rawKey)
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant parse key record: %s", err)
	}

	if kr.Version == crypt.KeyVersionLegacy {
		fmt.Println("Upgrading vault key format to v3 (per-vault salt)...")
	}

	// Always generate a fresh random salt — even if already v3, re-keying
	// rotates the salt which is good practice.
	salt, err := crypt.GenerateSalt()
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant generate salt: %s", err)
	}

	encryptedKey, err := crypt.EncryptKey(newPassword, salt, string(key))
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant encrypt key: %s", err)
	}

	err = s.Store.SetKey(ctx, crypt.EncodeKeyString(encryptedKey, salt))
	if err != nil {
		return ExitError(ExitUnknown, err, "Cant save to storage: %s", err)
	}

	out.SuccessfulChangedPassword()
	return nil
}
