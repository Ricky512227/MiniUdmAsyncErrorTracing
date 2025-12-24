# Getting Started

This guide will help you get started with MiniUdm Async Error Tracing.

## Quick Start

### Prerequisites

1. **Go 1.21 or higher**
   ```bash
   go version
   ```

2. **Kubernetes cluster access**
   - Ensure `kubectl` is configured
   ```bash
   kubectl cluster-info
   ```

3. **Make** (optional, for using Makefile)
   ```bash
   make --version
   ```

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ricky512227/MiniUdmAsyncErrorTracing.git
   cd MiniUdmAsyncErrorTracing
   ```

2. **Install dependencies**
   ```bash
   make deps
   # or
   go mod download
   ```

3. **Build the binaries**
   ```bash
   make build
   # or
   go build -o bin/list-deployments ./cmd/list-deployments
   go build -o bin/symptom-collection ./cmd/symptom-collection
   go build -o bin/apply-patch ./cmd/apply-patch
   ```

## Configuration

### Using Configuration File

Copy the example configuration:

```bash
cp configs/config.yaml ~/.miniumd/config.yaml
```

Edit the configuration file to match your environment:

```yaml
kubernetes:
  namespace: "your-namespace"

paths:
  tcn_vol_path: "/tcnVol"
  lib64_path: "/opt/SMAW/INTP/lib64"

logging:
  level: "info"
```

### Using Environment Variables

Set environment variables (prefixed with `MINIUDM_`):

```bash
export MINIUDM_KUBERNETES_NAMESPACE="your-namespace"
export MINIUDM_LOGGING_LEVEL="debug"
```

### Command-Line Flags

Override configuration with command-line flags:

```bash
./bin/list-deployments -n your-namespace -c /path/to/config.yaml
```

## Usage Examples

### List Deployments

```bash
# List deployments in default namespace
./bin/list-deployments

# List deployments in specific namespace
./bin/list-deployments -n production

# With custom config
./bin/list-deployments -n production -c ~/.miniumd/config.yaml
```

### Symptom Collection

```bash
# Collect symptoms from specific pods
./bin/symptom-collection -n default -p "uecm nim"

# Multiple pods
./bin/symptom-collection -n miniudm -p "uecm testclient"
```

### Apply Patch

```bash
# Apply a patch to a service
./bin/apply-patch -p /path/to/patch.so -s uecm

# With custom config
./bin/apply-patch -p /path/to/patch.so -s uecm -c ~/.miniumd/config.yaml
```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover
```

### Code Quality

```bash
# Format code
make format

# Check code quality
make lint

# Run all checks
make verify
```

### Using the Library Programmatically

See `examples/list_deployments.go` for an example of using the library:

```go
package main

import (
    "context"
    "github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/kubernetes"
)

func main() {
    ctx := context.Background()
    client, err := kubernetes.NewClient(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    deployments, err := client.GetDeployments("default")
    // ...
}
```

## Troubleshooting

### Common Issues

1. **Kubernetes connection error**
   - Verify `kubectl` is configured correctly
   - Check cluster connectivity: `kubectl cluster-info`

2. **Permission denied**
   - Ensure you have proper RBAC permissions
   - Check namespace access: `kubectl get namespaces`

3. **Configuration not found**
   - Use `-c` flag to specify config path
   - Or set environment variables

4. **Build errors**
   - Ensure Go 1.21+ is installed
   - Run `go mod tidy` to fix dependencies

## Next Steps

- Read the [Architecture Documentation](ARCHITECTURE.md)
- Check the [Contributing Guide](../CONTRIBUTING.md)
- Review the [README](../README.md)

