package fsutil

import (
	"fmt"
	"os"
)

// IsFile checks if a certain path is actually a file
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// not found
			return false
		}
		fmt.Printf("failed to check dir %s: %s\n", path, err)
		return false
	}

	return fi.Mode().IsRegular()
}
