# Bazel Build System

This project uses [Bazel](https://bazel.build/) as the build system for reproducible, fast builds.

## Why Bazel?

- **Reproducible builds**: Same inputs always produce the same outputs
- **Fast incremental builds**: Only rebuilds what changed
- **Multi-language support**: Easy to add other languages later
- **Remote execution support**: Can scale builds across machines
- **Dependency management**: Explicit dependency declarations

## Getting Started

### Installation

1. **Install Bazelisk** (recommended):
   ```bash
   # macOS
   brew install bazelisk
   
   # Linux
   wget https://github.com/bazelbuild/bazelisk/releases/download/v1.19.0/bazelisk-linux-amd64
   chmod +x bazelisk-linux-amd64
   sudo mv bazelisk-linux-amd64 /usr/local/bin/bazel
   
   # Or download from https://github.com/bazelbuild/bazelisk/releases
   ```

2. **Verify installation**:
   ```bash
   bazel version
   ```

### Basic Commands

```bash
# Build all binaries
bazel build //cmd/list-deployments:list-deployments
bazel build //cmd/symptom-collection:symptom-collection
bazel build //cmd/apply-patch:apply-patch

# Or build all at once
bazel build //cmd/...

# Run tests
bazel test //...

# Run a binary
bazel run //cmd/list-deployments:list-deployments -- -n default

# Clean build artifacts
bazel clean
```

## Project Structure

```
.
├── MODULE.bazel          # Bazel module definition (dependencies)
├── BUILD.bazel           # Root BUILD file (contains Gazelle)
├── .bazelrc              # Bazel configuration
├── .bazelversion         # Bazel version pin
├── pkg/
│   ├── BUILD.bazel       # Package-level BUILD file
│   └── ...
└── cmd/
    ├── BUILD.bazel       # Command-level BUILD file
    └── ...
```

## Dependency Management

Dependencies are managed via `go.mod` and automatically imported into Bazel using Gazelle's `go_deps` extension.

### Adding a Go Dependency

1. Add to `go.mod`:
   ```bash
   go get github.com/example/package@v1.0.0
   go mod tidy
   ```

2. Update BUILD files with Gazelle:
   ```bash
   bazel run //:gazelle
   ```

Gazelle will automatically:
- Read dependencies from `go.mod`
- Update BUILD files with correct dependency labels
- Resolve transitive dependencies

### Manual Dependency Updates

If you need to manually update dependencies:

```bash
# Sync Bazel dependencies
bazel sync

# Update BUILD files
bazel run //:gazelle
```

## BUILD Files

BUILD files define build targets. Example:

```starlark
load("@rules_go//go:def.bzl", "go_library", "go_binary")

go_library(
    name = "utils",
    srcs = ["file.go", "hash.go"],
    importpath = "github.com/Ricky512227/MiniUdmAsyncErrorTracing/pkg/utils",
    visibility = ["//visibility:public"],
)
```

### Using Gazelle

[Gazelle](https://github.com/bazelbuild/bazel-gazelle) auto-generates and updates BUILD files based on your Go code:

```bash
# Update all BUILD files
bazel run //:gazelle

# Or use the Make target
make gazelle
```

Gazelle will:
- Scan `.go` files for imports
- Generate BUILD files
- Update dependencies from `go.mod`
- Resolve import paths to Bazel labels

## Configuration

### .bazelrc

Contains Bazel configuration:

- Build modes (fastbuild, opt, dbg)
- Platform settings
- Remote execution settings
- Test configuration

### .bazelversion

Pins the Bazel version (used by Bazelisk):

```
7.0.0
```

### MODULE.bazel

Defines the Bazel module and uses Gazelle's `go_deps` extension to import dependencies from `go.mod`.

## Building for Different Platforms

```bash
# Linux
bazel build //cmd/list-deployments:list-deployments --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64

# macOS
bazel build //cmd/list-deployments:list-deployments --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64

# Windows
bazel build //cmd/list-deployments:list-deployments --platforms=@io_bazel_rules_go//go/toolchain:windows_amd64
```

## Testing

```bash
# Run all tests
bazel test //...

# Run specific test
bazel test //pkg/utils:utils_test

# Run with coverage
bazel test --collect_code_coverage //...

# View coverage
find bazel-testlogs -name "coverage.dat"
```

## Troubleshooting

### Build failures

1. **Clean and rebuild**:
   ```bash
   bazel clean
   bazel build //...
   ```

2. **Update BUILD files**:
   ```bash
   bazel run //:gazelle
   ```

3. **Sync dependencies**:
   ```bash
   bazel sync
   ```

### Cache issues

Clear Bazel cache:
```bash
bazel clean --expunge
```

### BUILD file issues

Regenerate with Gazelle:
```bash
bazel run //:gazelle
```

### Dependency resolution

If dependencies aren't resolving correctly:

1. Ensure `go.mod` is up to date:
   ```bash
   go mod tidy
   ```

2. Update BUILD files:
   ```bash
   bazel run //:gazelle
   ```

3. Sync Bazel:
   ```bash
   bazel sync
   ```

## CI/CD Integration

The GitHub Actions workflow uses Bazel for building and testing. See `.github/workflows/ci.yml`.

## Workflow

Typical development workflow:

1. **Add/update Go code**
2. **Update dependencies** (if needed):
   ```bash
   go get package@version
   go mod tidy
   ```
3. **Update BUILD files**:
   ```bash
   bazel run //:gazelle
   ```
4. **Build and test**:
   ```bash
   bazel build //...
   bazel test //...
   ```

## Resources

- [Bazel Documentation](https://bazel.build/docs)
- [Rules Go](https://github.com/bazelbuild/rules_go)
- [Gazelle](https://github.com/bazelbuild/bazel-gazelle)
- [Bazel Best Practices](https://bazel.build/contribute/guide)
