package patch

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/config"
	"github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/utils"
	"go.uber.org/zap"
)

// Manager handles patch application operations
type Manager struct {
	config  *config.Config
	logger  *zap.Logger
	service ServiceRestarter
}

// ServiceRestarter interface for restarting services
type ServiceRestarter interface {
	RestartService(ctx context.Context, serviceName string) error
	MonitorHealth(ctx context.Context, serviceName string, timeout time.Duration) error
}

// NewManager creates a new patch manager
func NewManager(cfg *config.Config, logger *zap.Logger, restarter ServiceRestarter) *Manager {
	return &Manager{
		config:  cfg,
		logger:  logger,
		service: restarter,
	}
}

// ApplyPatch applies a patch file to a service
//
// Dev will generate lib file and copy to any path in the node.
//
// Syntax: apply-patch -p "absolutepath of the patch" -s "servicename"
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
func (m *Manager) ApplyPatch(ctx context.Context, patchPath, serviceName string) error {
	m.logger.Info("Starting patch application",
		zap.String("service", serviceName),
		zap.String("patch", patchPath),
	)

	// Step 1: Validate and calculate MD5
	if !utils.FileExists(patchPath) {
		return fmt.Errorf("patch file does not exist: %s", patchPath)
	}

	patchMD5, err := utils.CalculateMD5(patchPath)
	if err != nil {
		return fmt.Errorf("failed to calculate MD5: %w", err)
	}
	m.logger.Debug("Calculated patch MD5", zap.String("md5", patchMD5))

	// Step 2: Handle patch file in /tcnVol
	serviceTcnVolPath := filepath.Join(m.config.Paths.TcnVolPath, serviceName)
	patchFileName := filepath.Base(patchPath)
	tcnVolPatchPath := filepath.Join(serviceTcnVolPath, patchFileName)
	libPath := filepath.Join(m.config.Paths.Lib64Path, patchFileName)

	patchExists := utils.FileExists(tcnVolPatchPath)

	if !patchExists {
		// Copy patch to /tcnVol
		m.logger.Info("Copying patch to tcnVol", zap.String("destination", tcnVolPatchPath))
		if err := utils.EnsureDirectory(serviceTcnVolPath); err != nil {
			return fmt.Errorf("failed to create tcnVol directory: %w", err)
		}

		if err := utils.CopyFile(patchPath, tcnVolPatchPath); err != nil {
			return fmt.Errorf("failed to copy patch: %w", err)
		}

		// Update library if it exists and MD5 differs
		if utils.FileExists(libPath) {
			if err := m.updateLibrary(tcnVolPatchPath, libPath); err != nil {
				return fmt.Errorf("failed to update library: %w", err)
			}
		}
	} else {
		// Patch exists, check MD5
		existingMD5, err := utils.CalculateMD5(tcnVolPatchPath)
		if err != nil {
			m.logger.Warn("Could not calculate MD5 of existing patch", zap.Error(err))
		}

		if existingMD5 != patchMD5 {
			m.logger.Info("Patch MD5 differs, updating", zap.String("existing", existingMD5), zap.String("new", patchMD5))
			
			if err := utils.CopyFile(patchPath, tcnVolPatchPath); err != nil {
				return fmt.Errorf("failed to update patch: %w", err)
			}

			// Update library if it exists
			if utils.FileExists(libPath) {
				if err := m.updateLibrary(tcnVolPatchPath, libPath); err != nil {
					return fmt.Errorf("failed to update library: %w", err)
				}
			}
		} else {
			m.logger.Info("Patch already exists with same MD5, skipping copy")
		}
	}

	// Step 3: Restart service
	m.logger.Info("Restarting service", zap.String("service", serviceName))
	if err := m.service.RestartService(ctx, serviceName); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	// Step 4: Monitor service health
	m.logger.Info("Monitoring service health", zap.String("service", serviceName))
	timeout := m.config.Patch.HealthTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	if err := m.service.MonitorHealth(ctx, serviceName, timeout); err != nil {
		return fmt.Errorf("service health check failed: %w", err)
	}

	m.logger.Info("Patch application completed successfully", zap.String("service", serviceName))
	return nil
}

// updateLibrary updates the library file with backup
func (m *Manager) updateLibrary(sourcePath, libPath string) error {
	// Create backup if enabled
	if m.config.Patch.BackupEnabled {
		backupPath, err := utils.BackupFile(libPath)
		if err != nil {
			m.logger.Warn("Failed to backup library", zap.Error(err))
		} else {
			m.logger.Info("Backed up library", zap.String("backup", backupPath))
		}
	}

	// Create symlink from /tcnVol to /opt/SMAW/INTP/lib64
	if err := utils.CreateSymlink(sourcePath, libPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	m.logger.Info("Created symlink", zap.String("link", libPath), zap.String("target", sourcePath))
	return nil
}

