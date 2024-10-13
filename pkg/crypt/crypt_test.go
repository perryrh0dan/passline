package crypt

import (
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
