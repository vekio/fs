package fs

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCopyFileToFile tests copying from one file to another.
func TestCopyFileToFile(t *testing.T) {
	// Create a temporary source file.
	srcFile, err := os.CreateTemp("", "source_file.txt")
	if err != nil {
		t.Fatalf("error creating temp source file: %v", err)
	}
	defer os.Remove(srcFile.Name()) // Clean up after test.

	// Write content to the source file.
	content := []byte("This is the source file content.")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("error writing to source file: %v", err)
	}
	srcFile.Close()

	// Create a temporary destination file.
	dstFile, err := os.CreateTemp("", "destination_file.txt")
	if err != nil {
		t.Fatalf("error creating temp destination file: %v", err)
	}
	defer os.Remove(dstFile.Name()) // Clean up after test.

	// Call the Copy function.
	err = Copy(srcFile.Name(), dstFile.Name())
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify the contents of the destination file.
	data, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatalf("error reading destination file: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("expected content %s, but got %s", string(content), string(data))
	}
}

// TestCopyFileToDirectory tests copying a file into a directory.
func TestCopyFileToDirectory(t *testing.T) {
	// Create a temporary source file.
	srcFile, err := os.CreateTemp("", "source_file.txt")
	if err != nil {
		t.Fatalf("error creating temp source file: %v", err)
	}
	defer os.Remove(srcFile.Name()) // Clean up after test.

	// Write content to the source file.
	content := []byte("This is the source file content.")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("error writing to source file: %v", err)
	}
	srcFile.Close()

	// Create a temporary destination directory.
	dstDir, err := os.MkdirTemp("", "destination_dir")
	if err != nil {
		t.Fatalf("error creating temp destination directory: %v", err)
	}
	defer os.RemoveAll(dstDir) // Clean up after test.

	// Call the Copy function, copying to the destination directory.
	err = Copy(srcFile.Name(), dstDir)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify the file was copied into the directory.
	dstFilePath := filepath.Join(dstDir, filepath.Base(srcFile.Name()))
	data, err := os.ReadFile(dstFilePath)
	if err != nil {
		t.Fatalf("error reading copied file in directory: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("expected content %s, but got %s", string(content), string(data))
	}
}

// TestCopySourceFileNotExist tests handling when the source file does not exist.
func TestCopySourceFileNotExist(t *testing.T) {
	// Call the Copy function with a non-existent source file.
	err := Copy("nonexistent_file.txt", "destination.txt")
	if err == nil {
		t.Fatalf("expected error when source file does not exist, but got nil")
	}
}

// TestCopyDestinationDirectoryNotExist tests handling when the destination directory does not exist.
func TestCopyDestinationDirectoryNotExist(t *testing.T) {
	// Create a temporary source file.
	srcFile, err := os.CreateTemp("", "source_file.txt")
	if err != nil {
		t.Fatalf("error creating temp source file: %v", err)
	}
	defer os.Remove(srcFile.Name()) // Clean up after test.

	// Write content to the source file.
	content := []byte("This is the source file content.")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("error writing to source file: %v", err)
	}
	srcFile.Close()

	// Call the Copy function with a non-existent destination directory.
	err = Copy(srcFile.Name(), "/nonexistent/destination/file.txt")
	if err == nil {
		t.Fatalf("expected error when destination directory does not exist, but got nil")
	}
}
