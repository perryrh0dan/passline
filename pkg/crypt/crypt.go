package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// AesGcmEncrypt takes an encryption key and a plaintext string and encrypts it with AES256 in GCM mode, which provides authenticated encryption. Returns the ciphertext and the used nonce.
func AesGcmEncrypt(key []byte, text string) (string, string, error) {
	// key := []byte(keyText)
	plaintextBytes := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := aesgcm.Seal(nil, nonce, plaintextBytes, nil)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), base64.URLEncoding.EncodeToString(nonce), nil
}

// AesGcmDecrypt takes an decryption key, a ciphertext and the corresponding nonce and decrypts it with AES256 in GCM mode. Returns the plaintext string.
func AesGcmDecrypt(key []byte, cryptoText string, nonce string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	nonceEncoded, _ := base64.URLEncoding.DecodeString(nonce)
	nonceBytes := []byte(nonceEncoded)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintextBytes, err := aesgcm.Open(nil, nonceBytes, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintextBytes), nil
}

func GenerateKey(password []byte) []byte {
	salt := []byte("This is the salt")
	dk := pbkdf2.Key(password, salt, 4096, 32, sha1.New)
	return dk
}
