package crypt

import (
	"crypto/rand"
	"math/big"
)

func secureIntn(n int) int {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return int(val.Int64())
}

func GeneratePassword(options *Options) (string, error) {
	lowercase := []rune("abcdefghijklmnopqrstuvwxyz")
	uppercase := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers := []rune("0123456789")
	symbols := []rune("!$%&()/?")
	var all = []rune{}

	if options.IncludeCharacters {
		all = append(all, lowercase...)
		all = append(all, uppercase...)
	}
	if options.IncludeNumbers {
		all = append(all, numbers...)
	}
	if options.IncludeSymbols {
		all = append(all, symbols...)
	}

	var a = []rune{}

	// get the requirements
	if options.IncludeCharacters {
		a = append(a, lowercase[secureIntn(len(lowercase))])
		a = append(a, uppercase[secureIntn(len(uppercase))])
	}
	if options.IncludeNumbers {
		a = append(a, numbers[secureIntn(len(numbers))])
	}
	if options.IncludeSymbols {
		a = append(a, symbols[secureIntn(len(symbols))])
	}

	// populate the rest with random chars
	for len(a) < options.Length {
		a = append(a, all[secureIntn(len(all))])
	}

	// shuffle up
	for i := 0; i < options.Length; i++ {
		randomPosition := secureIntn(options.Length)
		temp := a[i]
		a[i] = a[randomPosition]
		a[randomPosition] = temp
	}

	password := string(a)
	return password, nil
}
