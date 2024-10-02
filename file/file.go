package file

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Default perms for new file creation.
var DefaultFilePerms = os.FileMode(0644)

// Restricted perms for new file creation.
var RestrictedFilePerms = os.FileMode(0600)

// AppendToFile appends content to the file at the specified path.
// If the file does not exist, it will be created.
func AppendToFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s for appending: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error appending content to file %s: %w", path, err)
	}

	return nil
}

// CreateFile creates a new file or truncates an existing file.
func CreateFile(path string, perms os.FileMode) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return fmt.Errorf("error creating or truncating file %s: %w", path, err)
	}
	defer file.Close()

	return nil
}

// FileExists checks if the given path exists and is a file (not a directory).
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // File does not exist, no error.
		}
		return false, fmt.Errorf("error accessing %s: %w", path, err) // Other errors accessing the path.
	}

	if info.IsDir() {
		return false, fmt.Errorf("the path %s is a directory, not a file", path)
	}
	return true, nil
}

// GetFileSize returns the size of the file at the specified path.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("error getting file info for %s: %w", path, err)
	}
	return info.Size(), nil
}

// WriteFileContent writes the given content to a specified file.
// If the file does not exist, it creates it. If it exists, it truncates the content before writing.
func WriteFileContent(path string, content []byte, perms os.FileMode) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perms)
	if err != nil {
		return fmt.Errorf("error creating or opening file %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", path, err)
	}

	return nil
}

// MoveFile moves a file from src to dst. If dst exists, it will be overwritten.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("error moving file from %s to %s: %w", src, dst, err)
	}
	return nil
}

// ReadFile reads the entire content of the file at the given path.
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	return data, nil
}

// EditFile opens the specified file in the preferred editor.
// It first checks if the file exists before attempting to open it.
func EditFile(filePath string) error {
	// Check if the file exists before editing.
	exists, err := FileExists(filePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("file %s does not exist", filePath)
	}

	editor := getDefaultEditor()
	return executeCmd(editor, filePath)
}

// getDefaultEditor returns the user's preferred editor or a default one.
func getDefaultEditor() string {
	var editor string

	switch runtime.GOOS {
	case "linux", "darwin":
		// Check for preferred editor in environment variables
		editor = os.Getenv("VISUAL")
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		// If neither is set, use a sensible default
		if editor == "" {
			editor = "vi"
		}
	case "windows":
		editor = "notepad"
	default:
		// Default for unknown systems
		editor = "vi"
	}

	return editor
}

// executeCmd executes a command with the specified arguments.
func executeCmd(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command %s with arguments %v: %w", cmdName, args, err)
	}

	return nil
}
