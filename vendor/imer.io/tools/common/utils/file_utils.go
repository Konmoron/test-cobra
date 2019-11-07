package utils

import (
	"os"
)

// https://golangcode.com/check-if-a-file-exists/
func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
