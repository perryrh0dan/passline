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
func AesGcmEncrypt(password []byte, text string) (string, error) {
	// Generate key from password with kdf
	key := GenerateKey(password)
	plaintextBytes := []byte(text)

	// Creation of the new block cipher based on the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Wrap the block cipher in a Galois Counter Mode (GCM) with standard nonce length
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a random nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// The first parameter is the prefix value
	ciphertext := aesgcm.Seal(nonce, nonce, plaintextBytes, nil)

	// Convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// AesGcmDecrypt takes an decryption key, a ciphertext and the corresponding nonce and decrypts it with AES256 in GCM mode. Returns the plaintext string.
func AesGcmDecrypt(password []byte, cryptoText string) (string, error) {
	// Generate key from password with kdf
	key := GenerateKey(password)

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintextBytes), nil
}

func GenerateKey(password []byte) []byte {
	salt := []byte("This is the salt")
	dk := pbkdf2.Key(password, salt, 4096, 32, sha1.New)
	return dk
}
