package dir

import (
	"fmt"
	"os"
	"path/filepath"

	_fs "github.com/vekio/fs"
	_file "github.com/vekio/fs/file"
)

// Copy recursively the contents of the source directory to the
// destination directory.
func Copy(src, dst string) error {
	// Get information about the source directory
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory if it does not exist
	exists, err := Exists(dst)
	if err != nil {
		return err
	}

	if !exists {
		if err := _fs.Create(dst, info.Mode()); err != nil {
			return err
		}
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
		if isDir, err := _fs.IsDir(file); err == nil && isDir {
			if err := Copy(file, newPath); err != nil {
				return fmt.Errorf("copy %s failed in %s: %w", file, newPath, err)
			}
		} else {
			// If it's a file, copy it
			if err := _file.Copy(file, newPath); err != nil {
				return fmt.Errorf("copy %s failed in %s: %w", file, newPath, err)
			}
		}
	}

	return nil
}

// Exists checks if the given dir path exists and is a directory.
func Exists(path string) (bool, error) {
	exists, err := _fs.Exists(path)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	isDir, err := _fs.IsDir(path)
	if err != nil {
		return false, err
	}

	return isDir, nil
}
