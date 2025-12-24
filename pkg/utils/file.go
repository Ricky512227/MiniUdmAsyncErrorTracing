package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// FileExists checks if a file exists at the given path
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Sync to ensure data is written
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// CreateSymlink creates a symbolic link from target to linkPath
func CreateSymlink(target, linkPath string) error {
	// Remove existing link if it exists
	if FileExists(linkPath) {
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to remove existing link: %w", err)
		}
	}

	// Create parent directory if needed
	if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.Symlink(target, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// BackupFile creates a backup of a file with timestamp
func BackupFile(filePath string) (string, error) {
	if !FileExists(filePath) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup.%s", filePath, timestamp)
	
	if err := CopyFile(filePath, backupPath); err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return backupPath, nil
}

// EnsureDirectory ensures a directory exists, creating it if necessary
func EnsureDirectory(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

