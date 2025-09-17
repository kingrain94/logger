# Contributing to Go Logger

Thank you for your interest in contributing to the Go Logger project! We welcome contributions from the community and are pleased to have you join us.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/logger.git
   cd logger
   ```
3. **Add the upstream repository** as a remote:
   ```bash
   git remote add upstream https://github.com/kingrain94/logger.git
   ```

## How to Contribute

### Types of Contributions

We welcome several types of contributions:

- **Bug fixes**: Fix existing bugs or issues
- **Feature additions**: Add new functionality
- **Documentation**: Improve or add documentation
- **Tests**: Add or improve test coverage
- **Performance improvements**: Optimize existing code
- **Code quality**: Refactor code for better maintainability

### Before You Start

1. **Check existing issues** to see if your contribution is already being worked on
2. **Create an issue** to discuss major changes before implementing them
3. **Search existing pull requests** to avoid duplicating work

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Local Development

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Run tests** to ensure everything works:
   ```bash
   go test ./...
   ```

3. **Run tests with coverage**:
   ```bash
   go test -race -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

4. **Run benchmarks**:
   ```bash
   go test -bench=. -benchmem
   ```

5. **Check code formatting**:
   ```bash
   gofmt -s -w .
   ```

6. **Run linting** (if you have golangci-lint installed):
   ```bash
   golangci-lint run
   ```

## Pull Request Process

### Before Submitting

1. **Create a feature branch** from `master`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Add or update tests** for your changes

4. **Update documentation** if necessary

5. **Ensure all tests pass**:
   ```bash
   go test ./...
   ```

6. **Commit your changes** with a clear commit message:
   ```bash
   git commit -m "Add feature: description of your feature"
   ```

### Submitting the Pull Request

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a pull request** on GitHub with:
   - Clear title and description
   - Reference to related issues (if any)
   - List of changes made
   - Screenshots (if applicable)

3. **Respond to review feedback** promptly

### Pull Request Guidelines

- **One feature per PR**: Keep pull requests focused on a single feature or bug fix
- **Small, incremental changes**: Large PRs are harder to review
- **Clear commit messages**: Use descriptive commit messages
- **Update documentation**: Include documentation updates for new features
- **Add tests**: Ensure new code is properly tested

## Coding Standards

### Go Style Guidelines

We follow standard Go conventions:

- **Use `gofmt`** for code formatting
- **Follow effective Go guidelines**
- **Use meaningful variable and function names**
- **Write clear comments** for exported functions
- **Keep functions small and focused**

### Code Organization

```go
// Package comment should describe the package purpose
package logger

import (
    // Standard library imports first
    "fmt"
    "os"
    
    // Third-party imports
    "go.uber.org/zap"
)

// Constants should be grouped
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// Variables should be grouped
var (
    logger *zap.Logger
    config Config
)

// Public functions should have comments
// Initialize initializes the logger with the given configuration
func Initialize(config Config) error {
    // Implementation
}
```

### Error Handling

- **Always handle errors** explicitly
- **Use meaningful error messages**
- **Wrap errors** with context when appropriate:
  ```go
  if err != nil {
      return fmt.Errorf("failed to initialize logger: %w", err)
  }
  ```

### Documentation

- **All exported functions** must have comments
- **Use godoc format** for documentation
- **Include examples** for complex functions
- **Keep comments up to date** with code changes

## Testing

### Test Requirements

- **All new code** must include tests
- **Maintain or improve** test coverage
- **Include edge cases** in tests
- **Use table-driven tests** where appropriate

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Benchmark Tests

Include benchmark tests for performance-critical code:

```go
func BenchmarkFunctionName(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FunctionName(testInput)
    }
}
```

## Documentation

### README Updates

When adding new features:

1. **Update the README.md** with new examples
2. **Add to the API reference** section
3. **Update the feature list** if applicable

### Code Comments

- **Exported functions** must have comments
- **Complex logic** should be explained
- **Use examples** in comments when helpful

### Changelog

For significant changes, consider updating a changelog or mentioning the changes in your pull request description.

## Issue Reporting

### Before Creating an Issue

1. **Search existing issues** to avoid duplicates
2. **Check the documentation** for answers
3. **Try the latest version** to see if the issue is already fixed

### Creating a Good Issue

Include the following information:

- **Clear title** describing the issue
- **Steps to reproduce** the problem
- **Expected behavior**
- **Actual behavior**
- **Go version** and OS
- **Relevant code snippets** or logs
- **Minimal example** that demonstrates the issue

### Issue Templates

#### Bug Report
```
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Initialize logger with...
2. Call function...
3. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g., Linux, macOS, Windows]
- Go version: [e.g., 1.21.0]
- Logger version: [e.g., v1.0.0]

**Additional context**
Any other context about the problem.
```

#### Feature Request
```
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or features you've considered.

**Additional context**
Any other context about the feature request.
```

## Getting Help

If you need help with contributing:

1. **Check the documentation** first
2. **Look at existing code** for examples
3. **Create an issue** with the "question" label
4. **Join discussions** in existing issues

## Recognition

Contributors will be recognized in:

- **GitHub contributors list**
- **Release notes** for significant contributions
- **Special thanks** in documentation

Thank you for contributing to Go Logger! 🎉
