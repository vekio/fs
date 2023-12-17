package file

import (
	"fmt"
	"io"
	"os"

	_fs "github.com/vekio/fs"
)

// DefaultPerms for new file creation.
var DefaultPerms = 0600

// Copy copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func Copy(src, dst string) (err error) {
	// Get information about the source file
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}

	// Check if the source file is a regular file
	if !sfi.Mode().IsRegular() {
		// Cannot copy non-regular files
		// (e.g., directories, symbolic links, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)",
			sfi.Name(), sfi.Mode().String())
	}

	// Get information about the destination file, if it exists
	dfi, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check if the destination file is a regular file
	if err == nil && !dfi.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular destination file %s (%q)",
			dfi.Name(), dfi.Mode().String())
	}

	// Check if the source and destination files are the same
	if err == nil && os.SameFile(sfi, dfi) {
		return
	}

	// Try creating a hard link between the source and destination files
	if err = os.Link(src, dst); err == nil {
		return
	}

	// If creating a hard link fails,
	// copy the contents of the source file to the destination file
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	// Open the source file for reading
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	// Create the destination file for writing
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	// Copy the contents from the source file to the destination file
	if _, err = io.Copy(out, in); err != nil {
		return
	}

	// Ensure all contents are written to disk
	err = out.Sync()
	return
}

// Exists checks if the given file path exists and is not a directory.
func Exists(path string) (bool, error) {
	exists, err := _fs.Exists(path)
	if err != nil {
		return false, fmt.Errorf("error checking existence of %s: %w", path, err)
	}

	if !exists {
		return false, nil
	}

	isDir, err := _fs.IsDir(path)
	if err != nil {
		return false, fmt.Errorf("error checking if %s is a directory: %w", path, err)
	}

	return !isDir, nil
}
