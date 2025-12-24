.PHONY: help build test clean install lint format vet deps run-list run-symptom run-patch gazelle update-deps

# Bazel parameters
BAZEL=bazel
BAZEL_RUN=$(BAZEL) run
BAZEL_BUILD=$(BAZEL) build
BAZEL_TEST=$(BAZEL) test
BAZEL_CLEAN=$(BAZEL) clean

# Directories
BIN_DIR=./bazel-bin

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

deps: ## Update Bazel dependencies
	$(BAZEL) sync --only=com_github_spf13_cobra,com_github_spf13_viper,go_uber_org_zap

update-deps: ## Update all Bazel dependencies using Gazelle
	$(BAZEL) run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

gazelle: ## Run Gazelle to update BUILD files
	$(BAZEL) run //:gazelle

build: ## Build all binaries with Bazel
	$(BAZEL_BUILD) //cmd/list-deployments:list-deployments
	$(BAZEL_BUILD) //cmd/symptom-collection:symptom-collection
	$(BAZEL_BUILD) //cmd/apply-patch:apply-patch
	@echo "Build complete. Binaries are in $(BIN_DIR)/"

build-list: ## Build list-deployments binary
	$(BAZEL_BUILD) //cmd/list-deployments:list-deployments

build-symptom: ## Build symptom-collection binary
	$(BAZEL_BUILD) //cmd/symptom-collection:symptom-collection

build-patch: ## Build apply-patch binary
	$(BAZEL_BUILD) //cmd/apply-patch:apply-patch

test: ## Run tests with Bazel
	$(BAZEL_TEST) //...

test-cover: ## Run tests with coverage
	$(BAZEL_TEST) --collect_code_coverage //...
	@echo "Coverage reports generated in bazel-testlogs/"

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

format: ## Format code (Go fmt)
	@which gofmt > /dev/null || (echo "gofmt not found. Install Go." && exit 1)
	gofmt -s -w .
	@echo "Code formatted"

vet: ## Run go vet (if using Go directly)
	@which go > /dev/null || (echo "go not found" && exit 1)
	go vet ./...

fmt-check: ## Check if code is formatted
	@which gofmt > /dev/null || (echo "gofmt not found" && exit 1)
	@if [ "$$(gofmt -d .)" != "" ]; then \
		echo "Code is not formatted. Run 'make format'"; \
		exit 1; \
	fi

clean: ## Clean Bazel build artifacts
	$(BAZEL_CLEAN)
	rm -rf bazel-*

install: build ## Install binaries (copy from bazel-bin to PATH location)
	@echo "Binaries are in $(BIN_DIR)/"
	@echo "To install, copy manually or add $(BIN_DIR)/cmd/* to your PATH"

run-list: build-list ## Run list-deployments command
	$(BAZEL_RUN) //cmd/list-deployments:list-deployments -- -n default

run-symptom: build-symptom ## Run symptom-collection command (example)
	@echo "Example: $(BAZEL_RUN) //cmd/symptom-collection:symptom-collection -- -n default -p \"pod1 pod2\""

run-patch: build-patch ## Run apply-patch command (example)
	@echo "Example: $(BAZEL_RUN) //cmd/apply-patch:apply-patch -- -p /path/to/patch.so -s service-name"

verify: fmt-check test ## Run all verification checks

.DEFAULT_GOAL := help

