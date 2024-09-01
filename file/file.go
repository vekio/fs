package file

import (
	"fmt"
	"os"
	"os/exec"
)

// Default perms for new file creation.
var DefaultFilePerms = os.FileMode(0644)

// Restricted perms for new file creation.
var RestrictedFilePerms = os.FileMode(0600)

// Touch updates the access and modification times of a file at the
// specified path. If the file does not exist, it creates an empty file
// with the specified permissions.
// func Touch(path string, perms os.FileMode) error {
// 	// Check if the file already exists
// 	exists, err := Exists(path)
// 	if err != nil {
// 		return err
// 	}

// 	// If the file exists, update its access and modification times
// 	// to the current time
// 	if exists {
// 		now := time.Now().Local()
// 		if err := os.Chtimes(path, now, now); err != nil {
// 			return err
// 		}
// 	}

// 	// If the file does not exist, create the necessary directory
// 	// and the empty file
// 	if !exists {
// 		// Create the directory (if needed) using default directory permissions
// 		if err := _fs.CreateDir(filepath.Dir(path), _fs.DefaultDirPerms); err != nil {
// 			return err
// 		}

// 		// Create the empty file
// 		file, err := os.Create(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		// Set the file permissions to the specified permissions
// 		if err := os.Chmod(path, perms); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// Copy copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
// func Copy(src, dst string) error {
// 	// Get information about the source file
// 	sfi, err := os.Stat(src)
// 	if err != nil {
// 		return nil
// 	}

// 	// Check if the source file is a regular file
// 	if !sfi.Mode().IsRegular() {
// 		// Cannot copy non-regular files
// 		// (e.g., directories, symbolic links, devices, etc.)
// 		return fmt.Errorf("can't copy non-regular file %s (%q)",
// 			sfi.Name(), sfi.Mode().String())
// 	}

// 	// Get information about the destination file, if it exists
// 	exists, err := Exists(dst)
// 	if err != nil {
// 		return err
// 	}

// 	if exists {
// 		dfi, _ := os.Stat(dst)

// 		// Check if the destination file is a regular file
// 		if !dfi.Mode().IsRegular() {
// 			return fmt.Errorf("non-regular destination file %s (%q)",
// 				dfi.Name(), dfi.Mode().String())
// 		}

// 		// Check if the source and destination files are the same
// 		if os.SameFile(sfi, dfi) {
// 			return nil
// 		}
// 	}

// 	// Copy the contents of the source file to the destination file
// 	err = copyFileContents(src, dst)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// copies the contents of the file named src to the file named by dst.
// The file will be created if it does not already exist. If the destination
// file exists, all it's contents will be replaced by the contents
// of the source file.
// func copyFileContents(src, dst string) error {
// 	// Open the source file for reading
// 	in, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer in.Close()

// 	// Create the destination file for writing
// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		cerr := out.Close()
// 		if err == nil {
// 			err = cerr
// 		}
// 	}()

// 	// Copy the contents from the source file to the destination file
// 	if _, err = io.Copy(out, in); err != nil {
// 		return err
// 	}

// 	// Ensure all contents are written to disk
// 	err = out.Sync()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Create creates a new file and returns a file handle.
// If the file already exists, it simply opens it with truncation.
func Create(path string, perms os.FileMode) (*os.File, error) {
	// Try to open the file in read/write mode, create it if it does not exist, and truncate it if it does.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return nil, fmt.Errorf("error creating or opening file %s: %w", path, err)
	}
	return file, nil
}

// Exists checks if the given path exists and is a file (not a directory).
func Exists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Path does not exist, no error.
		}
		return false, fmt.Errorf("error accessing %s: %w", path, err) // Other errors accessing the path.
	}
	if info.IsDir() {
		return false, fmt.Errorf("path is a directory and not a file")
	}
	return true, nil // The path exists and is not a directory.
}

// Edit opens the specified file in the preferred editor.
func Edit(filePath string) error {
	editor := getDefaultEditor()
	return executeCmd(editor, filePath)
}

func getDefaultEditor() string {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi" // TODO Include "notepad" for Windows.
	}
	return editor
}

func executeCmd(editor, filePath string) error {
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", editor, err)
	}
	return nil
}
