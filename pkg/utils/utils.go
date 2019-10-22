package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/perryrh0dan/passline/pkg/structs"
)

func RemoveFromArray(s []structs.Item, i int) []structs.Item {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func GetIndexOfItem(slice []structs.Item, item structs.Item) int {
	for p, v := range slice {
		if v == item {
			return p
		}
		return -1
	}
	return -1
}

func GeneratePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!$%&()/?")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	password := b.String() // E.g. "ExcbsVQs"
	return password
}
