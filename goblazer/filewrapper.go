package goblazer

import (
	"os"
)

// IsFileExisted  a function wrapper to check file existed or not
func IsFileExisted(szFilePath string) bool {
	_, err := os.Stat(szFilePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
