package goblazer

import (
	"os"
)

// IsFileExisted is a function wrapper to check file existed or not
// parameter    : path - the path of the checked file
// return value : true - file existed, false - file not existed
func IsFileExisted(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
