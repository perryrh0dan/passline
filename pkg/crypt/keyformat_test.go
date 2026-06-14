package crypt

import (
	"bytes"
	"crypto/sha1"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

// TestParseKeyString_Legacy verifies that a bare encrypted CEK (no prefix) is
// parsed as a legacy record and returns the static salt.
func TestParseKeyString_Legacy(t *testing.T) {
	const encCEK = "someencryptedvalue"
	kr, err := ParseKeyString(encCEK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kr.Version != KeyVersionLegacy {
		t.Errorf("Version = %d; want %d (KeyVersionLegacy)", kr.Version, KeyVersionLegacy)
	}
	if kr.EncryptedCEK != encCEK {
		t.Errorf("EncryptedCEK = %q; want %q", kr.EncryptedCEK, encCEK)
	}
	if !bytes.Equal(kr.Salt, LegacyStaticSalt) {
		t.Errorf("Salt = %q; want LegacyStaticSalt", kr.Salt)
	}
}

// TestParseKeyString_V3RoundTrip verifies that EncodeKeyString → ParseKeyString is
// lossless for v3 records.
func TestParseKeyString_V3RoundTrip(t *testing.T) {
	salt := []byte("0123456789abcdef") // 16 bytes
	const encCEK = "anotherencryptedvalue"

	encoded := EncodeKeyString(encCEK, salt)
	kr, err := ParseKeyString(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kr.Version != KeyVersionV3 {
		t.Errorf("Version = %d; want %d (KeyVersionV3)", kr.Version, KeyVersionV3)
	}
	if kr.EncryptedCEK != encCEK {
		t.Errorf("EncryptedCEK = %q; want %q", kr.EncryptedCEK, encCEK)
	}
	if !bytes.Equal(kr.Salt, salt) {
		t.Errorf("Salt = %v; want %v", kr.Salt, salt)
	}
}

// TestParseKeyString_MalformedV3 checks that malformed v3 content returns errors.
func TestParseKeyString_MalformedV3(t *testing.T) {
	cases := []struct {
		name    string
		content string
	}{
		{"missing separator", "v3:AAAA"},
		{"bad base64 salt", "v3:!!!:encCEK"},
		{"empty salt", "v3::encCEK"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseKeyString(tc.content)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tc.content)
			}
		})
	}
}

// TestDeriveKEK_LegacyMatchesOldKDF verifies that a legacy KeyRecord produces
// exactly the same KEK as the original PBKDF2-SHA1/4096 derivation that was
// used to encrypt all v1/v2 vaults.
func TestDeriveKEK_LegacyMatchesOldKDF(t *testing.T) {
	password := []byte("hunter2")

	// Replicate the exact original key derivation.
	want := pbkdf2.Key(password, []byte("This is the salt"), 4096, 32, sha1.New)

	kr := KeyRecord{
		EncryptedCEK: "",
		Salt:         LegacyStaticSalt,
		Version:      KeyVersionLegacy,
	}
	got := kr.DeriveKEK(password)

	if !bytes.Equal(got, want) {
		t.Errorf("legacy DeriveKEK does not match original KDF\ngot:  %x\nwant: %x", got, want)
	}
}

// TestDeriveKEK_V3Distinct verifies that a v3 KeyRecord with a random salt
// produces a different KEK than the legacy KDF for the same password.
func TestDeriveKEK_V3Distinct(t *testing.T) {
	password := []byte("hunter2")
	salt := []byte("randomsalt000000") // 16 bytes, intentionally differs from LegacyStaticSalt

	legacyKR := KeyRecord{Salt: LegacyStaticSalt, Version: KeyVersionLegacy}
	v3KR := KeyRecord{Salt: salt, Version: KeyVersionV3}

	legacyKEK := legacyKR.DeriveKEK(password)
	v3KEK := v3KR.DeriveKEK(password)

	if bytes.Equal(legacyKEK, v3KEK) {
		t.Error("legacy and v3 KEKs should differ but are identical")
	}
}

// TestDecryptKey_LegacyRoundTrip verifies full encrypt-then-decrypt using the
// legacy KDF path — simulating decryption of a vault created before v3.
func TestDecryptKey_LegacyRoundTrip(t *testing.T) {
	password := []byte("correct horse battery staple")
	plainCEK := "the-actual-vault-key-content"

	// Encrypt using the old KEK (SHA1/4096/static-salt) directly.
	legacyKEK := pbkdf2.Key(password, LegacyStaticSalt, 4096, 32, sha1.New)
	encCEK, err := AesGcmEncrypt(legacyKEK, plainCEK)
	if err != nil {
		t.Fatalf("AesGcmEncrypt failed: %v", err)
	}

	// Decrypt via the new API using a legacy KeyRecord.
	kr := KeyRecord{
		EncryptedCEK: encCEK,
		Salt:         LegacyStaticSalt,
		Version:      KeyVersionLegacy,
	}
	got, err := DecryptKey(password, kr)
	if err != nil {
		t.Fatalf("DecryptKey failed: %v", err)
	}
	if got != plainCEK {
		t.Errorf("DecryptKey = %q; want %q", got, plainCEK)
	}
}
