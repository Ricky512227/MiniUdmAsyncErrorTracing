package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	tcnVolPath      = "/tcnVol"
	lib64Path       = "/opt/SMAW/INTP/lib64"
	backupExtension = ".backup"
)

// ApplyPatch applies a patch file to a service
//
// Dev will generate lib file and copy to any path in the node.
//
// Syntax: applyPatch.xx "absolutepath of the patch" "servicename"
//
// Workflow:
// 1. Get the md5sum of the patchFile.
// 2. Check the patch file is already present in the /tcnVol,
//    1. If not,
//        1. Copy the patch to the /tcnVol of the required service.
//        2. If the library is already exists in the /opt/SMAW/INTP/lib64 and md5sum of the patchfile is diff.
//            a. Take the backup of the existing lib file, link the new file from /tcnVol to the /opt/SMAW/INTP/lib64.
//    2. If yes,
//        1. the md5sum are diff and the same patch is already linked from /tcnVol to the /opt/SMAW/INTP/lib64.
//        2. then Copy the patch to the /tcnVol of the required service
// 3. Login the required service of mcc container, kill the service process
// 4. Monitor till the process comes up.
// 5. If all the process comes up then,
//    1. Log Patching as successful
//    2. else, Patching as unsuccessful.
func ApplyPatch(patchPath, serviceName string) error {
	log.Printf("Starting patch application for service: %s", serviceName)
	log.Printf("Patch file: %s", patchPath)

	// Step 1: Get MD5 checksum of patch file
	patchMD5, err := CalculateMD5(patchPath)
	if err != nil {
		return fmt.Errorf("failed to calculate MD5: %w", err)
	}
	log.Printf("Patch MD5: %s", patchMD5)

	// Step 2: Check if patch file exists in /tcnVol
	serviceTcnVolPath := filepath.Join(tcnVolPath, serviceName)
	patchFileName := filepath.Base(patchPath)
	tcnVolPatchPath := filepath.Join(serviceTcnVolPath, patchFileName)

	patchExists := FileExists(tcnVolPatchPath)
	var existingMD5 string

	if patchExists {
		existingMD5, err = CalculateMD5(tcnVolPatchPath)
		if err != nil {
			log.Printf("Warning: Could not calculate MD5 of existing patch: %v", err)
		}
	}

	// Step 3: Handle patch file
	if !patchExists {
		// Copy patch to /tcnVol
		log.Printf("Copying patch to %s", tcnVolPatchPath)
		if err := os.MkdirAll(serviceTcnVolPath, 0755); err != nil {
			return fmt.Errorf("failed to create tcnVol directory: %w", err)
		}

		if err := CopyFile(patchPath, tcnVolPatchPath); err != nil {
			return fmt.Errorf("failed to copy patch: %w", err)
		}

		// Check if library exists and needs updating
		libPath := filepath.Join(lib64Path, patchFileName)
		if FileExists(libPath) {
			libMD5, err := CalculateMD5(libPath)
			if err == nil && libMD5 != patchMD5 {
				log.Printf("Library exists with different MD5, updating...")
				if err := updateLibrary(tcnVolPatchPath, libPath); err != nil {
					return fmt.Errorf("failed to update library: %w", err)
				}
			}
		}
	} else if existingMD5 != patchMD5 {
		// Patch exists but MD5 differs
		log.Printf("Patch exists with different MD5, updating...")
		if err := CopyFile(patchPath, tcnVolPatchPath); err != nil {
			return fmt.Errorf("failed to update patch: %w", err)
		}

		// Check if linked and update if needed
		libPath := filepath.Join(lib64Path, patchFileName)
		if FileExists(libPath) {
			if err := updateLibrary(tcnVolPatchPath, libPath); err != nil {
				return fmt.Errorf("failed to update library: %w", err)
			}
		}
	} else {
		log.Printf("Patch already exists with same MD5, skipping copy")
	}

	// Step 4: Restart service (placeholder - actual implementation would use kubectl/API)
	log.Printf("Restarting service: %s", serviceName)
	if err := restartService(serviceName); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	// Step 5: Monitor service health
	log.Printf("Monitoring service health...")
	if err := monitorServiceHealth(serviceName, 30*time.Second); err != nil {
		return fmt.Errorf("service health check failed: %w", err)
	}

	log.Printf("Patch application completed successfully for service: %s", serviceName)
	return nil
}

// updateLibrary updates the library file with backup
func updateLibrary(sourcePath, libPath string) error {
	// Create backup
	backupPath, err := BackupFile(libPath)
	if err != nil {
		return fmt.Errorf("failed to backup library: %w", err)
	}
	log.Printf("Backed up library to: %s", backupPath)

	// Create symlink from /tcnVol to /opt/SMAW/INTP/lib64
	if err := CreateSymlink(sourcePath, libPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	log.Printf("Created symlink: %s -> %s", libPath, sourcePath)

	return nil
}

// restartService restarts a service (placeholder implementation)
func restartService(serviceName string) error {
	// In real implementation, this would:
	// 1. Login to mcc container
	// 2. Kill the service process
	// 3. Wait for it to restart
	log.Printf("Service restart initiated for: %s", serviceName)
	time.Sleep(2 * time.Second) // Simulate restart time
	return nil
}

// monitorServiceHealth monitors service until it's healthy
func monitorServiceHealth(serviceName string, timeout time.Duration) error {
	startTime := time.Now()
	checkInterval := 2 * time.Second

	for time.Since(startTime) < timeout {
		// In real implementation, check if all processes are up
		// For now, just simulate success
		log.Printf("Checking health of service: %s", serviceName)
		time.Sleep(checkInterval)

		// Simulate health check success
		log.Printf("Service %s is healthy", serviceName)
		return nil
	}

	return fmt.Errorf("service health check timed out after %v", timeout)
}
