package util

import (
	"strings"
)

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

func LevenshteinDistanceSubstring(target, pattern string) (string, int) {
	minDistance := len(pattern)
	bestMatch := ""

	if len(target) >= len(pattern) {
		for i := 0; i <= len(target)-len(pattern); i++ {
			substring := target[i : i+len(pattern)]
			distance := LevenshteinDistance(substring, pattern)

			if distance < minDistance {
				minDistance = distance
				bestMatch = substring
			}
		}
	} else {
		for i := 0; i <= len(pattern)-len(target); i++ {
			substring := pattern[i : i+len(target)]
			distance := LevenshteinDistance(target, substring)

			if distance < minDistance {
				minDistance = distance
				bestMatch = substring
			}
		}
	}

	return bestMatch, minDistance
}

func LevenshteinDistance(a string, b string) int {
	rows := len(b) + 1
	cols := len(a) + 1

	matrix := make([][]int, rows)

	for i := range matrix {
		matrix[i] = make([]int, cols)
		matrix[i][0] = i
	}

	for i := range matrix[0] {
		matrix[0][i] = i
	}

	for y := 1; y < len(matrix); y++ {
		for x := 1; x < len(matrix[y]); x++ {
			if a[x-1] == b[y-1] {
				matrix[y][x] = matrix[y-1][x-1]
			} else {
				matrix[y][x] = min(matrix[y][x-1], matrix[y-1][x-1], matrix[y-1][x]) + 1
			}
		}
	}

	return matrix[rows-1][cols-1]
}
