package storage

import "testing"

const (
	key = "01234567890123456789012345678912"
)

func TestEncryptCredentials(t *testing.T) {
	credential := Credential{
		Username:      "perry",
		Password:      "123456789",
		RecoveryCodes: []string{"test", "tee"},
	}
	err := EncryptCredential(&credential, []byte(key))
	if err != nil || credential.Password == "123456789" {
		t.Errorf("EncryptCredential() Password = %s; wanted != %s", credential.Password, "123456789")
	}
}

func TestDecryptCredentials(t *testing.T) {
	credential := Credential{
		Username:      "perry",
		Password:      "3IbVQJqqSavvXiRdjffXWh3Z2d-4oxp_0zJ_VIDEcZmJ8aT_5g==", //123456789
		RecoveryCodes: []string{},
	}
	err := DecryptCredential(&credential, []byte(key))
	if err != nil || credential.Password != "123456789" {
		t.Errorf("EncryptCredential() Password = %s; wanted %s", credential.Password, "123456789")
	}
}
