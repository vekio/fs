package file

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCreate tests the Create function to ensure it creates a new file if it doesn't exist.
func TestCreate(t *testing.T) {
	// Create a temporary directory to store the test file.
	tmpDir, err := os.MkdirTemp("", "test_create_file_dir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the directory and its content after the test.

	// Define the path to a non-existing file inside the temp directory.
	filePath := filepath.Join(tmpDir, "test_create_file.txt")

	// Ensure the file does not exist initially.
	if _, err := os.Stat(filePath); err == nil {
		t.Fatalf("expected file to not exist, but it already exists")
	}

	// Call CreateFile to create the file.
	err = CreateFile(filePath, DefaultFilePerms)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Check if the file was created.
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("expected file to be created, but got error: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("expected a file, but got a directory")
	}
}

// TestTruncate tests the Create function to ensure it truncates an existing file.
func TestTruncate(t *testing.T) {
	// Create a temporary file.
	tmpFile, err := os.CreateTemp("", "test_truncate_file.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test.

	// Write some initial content to the file.
	initialContent := []byte("This is some initial content.")
	if err := os.WriteFile(tmpFile.Name(), initialContent, 0644); err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify the file contains the initial content.
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected to read file, but got error: %v", err)
	}
	if string(data) != string(initialContent) {
		t.Fatalf("expected content %s, but got %s", string(initialContent), string(data))
	}

	// Call Create to truncate the file.
	err = CreateFile(tmpFile.Name(), 0644)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify that the file has been truncated (should be empty).
	data, err = os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected to read file, but got error: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("expected file to be empty after truncation, but it has content: %s", string(data))
	}
}

// TestFileExists tests the FileExists function.
func TestFileExists(t *testing.T) {
	// Create a temporary directory to store the test file.
	tmpDir, err := os.MkdirTemp("", "test_exists_dir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the directory and its content after the test.

	// Define the path to a non-existing file inside the temp directory.
	filePath := filepath.Join(tmpDir, "test_exists.txt")

	// Initially, the file should not exist.
	exists, err := FileExists(filePath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if exists {
		t.Fatalf("expected file to not exist, but it does")
	}

	// Create the file.
	tmpFile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	tmpFile.Close()

	// Now the file should exist.
	exists, err = FileExists(filePath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if !exists {
		t.Fatalf("expected file to exist, but it does not")
	}
}

// TestAppendToFile tests the AppendToFile function by appending content to a file.
func TestAppendToFile(t *testing.T) {
	// Create a temporary file.
	tmpFile, err := os.CreateTemp("", "test_append.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write initial content.
	initialContent := []byte("Initial content.\n")
	if err := os.WriteFile(tmpFile.Name(), initialContent, 0644); err != nil {
		t.Fatalf("failed to write initial content: %v", err)
	}

	// Append new content.
	appendContent := []byte("Appended content.\n")
	err = AppendToFile(tmpFile.Name(), appendContent)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify the content.
	expectedContent := "Initial content.\nAppended content.\n"
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != expectedContent {
		t.Fatalf("expected content: %s, but got: %s", expectedContent, string(data))
	}
}

// TestGetFileSize tests the GetFileSize function to ensure it returns the correct file size.
func TestGetFileSize(t *testing.T) {
	// Create a temporary file.
	tmpFile, err := os.CreateTemp("", "test_size.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the file.
	content := []byte("This is test content.\n")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}
	tmpFile.Close()

	// Get the file size.
	size, err := GetFileSize(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify the size.
	expectedSize := int64(len(content))
	if size != expectedSize {
		t.Fatalf("expected size: %d, but got: %d", expectedSize, size)
	}
}

// TestMoveFile tests the MoveFile function to ensure it moves a file correctly.
func TestMoveFile(t *testing.T) {
	// Create a temporary source file.
	srcFile, err := os.CreateTemp("", "test_move_src.txt")
	if err != nil {
		t.Fatalf("failed to create temp source file: %v", err)
	}
	defer os.Remove(srcFile.Name())

	// Write content to the source file.
	content := []byte("This is test content.\n")
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}
	srcFile.Close()

	// Create a temporary destination file path.
	dstFile, err := os.CreateTemp("", "test_move_dst.txt")
	if err != nil {
		t.Fatalf("failed to create temp destination file: %v", err)
	}
	dstFile.Close()
	os.Remove(dstFile.Name()) // Remove the empty file, we just want the path.

	// Move the source file to the destination.
	err = MoveFile(srcFile.Name(), dstFile.Name())
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Check if the source file is gone and the destination file exists.
	if _, err := os.Stat(srcFile.Name()); !os.IsNotExist(err) {
		t.Fatalf("expected source file to be moved, but it still exists")
	}
	if _, err := os.Stat(dstFile.Name()); err != nil {
		t.Fatalf("expected destination file to exist, but got: %v", err)
	}

	// Verify the content of the destination file.
	data, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("expected content: %s, but got: %s", string(content), string(data))
	}
}

// TestReadFile tests the ReadFile function to ensure it reads the content correctly.
func TestReadFile(t *testing.T) {
	// Create a temporary file.
	tmpFile, err := os.CreateTemp("", "test_read.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to the file.
	content := []byte("This is test content.\n")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}
	tmpFile.Close()

	// Use the ReadFile function to read the file content.
	data, err := ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Verify the content.
	if string(data) != string(content) {
		t.Fatalf("expected content: %s, but got: %s", string(content), string(data))
	}
}

// TestListFile tests the ListFile function to ensure it lists only files.
func TestListFile(t *testing.T) {
	// Create a temporary directory.
	tmpDir, err := os.MkdirTemp("", "test_list_file_dir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directories and files inside the temporary directory.
	directories := []string{"dir1", "dir2", "dir3"}
	for _, d := range directories {
		if err := os.Mkdir(filepath.Join(tmpDir, d), 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", d, err)
		}
	}
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, f := range files {
		if _, err := os.Create(filepath.Join(tmpDir, f)); err != nil {
			t.Fatalf("failed to create file %s: %v", f, err)
		}
	}

	// Call ListFile and check the result.
	listedFiles, err := ListFile(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if len(listedFiles) != len(files) {
		t.Fatalf("expected %d files, but got %d", len(files), len(listedFiles))
	}

	// Verify that the correct files were listed.
	for _, f := range files {
		found := false
		for _, listedFile := range listedFiles {
			if listedFile == f {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected file %s to be listed, but it was not", f)
		}
	}
}
