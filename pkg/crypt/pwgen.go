package crypt

import (
	random "math/rand"
	"time"
)

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

	random.Seed(time.Now().UnixNano())
	var a = []rune{}

	// get the requirements
	if options.IncludeCharacters {
		a = append(a, lowercase[random.Intn(len(lowercase))])
		a = append(a, uppercase[random.Intn(len(uppercase))])
	}
	if options.IncludeNumbers {
		a = append(a, numbers[random.Intn(len(numbers))])
	}
	if options.IncludeSymbols {
		a = append(a, symbols[random.Intn(len(symbols))])
	}

	// populate the rest with random chars
	for len(a) < options.Length {
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
