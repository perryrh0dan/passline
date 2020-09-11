package crypt

import (
	"regexp"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	options := DefaultOptions()
	password, err := GeneratePassword(&options)
	if err != nil {
		t.Error(err)
	}

	if len(password) != 20 {
		t.Errorf("GeneratePassword() = %s; wanted length %v", password, len(password))
	}
}

func TestGeneratePasswordWithCustomLength(t *testing.T) {
	options := DefaultOptions()
	options.Length = 10
	password, err := GeneratePassword(&options)
	if err != nil {
		t.Error(err)
	}

	if len(password) != 10 {
		t.Errorf("GeneratePassword() = %s; wanted length %v", password, len(password))
	}
}

func TestGeneratePasswordWithCharacters(t *testing.T) {
	options := DefaultOptions()
	options.IncludeNumbers = false
	options.IncludeSymbols = false
	password, err := GeneratePassword(&options)
	if err != nil {
		t.Error(err)
	}

	rgx := regexp.MustCompile("\\d")
	matches := rgx.FindAllStringIndex(password, -1)

	if len(matches) != 0 {
		t.Errorf("GeneratePassword() = %s; wanted only character", password)
	}
}

func TestGeneratePasswordWithNumbers(t *testing.T) {
	options := DefaultOptions()
	options.IncludeCharacters = false
	options.IncludeSymbols = false
	password, err := GeneratePassword(&options)
	if err != nil {
		t.Error(err)
	}

	rgx := regexp.MustCompile("\\D")
	matches := rgx.FindAllStringIndex(password, -1)

	if len(matches) != 0 {
		t.Errorf("GeneratePassword() = %s; wanted only numbers", password)
	}
}

func TestGeneratePasswordWithSymbols(t *testing.T) {
	options := DefaultOptions()
	options.IncludeCharacters = false
	options.IncludeNumbers = false
	password, err := GeneratePassword(&options)
	if err != nil {
		t.Error(err)
	}

	rgx := regexp.MustCompile("[^!$%&()/?]")
	matches := rgx.FindAllStringIndex(password, -1)

	if len(matches) != 0 {
		t.Errorf("GeneratePassword() = %s; wanted only symbols", password)
	}
}
