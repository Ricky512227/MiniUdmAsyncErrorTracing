# MiniUdm Async Error Tracing

[![CI](https://github.com/Ricky512227/MiniUdmAsyncErrorTracing/workflows/CI/badge.svg)](https://github.com/Ricky512227/MiniUdmAsyncErrorTracing/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ricky512227/MiniUdmAsyncErrorTracing)](https://goreportcard.com/report/github.com/Ricky512227/MiniUdmAsyncErrorTracing)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A production-grade Kubernetes-based tool for tracing and monitoring errors in MiniUdm deployments. This tool helps collect symptoms, apply patches, and monitor async error patterns in Kubernetes clusters.

## Features

- **Deployment Monitoring**: List and monitor Kubernetes deployments with detailed status
- **Error Tracing**: Collect async error traces from pods with configurable keyword detection
- **Patch Management**: Apply patches to services with MD5 validation and health monitoring
- **Symptom Collection**: Monitor multiple log files and collect error symptoms in parallel
- **Production-Ready**: Comprehensive error handling, logging, configuration management, and CI/CD

## Architecture

The project follows Go best practices with a clear separation of concerns:

```
.
├── cmd/                    # Command-line applications
│   ├── list-deployments/  # List Kubernetes deployments
│   ├── symptom-collection/ # Start symptom collection
│   └── apply-patch/       # Apply patches to services
├── pkg/                    # Reusable packages
│   ├── kubernetes/        # Kubernetes client wrapper
│   ├── utils/             # Utility functions (file, hash, command)
│   ├── patch/             # Patch application logic
│   ├── symptom/           # Symptom collection logic
│   └── config/            # Configuration management
├── internal/               # Internal packages (not for external use)
│   ├── logger/            # Logging utilities
│   └── validator/         # Validation utilities
├── configs/                # Configuration files
├── examples/               # Example code
├── docs/                   # Documentation
└── scripts/                # Build and utility scripts
```

## Prerequisites

- Go 1.21 or higher
- Kubernetes cluster access
- kubectl configured with cluster access
- Make (optional, for using Makefile)

## Installation

### Prerequisites

- Bazel 7.0+ (install via [Bazelisk](https://github.com/bazelbuild/bazelisk) or [Bazel](https://bazel.build/install))
- Go 1.21+ (for local development and linting)

### From Source

```bash
git clone https://github.com/Ricky512227/MiniUdmAsyncErrorTracing.git
cd MiniUdmAsyncErrorTracing

# Initial setup: Update BUILD files with Gazelle
bazel run //:gazelle

# Build all binaries
make build
# or
bazel build //cmd/...
```

The binaries will be available in the `bazel-bin/` directory.

**Note**: On first setup, run `bazel run //:gazelle` to ensure BUILD files are properly configured with dependencies from `go.mod`.

### Using Bazel Directly

```bash
# Build all binaries
bazel build //cmd/list-deployments:list-deployments
bazel build //cmd/symptom-collection:symptom-collection
bazel build //cmd/apply-patch:apply-patch

# Run tests
bazel test //...

# Run a specific binary
bazel run //cmd/list-deployments:list-deployments -- -n default
```

### Using Make (Bazel wrapper)

```bash
# Build all binaries
make build

# Build specific binary
make build-list
make build-symptom
make build-patch

# Run tests
make test

# Update BUILD files with Gazelle
make gazelle

# Format code (Go)
make format

# Run linter
make lint
```

## Configuration

The tool can be configured via:

1. **Configuration file** (`configs/config.yaml`)
2. **Environment variables** (prefixed with `MINIUDM_`)
3. **Command-line flags**

### Example Configuration

See `configs/config.yaml` for a complete example:

```yaml
kubernetes:
  namespace: "default"
  timeout: "30s"

paths:
  tcn_vol_path: "/tcnVol"
  lib64_path: "/opt/SMAW/INTP/lib64"

logging:
  level: "info"
  format: "json"

symptom:
  error_keywords:
    - "error"
    - "ERROR"
    - "fatal"
  check_interval: "1s"
  collection_timeout: "10m"

patch:
  backup_enabled: true
  health_timeout: "30s"
```

## Usage

### List Deployments

List all deployments in a namespace:

```bash
./bin/list-deployments -n default

# Or with custom config
./bin/list-deployments -n my-namespace -c /path/to/config.yaml
```

### Symptom Collection

Start symptom collection for specific pods:

```bash
./bin/symptom-collection -n default -p "uecm nim"

# With custom config
./bin/symptom-collection -n miniudm -p "uecm" -c /path/to/config.yaml
```

The symptom collection process:
1. Enables traces for processes
2. Enables pcap capture
3. Executes test commands (pybot)
4. Monitors log files for errors:
   - `/cmconfig.log`
   - `/logstore/TspCore`
   - `/RTPTraceError`
   - `/Envoy`
   - `/dumplog`
5. Collects and stores traces for analysis

### Apply Patch

Apply a patch file to a service:

```bash
./bin/apply-patch -p /path/to/patch.so -s service-name

# With custom config
./bin/apply-patch -p /path/to/patch.so -s uecm -c /path/to/config.yaml
```

The patch application process:
1. Validates patch file (MD5 checksum)
2. Copies patch to `/tcnVol`
3. Links library files to `/opt/SMAW/INTP/lib64`
4. Restarts service processes
5. Monitors process health
6. Logs success/failure status

## Development

### Running Tests

```bash
# Run all tests with Bazel
bazel test //...
# or
make test

# Run tests with coverage
make test-cover
```

### Code Quality

```bash
# Format code
make format

# Check formatting
make fmt-check

# Run linter
make lint

# Run all checks
make verify

# Update BUILD files after adding dependencies
make gazelle
```

### Project Structure Guidelines

- **cmd/**: Each subdirectory contains a main package for a command-line tool
- **pkg/**: Public packages that can be imported by other projects
- **internal/**: Private packages specific to this project
- **configs/**: Configuration file templates
- **examples/**: Example code demonstrating usage

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Format code (`make format`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Code Style

- Follow Go code review comments: https://github.com/golang/go/wiki/CodeReviewComments
- Use `gofmt` for formatting
- Add comments for exported functions and types
- Write tests for new features

## CI/CD

The project uses GitHub Actions for continuous integration:

- **Test**: Runs tests on every push and pull request
- **Lint**: Validates code quality using golangci-lint
- **Build**: Builds binaries for multiple platforms

## License

This project is open source and available for use. See [LICENSE](LICENSE) file for details.

## Author

Ricky512227

## Acknowledgments

Built with:
- [Kubernetes client-go](https://github.com/kubernetes/client-go)
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Zap](https://github.com/uber-go/zap) - Logging

