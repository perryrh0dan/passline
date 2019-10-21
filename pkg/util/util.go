package util

import "github.com/perryrh0dan/passline/pkg/structs"

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
