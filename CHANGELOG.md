# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Production-grade project structure with cmd/, pkg/, and internal/ directories
- Command-line tools: list-deployments, symptom-collection, apply-patch
- Comprehensive configuration management with Viper
- Structured logging with Zap
- Kubernetes client wrapper with high-level operations
- Patch management with MD5 validation and health monitoring
- Symptom collection with parallel log monitoring
- Unit tests for utility functions
- Makefile with common development tasks
- CI/CD pipeline with GitHub Actions
- Comprehensive documentation (README, CONTRIBUTING, ARCHITECTURE)
- Example code demonstrating usage
- Code quality tools (golangci-lint configuration)

### Changed
- Reorganized codebase from flat structure to modular packages
- Improved error handling with proper error wrapping
- Enhanced logging throughout the application
- Better separation of concerns with interface-based design

## [0.1.0] - Initial Release

### Added
- Basic deployment listing functionality
- Symptom collection workflow (initial implementation)
- Patch application workflow (initial implementation)
- Common utility functions

