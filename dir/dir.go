package dir

import (
	"fmt"
	"os"
)

// Default perms for new directory creation.
var DefaultDirPerms = os.FileMode(0755)

// Restricted perms for new directory creation.
var RestrictedDirPerms = os.FileMode(0700)

// CreateDir creates a new directory along with any necessary parents.
func CreateDir(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}

// DirExists checks if the given path exists and is a directory.
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Path does not exist, no error.
		}
		return false, fmt.Errorf("error accessing %s: %w", path, err) // Other errors accessing the path.
	}
	return info.IsDir(), nil
}

// EnsureDir checks if a directory exists at the specified path;
// if not, it creates it with the specified permissions.
func EnsureDir(path string, perms os.FileMode) error {
	exists, err := DirExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return CreateDir(path, perms) // The directory does not exist, create it.
	}
	return nil // The directory already exists, nothing to do.
}

// ListDir lists only the directories in the specified path.
func ListDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", path, err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}

// IsEmptyDir checks if the specified directory is empty.
func IsEmptyDir(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("error reading directory %s: %w", path, err)
	}
	return len(entries) == 0, nil
}
