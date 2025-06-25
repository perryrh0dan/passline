package action

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"passline/pkg/cli/input"
	"passline/pkg/crypt"
	"passline/pkg/ctxutil"
	"passline/pkg/out"
	"passline/pkg/storage"

	ucli "github.com/urfave/cli/v2"
)

func (s *Action) Restore(c *ucli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	args := c.Args()
	out.RestoreMessage()

	// User input path
	path, err := input.ArgOrInput(args, 0, "Path", "", "required")
	if err != nil {
		return err
	}

	if !ctxutil.HasAlwaysYes(ctx) {
		message := fmt.Sprintf("Are you sure you want to restore this  backup: %s (y/n): ", path)
		confirm := input.Confirmation(message)

		if !confirm {
			return nil
		}
	}

	err = s.restore(ctx, path)
	if err != nil {
		return err
	}

	out.SuccessfulRestoredBackup(path)
	return nil
}

func (s *Action) restore(ctx context.Context, path string) error {
	type Alias storage.Backup
	aux := struct {
		Items json.RawMessage `json:"items"`
		*Alias
	}{}

	_, err := os.Stat(path)
	if err != nil {
		out.InvalidFilePath()
		return err
	}

	file, _ := os.ReadFile(path)
	_ = json.Unmarshal([]byte(file), &aux)

	rawItems := json.RawMessage(aux.Items)

	globalPassword, err := input.MasterPassword(aux.Key, "to decrypt your vault")
	if err != nil {
		return err
	}

	var js json.RawMessage
	err = json.Unmarshal(aux.Items, &js)
	if err != nil {
		decryptedItems, err := crypt.AesGcmDecrypt(globalPassword, removeQuotes(string(aux.Items)))
		if err != nil {
			return err
		}

		rawItems = json.RawMessage(decryptedItems)
	}

	items := []storage.Item{}
	err = json.Unmarshal(rawItems, &items)
	if err != nil {
		return err
	}

	err = s.Store.SetKey(ctx, aux.Key)
	if err != nil {
		return err
	}

	err = s.Store.SetItems(ctx, items, globalPassword)
	if err != nil {
		return err
	}

	return nil
}

func removeQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	return s
}
