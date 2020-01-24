package utils

import (
	"os"
)

func IsPathExists(path string) bool {
	if _, err := os.Stat(path); path != "" && !os.IsNotExist(err) && err == nil {
		return true
	}
	return false
}
