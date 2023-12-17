package fs

import (
	"fmt"
	"os"
)

// Default perms for new directory creation.
var DefaultDirPerms = os.FileMode(0755)

// Default perms for new file creation.
var DefaultFilePerms = os.FileMode(0644)

// Create a new directory creating any new directories as well.
func Create(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}

// Exists checks if the given path exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	// The error indicates that the file or directory does not exist
	if os.IsNotExist(err) {
		return false, nil
	}

	// Unexpected error
	return false, fmt.Errorf("error checking existence of %s: %w", path, err)
}

// IsDir checks if the given path is a directory.
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("error checking if %s is a directory: %w", path, err)
	}

	return info.IsDir(), nil
}
