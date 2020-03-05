package crypt

import (
	"testing"
)

var decryptedText string = "Hallo Perry"
var encryptedText string

func TestEncrypt(t *testing.T) {
	var err error
	encryptedText, err = AesGcmEncrypt([]byte("1234567891011123"), decryptedText)
	if err != nil {
		t.Errorf("Encrypt(%s) Error occured", encryptedText)
	}
}

func TestDecrypt(t *testing.T) {
	got, err := AesGcmDecrypt([]byte("1234567891011123"), encryptedText)
	if err != nil || got != decryptedText {
		t.Errorf("Decrypt(%s) = %s; wanted %s", encryptedText, got, decryptedText)
	}
}

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword(20)
	if err != nil {
		t.Error(err)
	}

	if len(password) != 20 {
		t.Errorf("GeneratePassword() = %s; wanted length %v", password, len(password))
	}
}
