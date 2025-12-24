package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// CalculateMD5 calculates MD5 checksum of a file
func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CalculateSHA256 calculates SHA256 checksum of a file
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// VerifyChecksum verifies a file's checksum matches the expected value
func VerifyChecksum(filePath, expectedMD5 string) (bool, error) {
	actualMD5, err := CalculateMD5(filePath)
	if err != nil {
		return false, err
	}
	return actualMD5 == expectedMD5, nil
}

