package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	random "math/rand"
	"time"

	"passline/pkg/storage"

	"golang.org/x/crypto/pbkdf2"
)

type GeneratorOptions struct {
	Length int
}

func DecryptCredential(credential *storage.Credential, globalPassword []byte) error {
	err := DecryptPassword(credential, globalPassword)
	if err != nil {
		return err
	}

	err = DecryptRecoveryCodes(credential, globalPassword)
	if err != nil {
		return err
	}

	return nil
}

func DecryptPassword(credential *storage.Credential, globalPassword []byte) error {
	// Decrypt passwords
	var err error
	credential.Password, err = AesGcmDecrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func DecryptRecoveryCodes(credential *storage.Credential, globalPassword []byte) error {
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

func EncryptCredential(credential *storage.Credential, globalPassword []byte) error {
	err := EncryptPassword(credential, globalPassword)
	if err != nil {
		return err
	}

	err = EncryptRecoveryCodes(credential, globalPassword)
	if err != nil {
		return err
	}

	return nil
}

func EncryptPassword(credential *storage.Credential, globalPassword []byte) error {
	var err error
	credential.Password, err = AesGcmEncrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func EncryptRecoveryCodes(credential *storage.Credential, globalPassword []byte) error {
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

func GeneratePassword(options *GeneratorOptions) (string, error) {
	lowercase := []rune("abcdefghijklmnopqrstuvwxyz")
	uppercase := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers := []rune("0123456789")
	symbols := []rune("!$%&()/?")
	all := append(lowercase, uppercase...)
	all = append(all, numbers...)
	all = append(all, symbols...)
	random.Seed(time.Now().UnixNano())
	var a = []rune{}

	// get the requirements
	a = append(a, lowercase[random.Intn(len(lowercase))])
	a = append(a, uppercase[random.Intn(len(uppercase))])
	a = append(a, numbers[random.Intn(len(numbers))])
	a = append(a, symbols[random.Intn(len(symbols))])

	// populate the rest with random chars
	for i := 0; i < options.Length-4; i++ {
		a = append(a, all[random.Intn(len(all))])
	}

	// shuffle up
	for i := 0; i < options.Length; i++ {
		randomPosition := random.Intn(options.Length)
		temp := a[i]
		a[i] = a[randomPosition]
		a[randomPosition] = temp
	}

	password := string(a)
	return password, nil
}
