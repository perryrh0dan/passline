package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// LegacyStaticSalt is the salt used by v1/v2 vaults. It is kept here solely
// for backward-compatible decryption of existing key files. New vaults always
// use a randomly generated per-vault salt (see GenerateSalt).
var LegacyStaticSalt = []byte("This is the salt")

type Options struct {
	Length            int
	IncludeSymbols    bool
	IncludeNumbers    bool
	IncludeCharacters bool
}

func (options *Options) Validate() bool {
	if options.IncludeCharacters || options.IncludeNumbers || options.IncludeSymbols {
		return true
	}
	return false
}

func DefaultOptions() Options {
	return Options{
		Length:            20,
		IncludeCharacters: true,
		IncludeNumbers:    true,
		IncludeSymbols:    true,
	}
}

// GenerateSalt returns 16 cryptographically random bytes for use as a
// per-vault PBKDF2 salt.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// EncryptKey derives a KEK from password+salt (v3 KDF) and uses it to
// AES-GCM-encrypt key. Only call this when creating or re-keying a vault;
// decryption must go through DecryptKey which dispatches on KeyRecord.Version.
func EncryptKey(password, salt []byte, key string) (string, error) {
	kek := KeyRecord{Salt: salt, Version: KeyVersionV3}.DeriveKEK(password)
	encryptedKey, err := AesGcmEncrypt(kek, key)
	if err != nil {
		return "", err
	}
	return encryptedKey, nil
}

// DecryptKey uses the KDF parameters encoded in kr to derive the KEK, then
// AES-GCM-decrypts kr.EncryptedCEK. This is the single correct entry point
// for decryption — it handles both legacy and v3 vaults transparently.
func DecryptKey(password []byte, kr KeyRecord) (string, error) {
	kek := kr.DeriveKEK(password)
	key, err := AesGcmDecrypt(kek, kr.EncryptedCEK)
	if err != nil {
		return "", err
	}
	return key, nil
}

// AesGcmEncrypt takes an encryption key and a plaintext string and encrypts it with AES256 in GCM mode, which provides authenticated encryption. Returns the ciphertext and the used nonce.
func AesGcmEncrypt(key []byte, text string) (string, error) {
	plaintextBytes := []byte(text)

	if len(key) != 32 {
		return "", errors.New("Wrong key length")
	}

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
func AesGcmDecrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	if len(key) != 32 {
		return "", errors.New("Wrong key length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}
