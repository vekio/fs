package fs

import (
	"fmt"
	"os"
)

// Exists checks if the given path exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// El archivo o directorio existe
		return true, nil
	}

	if os.IsNotExist(err) {
		// El error indica que el archivo o directorio no existe
		return false, nil
	}

	// Otro tipo de error
	return false, fmt.Errorf("error checking existence of %s: %w", path, err)
}

// IsDir checks if the given path is a directory.
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}
	return info.IsDir(), nil
}
