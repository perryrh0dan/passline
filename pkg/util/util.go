package util

import "strings"

func ArrayContains(l []string, i string) bool {
	for _, li := range l {
		if li == i {
			return true
		}
	}

	return false
}

func FilterArray(l []string, filter string) []string {
	filteredNames := make([]string, 0)
	for _, i := range l {
		if strings.Contains(i, filter) {
			filteredNames = append(filteredNames, i)
		}
	}
	return filteredNames
}

func ArrayToString(l []string) string {
	return strings.Join(l, ",")
}

func StringToArray(s string) []string {
	return strings.Split(s, ",")
}
