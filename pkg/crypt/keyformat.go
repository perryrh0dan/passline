package crypt

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// KeyVersion identifies the key derivation scheme used by a vault.
type KeyVersion int

const (
	// KeyVersionLegacy represents v1/v2 vaults: PBKDF2-HMAC-SHA1, 4096
	// iterations, static salt. Kept solely for backward-compatible decryption.
	KeyVersionLegacy KeyVersion = 1

	// KeyVersionV3 represents v3 vaults: PBKDF2-HMAC-SHA256, 600 000
	// iterations, per-vault random salt.
	KeyVersionV3 KeyVersion = 3
)

const keyRecordV3Prefix = "v3:"

// KeyRecord holds everything needed to derive the KEK and decrypt the CEK,
// independent of how or where the record is stored (file, Firestore, etc.).
type KeyRecord struct {
	EncryptedCEK string
	Salt         []byte
	Version      KeyVersion
}

// DeriveKEK derives the Key Encryption Key from password using the KDF
// parameters that match the record's Version.
func (kr KeyRecord) DeriveKEK(password []byte) []byte {
	switch kr.Version {
	case KeyVersionLegacy:
		return pbkdf2.Key(password, LegacyStaticSalt, 4096, 32, sha1.New)
	case KeyVersionV3:
		return pbkdf2.Key(password, kr.Salt, 600_000, 32, sha256.New)
	default:
		panic(fmt.Sprintf("crypt: unsupported key version %d", kr.Version))
	}
}

// ParseKeyString parses the raw content of a stored key record and returns a
// KeyRecord ready for decryption.
//
// Supported formats:
//
//	Legacy (v1/v2): <encryptedCEK>
//	    → KeyVersionLegacy; DeriveKEK will use PBKDF2-SHA1/4096/static-salt.
//
//	v3: v3:<base64url(salt)>:<encryptedCEK>
//	    → KeyVersionV3; DeriveKEK will use PBKDF2-SHA256/600k/per-vault-salt.
func ParseKeyString(content string) (KeyRecord, error) {
	if !strings.HasPrefix(content, keyRecordV3Prefix) {
		// Legacy vault: no prefix means v1/v2.
		return KeyRecord{
			EncryptedCEK: content,
			Salt:         LegacyStaticSalt,
			Version:      KeyVersionLegacy,
		}, nil
	}

	// v3:<base64url(salt)>:<encryptedCEK>
	rest := strings.TrimPrefix(content, keyRecordV3Prefix)
	idx := strings.Index(rest, ":")
	if idx == -1 {
		return KeyRecord{}, fmt.Errorf("malformed v3 key record: missing salt/CEK separator")
	}

	saltB64 := rest[:idx]
	encryptedCEK := rest[idx+1:]

	salt, err := base64.URLEncoding.DecodeString(saltB64)
	if err != nil {
		return KeyRecord{}, fmt.Errorf("malformed v3 key record: invalid salt encoding: %w", err)
	}
	if len(salt) == 0 {
		return KeyRecord{}, fmt.Errorf("malformed v3 key record: salt is empty")
	}

	return KeyRecord{
		EncryptedCEK: encryptedCEK,
		Salt:         salt,
		Version:      KeyVersionV3,
	}, nil
}

// EncodeKeyString encodes an encrypted CEK and its salt into the v3 format:
//
//	v3:<base64url(salt)>:<encryptedCEK>
func EncodeKeyString(encryptedCEK string, salt []byte) string {
	return keyRecordV3Prefix + base64.URLEncoding.EncodeToString(salt) + ":" + encryptedCEK
}
