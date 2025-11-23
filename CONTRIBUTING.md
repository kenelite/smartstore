# Contributing to SmartStore

Thank you for your interest in contributing to SmartStore! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (for local development)
- Git
- Make

### Initial Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/smartstore.git
   cd smartstore
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/kenelite/smartstore.git
   ```

4. Run the setup script:
   ```bash
   make setup
   ```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications

### 2. Make Your Changes

- Write clean, readable code
- Follow Go best practices
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests and Checks

```bash
# Format code
make fmt

# Run linter
make lint

# Run all tests
make test

# Or use the comprehensive test script
make test-all
```

### 4. Commit Your Changes

Use clear and descriptive commit messages:

```bash
git commit -m "feat: add support for Azure Blob Storage"
git commit -m "fix: resolve race condition in cache layer"
git commit -m "docs: update API documentation"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Build process or auxiliary tool changes

### 5. Keep Your Branch Updated

```bash
git fetch upstream
git rebase upstream/main
```

### 6. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines and these additional conventions:

1. **Naming**
   - Use clear, descriptive names
   - Follow Go naming conventions (camelCase for private, PascalCase for public)
   - Avoid abbreviations unless widely understood

2. **Error Handling**
   - Always check and handle errors
   - Wrap errors with context using `fmt.Errorf` with `%w`
   - Don't panic in library code

3. **Comments**
   - Write godoc comments for all exported functions and types
   - Keep comments concise and informative
   - Update comments when code changes

4. **Code Organization**
   - Keep functions small and focused
   - Group related functionality
   - Use interfaces for abstractions

### Example

```go
// ObjectStore defines the interface for object storage operations.
// All implementations must be safe for concurrent use.
type ObjectStore interface {
    // Put stores an object and returns its metadata.
    // Returns an error if the operation fails.
    Put(ctx context.Context, key string, data io.Reader) (*Metadata, error)
    
    // Get retrieves an object by key.
    // Returns ErrNotFound if the object doesn't exist.
    Get(ctx context.Context, key string) (io.ReadCloser, error)
}
```

## Testing

### Writing Tests

1. **Unit Tests**
   - Test individual functions and methods
   - Use table-driven tests when appropriate
   - Mock external dependencies

2. **Integration Tests**
   - Test component interactions
   - Use docker-compose for test dependencies
   - Clean up resources after tests

3. **Test Coverage**
   - Aim for at least 80% coverage for new code
   - Focus on critical paths and edge cases

### Example Test

```go
func TestObjectStore_Put(t *testing.T) {
    tests := []struct {
        name    string
        key     string
        data    string
        wantErr bool
    }{
        {
            name:    "valid object",
            key:     "test/file.txt",
            data:    "test content",
            wantErr: false,
        },
        {
            name:    "empty key",
            key:     "",
            data:    "test content",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            store := setupTestStore(t)
            defer cleanupTestStore(t, store)
            
            reader := strings.NewReader(tt.data)
            _, err := store.Put(context.Background(), tt.key, reader)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make test-bench

# Run short tests only
make test-short
```

## Pull Request Process

### Before Submitting

1. Ensure all tests pass: `make test-all`
2. Run linter: `make lint`
3. Update documentation if needed
4. Add or update tests for your changes
5. Rebase on the latest upstream/main

### Submitting

1. Push your changes to your fork
2. Create a Pull Request to the main repository
3. Fill out the PR template completely
4. Link any related issues

### PR Title Format

Use the same format as commit messages:

```
feat: add Azure Blob Storage adapter
fix: resolve memory leak in cache layer
docs: improve API documentation
```

### PR Description Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe the tests you ran and how to reproduce

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests added/updated
- [ ] All tests pass
```

### Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR

## Issue Reporting

### Bug Reports

Use the bug report template and include:

1. **Description**: Clear description of the bug
2. **Steps to Reproduce**: Detailed steps
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happens
5. **Environment**: OS, Go version, etc.
6. **Logs**: Relevant logs or error messages

### Feature Requests

Use the feature request template and include:

1. **Description**: Clear description of the feature
2. **Use Case**: Why is this feature needed
3. **Proposed Solution**: How you think it should work
4. **Alternatives**: Alternative solutions considered

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Project README](README.md)

## Questions?

If you have questions, feel free to:
- Open an issue for discussion
- Reach out to maintainers
- Check existing issues and PRs

Thank you for contributing to SmartStore! 🚀

