package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vekio/fs/dir"
	"github.com/vekio/fs/file"
)

// Copy copies a file from src to dst, mimicking the behavior of the 'cp' command.
// If the destination is a directory, the file will be copied into it.
func Copy(src, dst string) error {
	// Check if the source file exists and is a file.
	exists, err := file.FileExists(src)
	if err != nil {
		return fmt.Errorf("error checking source file: %w", err)
	}
	if !exists {
		return fmt.Errorf("source file %s does not exist", src)
	}

	// Check if the destination is a directory.
	dstExists, err := dir.DirExists(dst)
	if err != nil {
		return fmt.Errorf("error checking destination: %w", err)
	}

	if dstExists {
		// If the destination is a directory, copy the file into the directory.
		dst = filepath.Join(dst, filepath.Base(src))
	} else {
		// Check if the parent directory of the destination exists.
		dstDir := filepath.Dir(dst)
		dirExists, err := dir.DirExists(dstDir)
		if err != nil {
			return fmt.Errorf("error checking destination directory: %w", err)
		}
		if !dirExists {
			return fmt.Errorf("destination directory %s does not exist", dstDir)
		}
	}

	return copyFileContents(src, dst)
}

// copyFileContents copies the contents of the source file to the destination file.
// It creates the destination file with the same permissions as the source.
func copyFileContents(src, dst string) error {
	// Open the source file for reading.
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer srcFile.Close()

	// Get file info for the source file to copy permissions.
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source file info: %w", err)
	}

	// Open the destination file for writing, creating or truncating it.
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("error creating or truncating destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}
