package crypt

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	got, err := GenerateKey()
	if err != nil || len(got) != 32 {
		t.Errorf("GenerateKey() = length %d; wanted length %d", len(got), 32)
	}
}
