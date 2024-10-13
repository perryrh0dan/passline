package storage

import (
	"context"
	"testing"

	"passline/pkg/config"
)

func TestAddItem(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Errorf("Unable to initialize config")
	}
	ctx := cfg.WithContext(context.Background())

	s, err := NewLocalStorage()
	if err != nil {
		t.Errorf("Unable to initialize storage")
	}

	credential := Credential{
		Username: "tpoe",
		Password: "1234",
	}

	s.AddCredential(ctx, "test", credential, []byte{})

	i, _ := s.GetItemByName(ctx, "test")
	c, _ := i.GetCredentialByUsername("tpoe")

	if c.Password != credential.Password || c.Username != c.Username {
		t.Errorf("Credential was not added correctly")
	}
}
