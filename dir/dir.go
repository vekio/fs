package dir

import (
	"io/fs"
	"os"
	"path/filepath"

	_file "github.com/vekio/fs/file"
)

// DefaultPerms are defaults for new directory creation.
var DefaultPerms = 0755

// Create a new directory with the DefaultPerms creating any new
// directories as well (see os.MkdirAll).
func Create(path string) error {
	return os.MkdirAll(path, fs.FileMode(DefaultPerms))
}

// Copy recursively the contents of the source directory
// to the destination directory.
func Copy(src, dst string) error {
	// Get information about the source directory
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory if it does not exist
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	// List files in the source directory
	files, err := filepath.Glob(filepath.Join(src, "*"))
	if err != nil {
		return err
	}

	// Copy each file/directory to the destination directory
	for _, file := range files {
		newPath := filepath.Join(dst, filepath.Base(file))

		// If it's a directory, recursively call Copy
		if info, err := os.Stat(file); err == nil && info.IsDir() {
			if err := Copy(file, newPath); err != nil {
				return err
			}
		} else {
			// If it's a file, copy it
			if err := _file.Copy(file, newPath); err != nil {
				return err
			}
		}
	}

	return nil
}
