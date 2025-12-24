package utils

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "existing file",
			filePath: tmpFile.Name(),
			expected: true,
		},
		{
			name:     "non-existing file",
			filePath: "/nonexistent/file/path",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FileExists(tt.filePath)
			if result != tt.expected {
				t.Errorf("FileExists() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	// Create source file
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	content := "test content"
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}

	// Verify destination exists
	if !FileExists(dst) {
		t.Error("Destination file was not created")
	}

	// Verify content
	destContent, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(destContent) != content {
		t.Errorf("Copied content = %s, want %s", string(destContent), content)
	}
}

func TestCalculateMD5(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "test content"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}
	tmpFile.Close()

	hash, err := CalculateMD5(tmpFile.Name())
	if err != nil {
		t.Fatalf("CalculateMD5() error = %v", err)
	}

	if hash == "" {
		t.Error("CalculateMD5() returned empty hash")
	}

	// Test with non-existent file
	_, err = CalculateMD5("/nonexistent/file")
	if err == nil {
		t.Error("CalculateMD5() should return error for non-existent file")
	}
}

func TestGetAge(t *testing.T) {
	// This is a basic test - in real scenario, you'd use a fixed time
	// For now, just verify it returns a non-empty string
	result := GetAge(time.Now().Add(-5 * time.Minute))
	if result == "" {
		t.Error("GetAge() returned empty string")
	}
}

