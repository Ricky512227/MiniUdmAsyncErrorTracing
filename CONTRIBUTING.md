# Contributing to MiniUdm Async Error Tracing

Thank you for your interest in contributing to MiniUdm Async Error Tracing! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and considerate of others when contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/MiniUdmAsyncErrorTracing.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes: `make test`
6. Format your code: `make format`
7. Commit your changes: `git commit -m 'Add your feature'`
8. Push to your fork: `git push origin feature/your-feature-name`
9. Open a Pull Request

## Development Workflow

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/Ricky512227/MiniUdmAsyncErrorTracing.git
cd MiniUdmAsyncErrorTracing

# Install dependencies
make deps

# Build the project
make build

# Run tests
make test
```

### Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (run `make format`)
- Add comments for exported functions, types, and packages
- Keep functions small and focused
- Use meaningful variable and function names

### Testing

- Write tests for new functionality
- Ensure existing tests pass: `make test`
- Aim for high test coverage
- Use table-driven tests where appropriate

### Commits

- Write clear, descriptive commit messages
- Use present tense ("Add feature" not "Added feature")
- Reference issues in commit messages when applicable: "Fix #123"

## Pull Request Process

1. Ensure your code follows the project's code style
2. Update documentation if needed
3. Add tests for new features
4. Ensure all tests pass
5. Update CHANGELOG.md if applicable
6. Request review from maintainers

### Pull Request Checklist

- [ ] Code follows the project's style guidelines
- [ ] Tests have been added/updated
- [ ] All tests pass
- [ ] Documentation has been updated
- [ ] Commit messages are clear and descriptive

## Project Structure

- `cmd/` - Command-line applications
- `pkg/` - Reusable packages (can be imported by other projects)
- `internal/` - Internal packages (not for external use)
- `configs/` - Configuration files
- `examples/` - Example code
- `docs/` - Documentation

## Reporting Issues

When reporting issues, please include:

- Description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version, Kubernetes version)
- Relevant logs or error messages

## Feature Requests

For feature requests, please:

- Describe the feature and its use case
- Explain why it would be useful
- Consider potential implementation approaches
- Check if a similar feature already exists

## Questions?

Feel free to open an issue for questions or reach out to the maintainers.

Thank you for contributing!

