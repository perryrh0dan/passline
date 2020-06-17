package crypt

import (
	"passline/pkg/storage"
	"testing"
)

const (
	key = "01234567890123456789012345678912"
)

var (
	secret        string = "Hallo Perry"
	encryptedText string = "WUZJbLlKzYUJAx57rzUqUe4tB32cRwb7PoQceM-ad4LtPpO_ALHo"
)

func TestEncrypt(t *testing.T) {
	var err error
	encryptedText, err = AesGcmEncrypt([]byte(key), secret)
	if err != nil {
		t.Errorf("Encrypt(%s) Error occured", encryptedText)
	}
}

func TestDecrypt(t *testing.T) {
	got, err := AesGcmDecrypt([]byte(key), encryptedText)
	if err != nil || got != secret {
		t.Errorf("Decrypt(%s) = %s; wanted %s", encryptedText, got, secret)
	}
}

func TestEncryptWithShortKey(t *testing.T) {
	_, err := AesGcmEncrypt([]byte(key)[:29], secret)
	if err == nil {
		t.Errorf("Encrypt with short key should throw error")
	}
}

func TestEncryptCredentials(t *testing.T) {
	credential := storage.Credential{
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
	credential := storage.Credential{
		Username:      "perry",
		Password:      "3IbVQJqqSavvXiRdjffXWh3Z2d-4oxp_0zJ_VIDEcZmJ8aT_5g==", //123456789
		RecoveryCodes: []string{},
	}
	err := DecryptCredential(&credential, []byte(key))
	if err != nil || credential.Password != "123456789" {
		t.Errorf("EncryptCredential() Password = %s; wanted %s", credential.Password, "123456789")
	}
}
