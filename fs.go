package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Default perms for new directory creation.
var DefaultDirPerms = os.FileMode(0755)

// Restricted perms for new directory creation.
var RestrictedDirPerms = os.FileMode(0700)

// Default perms for new file creation.
var DefaultFilePerms = os.FileMode(0644)

// Restricted perms for new file creation.
var RestrictedFilePerms = os.FileMode(0600)

// CreateDir a new directory creating any new directories as well.
func CreateDir(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}

// CreateFile creates a new file, also creating any new directories
// if necessary. If the file already exists it is truncated.
func CreateFile(path string, perms os.FileMode) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), perms); err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if err := file.Chmod(perms); err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

// Exists checks if the given path exists on the filesystem.
// It returns true if the path exists, false if the path does not exist.
// Any other error encountered during the check (e.g., permission errors) is returned as an error.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("error accessing %s: %w", path, err) // Unexpected error.
}

// IsDir checks if the given path is a directory.
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("error checking if %s is a directory: %w", path, err)
	}

	return info.IsDir(), nil
}
