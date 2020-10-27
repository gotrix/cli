package helpers

import (
	"os"
)

// Exist returns true if file/dir exist.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// MakeDir create a dir recursively. Equivalent to mkdir -p.
func MakeDir(path string) error {
	if !Exists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// Rm removes file/dir. Equivalent to rm -rf.
func Rm(path string) error {
	if Exists(path) {
		return os.RemoveAll(path)
	}
	return nil
}
