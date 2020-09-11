package crypt

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword(&Options{Length: 20})
	if err != nil {
		t.Error(err)
	}

	if len(password) != 20 {
		t.Errorf("GeneratePassword() = %s; wanted length %v", password, len(password))
	}
}

func TestGeneratePasswordWithCustomLength(t *testing.T) {
	password, err := GeneratePassword(&Options{Length: 10})
	if err != nil {
		t.Error(err)
	}

	if len(password) != 10 {
		t.Errorf("GeneratePassword() = %s; wanted length %v", password, len(password))
	}
}
