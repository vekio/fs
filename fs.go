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

// CreateFileWithDirs creates a new file at the specified path, ensuring that all
// necessary parent directories are created with the specified permissions.
// If the file already exists, it is truncated. The file's permissions are also set accordingly.
func CreateFileWithDirs(path string, perms os.FileMode) (*os.File, error) {
	// Create all necessary parent directories.
	if err := os.MkdirAll(filepath.Dir(path), perms); err != nil {
		return nil, fmt.Errorf("creating directories for %s failed: %w", path, err)
	}

	// Create or truncate the file at the specified path.
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("creating file %s failed: %w", path, err)
	}

	// Set the permissions for the newly created file.
	if err := file.Chmod(perms); err != nil {
		file.Close() // Ensure the file is closed if setting permissions fails.
		return nil, fmt.Errorf("setting permissions for %s failed: %w", path, err)
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
