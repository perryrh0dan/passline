package crypt

import (
	"testing"
)

var encryptedText string = "Hallo Perry"
var nonce string
var cryptedText string

func TestEncrypt(t *testing.T) {
	var err error
	cryptedText, nonce, err = AesGcmEncrypt([]byte("1234567891011123"), encryptedText)
	if err != nil {
		t.Errorf("Encrypt(1234, %s) Error occured", encryptedText)
	}
}

func TestDecrypt(t *testing.T) {
	got, err := AesGcmDecrypt([]byte("1234567891011123"), cryptedText, nonce)
	if err != nil || got != encryptedText {
		t.Errorf("Encrypt(1234, %s) = %s; wanted %s", cryptedText, got, encryptedText)
	}
}
