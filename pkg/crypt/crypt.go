package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"passline/pkg/storage"

	"golang.org/x/crypto/pbkdf2"
)

type GeneratorOptions struct {
	Length int
}

func DecryptCredential(credential *storage.Credential, globalPassword []byte) error {
	err := decryptPassword(credential, globalPassword)
	if err != nil {
		return err
	}

	err = decryptRecoveryCodes(credential, globalPassword)
	if err != nil {
		return err
	}

	return nil
}

func EncryptCredential(credential *storage.Credential, key []byte) error {
	err := encryptPassword(credential, key)
	if err != nil {
		return err
	}

	err = encryptRecoveryCodes(credential, key)
	if err != nil {
		return err
	}

	return nil
}

func EncryptKey(password []byte, key string) (string, error) {
	pwKey := GetKey(password)
	encryptedKey, err := AesGcmEncrypt(pwKey, key)
	if err != nil {
		return "", err
	}
	return encryptedKey, nil
}

func DecryptKey(password []byte, encryptedKey string) (string, error) {
	pwKey := GetKey(password)
	key, err := AesGcmDecrypt(pwKey, encryptedKey)
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
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

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
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintextBytes), nil
}

func GetKey(password []byte) []byte {
	salt := []byte("This is the salt")
	dk := pbkdf2.Key(password, salt, 4096, 32, sha1.New)
	return dk
}

func encryptPassword(credential *storage.Credential, key []byte) error {
	var err error
	credential.Password, err = AesGcmEncrypt(key, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func encryptRecoveryCodes(credential *storage.Credential, globalPassword []byte) error {
	var encryptedRecoveryCodes = make([]string, 0)

	for _, c := range credential.RecoveryCodes {
		encryptedRecoveryCode, err := AesGcmEncrypt(globalPassword, c)
		if err != nil {
			return err
		}
		encryptedRecoveryCodes = append(encryptedRecoveryCodes, encryptedRecoveryCode)
	}

	credential.RecoveryCodes = encryptedRecoveryCodes
	return nil
}

func decryptPassword(credential *storage.Credential, globalPassword []byte) error {
	// Decrypt passwords
	var err error
	credential.Password, err = AesGcmDecrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func decryptRecoveryCodes(credential *storage.Credential, globalPassword []byte) error {
	var decryptedRecoveryCodes = make([]string, 0)
	for _, c := range credential.RecoveryCodes {
		decryptedRecoveryCode, err := AesGcmDecrypt(globalPassword, c)
		if err != nil {
			return err
		}
		decryptedRecoveryCodes = append(decryptedRecoveryCodes, decryptedRecoveryCode)
	}

	credential.RecoveryCodes = decryptedRecoveryCodes
	return nil
}
