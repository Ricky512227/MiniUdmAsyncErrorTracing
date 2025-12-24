# MiniUdm Async Error Tracing

A Kubernetes-based tool for tracing and monitoring errors in MiniUdm deployments. This tool helps collect symptoms, apply patches, and monitor async error patterns in Kubernetes clusters.

## Features

- **Deployment Monitoring**: List and monitor Kubernetes deployments
- **Error Tracing**: Collect async error traces from pods
- **Patch Management**: Apply patches to services with validation
- **Symptom Collection**: Monitor logs and collect error symptoms

## Prerequisites

- Go 1.18 or higher
- Kubernetes cluster access
- kubectl configured with cluster access
- Kubernetes client-go libraries

## Installation

```bash
git clone https://github.com/Ricky512227/MiniUdmAsyncErrorTracing.git
cd MiniUdmAsyncErrorTracing
go mod download
go build
```

## Usage

### Sample Deployment Listing

```bash
go run sample.go
```

This will list all deployments in the default namespace.

### Get Deployments Programmatically

```go
import (
    "context"
    "github.com/Ricky512227/MiniUdmAsyncErrorTracing"
)

deployments, err := GetDeployments(clientset, ctx, "default")
```

## Project Structure

- `sample.go` - Example code for listing deployments
- `startSymptionCollection.go` - Symptom collection workflow (planned)
- `applyPatch.go` - Patch application logic (planned)
- `applyImage.go` - Image application logic (planned)
- `commonUtility.go` - Common utility functions (planned)

## Planned Features

### Symptom Collection Workflow

The symptom collection process will:
1. Enable traces for processes
2. Enable pcap capture
3. Execute test commands
4. Monitor log files for errors:
   - `/cmconfig.log`
   - `/logstore/TspCore`
   - `/RTPTraceError`
   - `/Envoy`
   - `/dumplog`
5. Collect and store traces for analysis

### Patch Application

The patch application process will:
1. Validate patch file (MD5 checksum)
2. Copy patch to `/tcnVol`
3. Link library files to `/opt/SMAW/INTP/lib64`
4. Restart service processes
5. Monitor process health
6. Log success/failure status

## Development

```bash
# Run tests
go test ./...

# Build
go build

# Format code
go fmt ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available for use.

## Author

Ricky512227

