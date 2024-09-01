package dir

import (
	"fmt"
	"os"
)

// Default perms for new directory creation.
var DefaultDirPerms = os.FileMode(0755)

// Restricted perms for new directory creation.
var RestrictedDirPerms = os.FileMode(0700)

// Copy recursively the contents of the source directory to the
// destination directory.
// func Copy(src, dst string) error {
// 	// Get information about the source directory
// 	info, err := os.Stat(src)
// 	if err != nil {
// 		return err
// 	}

// 	// Create the destination directory if it does not exist
// 	exists, err := Exists(dst)
// 	if err != nil {
// 		return err
// 	}

// 	if !exists {
// 		if err := _fs.CreateDir(dst, info.Mode()); err != nil {
// 			return err
// 		}
// 	}

// 	// List files in the source directory
// 	files, err := filepath.Glob(filepath.Join(src, "*"))
// 	if err != nil {
// 		return err
// 	}

// 	// Copy each file/directory to the destination directory
// 	for _, file := range files {
// 		newPath := filepath.Join(dst, filepath.Base(file))

// 		// If it's a directory, recursively call Copy
// 		if isDir, err := _fs.IsDir(file); err == nil && isDir {
// 			if err := Copy(file, newPath); err != nil {
// 				return fmt.Errorf("copy %s failed in %s: %w", file, newPath, err)
// 			}
// 		} else {
// 			// If it's a file, copy it
// 			if err := _file.Copy(file, newPath); err != nil {
// 				return fmt.Errorf("copy %s failed in %s: %w", file, newPath, err)
// 			}
// 		}
// 	}

// 	return nil
// }

// CreateDir a new directory creating any new directories as well.
func Create(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}

// Exists checks if the given path exists and is a directory.
func Exists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Path does not exist, no error.
		}
		return false, fmt.Errorf("error accessing %s: %w", path, err) // Other errors accessing the path.
	}
	return info.IsDir(), nil // Return true if it's a directory, false otherwise.
}

// EnsureDir checks if a directory exists at the specified path;
// if not, it creates it with the specified permissions.
func EnsureDir(path string, perms os.FileMode) error {
	exists, err := Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		return Create(path, perms) // The directory does not exist, create it.
	}
	return nil // The directory already exists, nothing to do.
}
