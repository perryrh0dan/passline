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
