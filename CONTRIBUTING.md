# Contributing to postalclient-go

Thank you for considering contributing to postalclient-go! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct: be respectful, considerate, and collaborative.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue tracker to avoid duplicates. When you create a bug report, include as many details as possible:

- Use a clear and descriptive title
- Describe the exact steps to reproduce the problem
- Provide specific examples to demonstrate the steps
- Describe the behavior you observed and what you expected to see
- Include Go version, package version, and OS information
- If possible, include a minimal code example that reproduces the issue

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful to most users
- If possible, include code examples showing how the feature would be used

### Pull Requests

- Fill in the required template
- Follow the Go style guide
- Include tests for new features or bug fixes
- Update documentation as needed
- End all files with a newline
- Make sure your code passes all tests and linting

## Development Workflow

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes
4. Run tests and linting locally
5. Commit your changes with a descriptive commit message
6. Push your branch to your fork
7. Submit a pull request

### Setting Up the Development Environment

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/postalclient-go.git
cd postalclient-go

# Add the original repository as a remote
git remote add upstream https://github.com/Suhaibinator/postalclient-go.git

# Get dependencies
go mod download
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linting
golangci-lint run
```

## Style Guide

This project follows the standard Go style guide and best practices:

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `gofmt` on your code before committing
- Use meaningful variable and function names
- Write comprehensive comments and documentation
- Keep functions small and focused on a single responsibility
- Write tests for new functionality

## Documentation

- Update the README.md if your changes affect the usage or installation instructions
- Add or update GoDoc comments for public functions, types, and methods
- Include code examples in documentation when helpful

## Releasing

The project maintainers will handle releases. If you believe a new release is needed, please open an issue.

## Questions?

If you have any questions or need help, please open an issue or reach out to the maintainers.

Thank you for contributing to postalclient-go!
