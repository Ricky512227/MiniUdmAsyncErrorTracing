# Architecture Documentation

## Overview

MiniUdm Async Error Tracing is built with a modular architecture following Go best practices. The project is organized into clear layers with separation of concerns.

## Directory Structure

```
.
├── cmd/                    # Command-line applications
│   ├── list-deployments/  # List Kubernetes deployments
│   ├── symptom-collection/ # Start symptom collection
│   └── apply-patch/       # Apply patches to services
├── pkg/                    # Reusable packages (public API)
│   ├── kubernetes/        # Kubernetes client wrapper
│   ├── utils/             # Utility functions
│   ├── patch/             # Patch application logic
│   ├── symptom/           # Symptom collection logic
│   └── config/            # Configuration management
├── internal/               # Internal packages (private)
│   ├── logger/            # Logging utilities
│   └── validator/         # Validation utilities
├── configs/                # Configuration files
├── examples/               # Example code
├── docs/                   # Documentation
└── scripts/                # Build and utility scripts
```

## Package Descriptions

### cmd/

Command-line applications that provide user-facing functionality. Each subdirectory contains a `main` package with a Cobra-based CLI.

### pkg/kubernetes

Provides a high-level wrapper around the Kubernetes client-go library. Encapsulates common operations like listing deployments, checking namespace existence, and monitoring deployment readiness.

**Key Types:**
- `Client`: Main client struct wrapping `kubernetes.Clientset`

**Key Functions:**
- `NewClient()`: Creates a new Kubernetes client
- `GetDeployments()`: Lists deployments in a namespace
- `NamespaceExists()`: Checks if a namespace exists
- `IsDeploymentReady()`: Checks deployment readiness

### pkg/utils

Utility functions for common operations.

**Subpackages:**
- `file.go`: File operations (copy, symlink, backup, existence checks)
- `hash.go`: Hash calculation (MD5, SHA256)
- `command.go`: Command execution with timeout support
- `time.go`: Time formatting utilities

### pkg/config

Configuration management using Viper. Supports:
- YAML configuration files
- Environment variables (prefixed with `MINIUDM_`)
- Default values

**Key Types:**
- `Config`: Main configuration struct
- `KubernetesConfig`: Kubernetes-specific settings
- `PathsConfig`: Path configuration
- `SymptomConfig`: Symptom collection settings
- `PatchConfig`: Patch application settings

### pkg/patch

Handles patch application to Kubernetes services. Implements the workflow:
1. Validate patch file (MD5 checksum)
2. Copy to `/tcnVol`
3. Link to `/opt/SMAW/INTP/lib64`
4. Restart service
5. Monitor health

**Key Types:**
- `Manager`: Main patch manager
- `ServiceRestarter`: Interface for service restart operations

### pkg/symptom

Manages symptom collection from Kubernetes pods. Coordinates parallel operations:
- Enabling traces
- Enabling pcap capture
- Executing test commands
- Monitoring log files
- Collecting error events

**Key Types:**
- `Collector`: Main symptom collector
- `SymptomCollectionConfig`: Collection configuration
- `ErrorEvent`: Error event representation

### internal/logger

Internal logging utilities using Zap. Provides structured logging with configurable levels.

## Design Principles

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Dependency Injection**: Dependencies are injected rather than created internally
3. **Interface-Based Design**: Interfaces allow for testability and flexibility
4. **Error Handling**: Comprehensive error handling with proper error wrapping
5. **Configuration Management**: Centralized configuration with multiple sources
6. **Logging**: Structured logging throughout the application

## Data Flow

### List Deployments

```
CLI → Config → Kubernetes Client → Kubernetes API → Response → Output
```

### Symptom Collection

```
CLI → Config → Collector → {
  - Kubernetes Client (validation)
  - Parallel routines (traces, pcap, logs)
  - Error channel (error events)
  - Cleanup (storage)
}
```

### Patch Application

```
CLI → Config → Patch Manager → {
  - File operations (copy, link)
  - Service restarter (restart)
  - Health monitor (verification)
}
```

## Extension Points

### Adding a New Command

1. Create a new directory under `cmd/`
2. Implement `main.go` with Cobra command structure
3. Use packages from `pkg/` for functionality
4. Add build target to Makefile

### Adding a New Package

1. Create directory under `pkg/` (public) or `internal/` (private)
2. Define clear interfaces
3. Add comprehensive tests
4. Document exported types and functions

### Custom Service Restarter

Implement the `patch.ServiceRestarter` interface to customize service restart behavior.

## Testing Strategy

- Unit tests for each package
- Integration tests for CLI commands
- Mock Kubernetes API for testing
- Table-driven tests for utilities

## Future Enhancements

- Metrics collection (Prometheus)
- Distributed tracing support
- Plugin system for custom collectors
- Web UI for monitoring
- Advanced error pattern detection

