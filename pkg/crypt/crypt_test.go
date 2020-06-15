package crypt

import (
	"testing"
)

var decryptedText string = "Hallo Perry"
var encryptedText string = "iY__LhvawdPwCjP43cApeRPPAbMhFHtR4oDLcPnKjXxLwSMkGQsw"

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

func TestGenerateKey(t *testing.T) {
	got, err := GenerateKey()
	if err != nil || len(got) != 32 {
		t.Errorf("GenerateKey() = length %d; wanted length %d", len(got), 32)
	}
}
