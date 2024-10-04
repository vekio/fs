package dir

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCreateDir tests the CreateDir function to ensure it creates a directory correctly.
func TestCreateDir(t *testing.T) {
	// Define a temporary directory path.
	dirPath := "test_create_dir"
	defer os.RemoveAll(dirPath) // Clean up after the test.

	// Call CreateDir to create the directory.
	err := CreateDir(dirPath, DefaultDirPerms)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Check if the directory was created.
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("expected directory to be created, but got error: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected a directory, but got a file")
	}
}

// TestDirExists tests the DirExists function to ensure it correctly identifies if a directory exists.
func TestDirExists(t *testing.T) {
	// Define a temporary directory path.
	dirPath := "test_exists_dir"
	defer os.RemoveAll(dirPath) // Clean up after the test.

	// Directory should not exist initially.
	exists, err := DirExists(dirPath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if exists {
		t.Fatalf("expected directory to not exist, but it does")
	}

	// Create the directory.
	err = os.Mkdir(dirPath, DefaultDirPerms)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Now the directory should exist.
	exists, err = DirExists(dirPath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if !exists {
		t.Fatalf("expected directory to exist, but it does not")
	}
}

// TestEnsureDir tests the EnsureDir function to ensure it creates a directory if it doesn't exist.
func TestEnsureDir(t *testing.T) {
	// Define a temporary directory path.
	dirPath := "test_ensure_dir"
	defer os.RemoveAll(dirPath) // Clean up after the test.

	// Directory should not exist initially.
	exists, err := DirExists(dirPath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if exists {
		t.Fatalf("expected directory to not exist, but it does")
	}

	// Call EnsureDir to create the directory.
	err = EnsureDir(dirPath, DefaultDirPerms)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// Now the directory should exist.
	exists, err = DirExists(dirPath)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if !exists {
		t.Fatalf("expected directory to exist, but it does not")
	}

	// Call EnsureDir again, it should not return an error.
	err = EnsureDir(dirPath, DefaultDirPerms)
	if err != nil {
		t.Fatalf("expected no error when ensuring existing directory, but got: %v", err)
	}
}

// TestListDir tests the ListDir function to ensure it lists only directories.
func TestListDir(t *testing.T) {
	// Create a temporary directory.
	tmpDir, err := os.MkdirTemp("", "test_list_dir")
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
	files := []string{"file1.txt", "file2.txt"}
	for _, f := range files {
		if _, err := os.Create(filepath.Join(tmpDir, f)); err != nil {
			t.Fatalf("failed to create file %s: %v", f, err)
		}
	}

	// Call ListDir and check the result.
	dirs, err := ListDir(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if len(dirs) != len(directories) {
		t.Fatalf("expected %d directories, but got %d", len(directories), len(dirs))
	}

	// Verify that the correct directories were listed.
	for _, d := range directories {
		found := false
		for _, listedDir := range dirs {
			if listedDir == d {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected directory %s to be listed, but it was not", d)
		}
	}
}

// TestIsEmptyDir tests the IsEmptyDir function to ensure it detects empty directories correctly.
func TestIsEmptyDir(t *testing.T) {
	// Create a temporary directory.
	tmpDir, err := os.MkdirTemp("", "test_empty_dir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// The directory should be empty.
	empty, err := IsEmptyDir(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if !empty {
		t.Fatalf("expected directory to be empty, but it is not")
	}

	// Create a file in the directory.
	tmpFile, err := os.CreateTemp(tmpDir, "test_file.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	// Now the directory should not be empty.
	empty, err = IsEmptyDir(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if empty {
		t.Fatalf("expected directory to not be empty, but it is")
	}
}
