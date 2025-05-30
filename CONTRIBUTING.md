# Contributing to UploadThing CLI

Thank you for your interest in contributing to the UploadThing CLI! We welcome contributions from the community and are grateful for your help in making this tool better.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- UploadThing account for testing

### Development Setup

1. **Fork the repository**
   ```bash
   # Click the "Fork" button on GitHub, then clone your fork
   git clone https://github.com/yourusername/uploadthing-cli.git
   cd uploadthing-cli
   ```

2. **Set up your development environment**
   ```bash
   # Install dependencies
   go mod download
   
   # Build the project
   go build -o ut .
   
   # Run tests
   go test ./...
   ```

3. **Configure for testing**
   ```bash
   # Set up your test configuration
   ./ut config set-secret YOUR_TEST_SECRET_KEY
   ```

## 📝 How to Contribute

### Reporting Issues

- **Search existing issues** first to avoid duplicates
- **Use the issue templates** when available
- **Provide clear reproduction steps** for bugs
- **Include system information** (OS, Go version, etc.)

### Suggesting Features

- **Check existing feature requests** first
- **Explain the use case** and why it would be valuable
- **Provide examples** of how the feature would work
- **Consider backward compatibility**

### Code Contributions

1. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

2. **Make your changes**
   - Follow the coding standards (see below)
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   # Run tests
   go test ./...
   
   # Test the CLI manually
   go build -o ut .
   ./ut --help
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add amazing feature"
   ```

5. **Push and create a Pull Request**
   ```bash
   git push origin feature/amazing-feature
   ```

## 🎯 Coding Standards

### Go Code Style

- Follow standard Go formatting: `go fmt ./...`
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and small
- Handle errors appropriately with proper error wrapping

### Code Organization

```
cmd/           # CLI commands
├── root.go    # Root command setup
├── push.go    # Upload functionality
├── fetch.go   # Download functionality
├── list.go    # List files functionality
├── config.go  # Configuration management
└── utils.go   # Shared utilities

config/        # Configuration package
main.go        # Application entry point
```

### Error Handling

- Use `fmt.Errorf` with `%w` verb for error wrapping
- Provide meaningful error messages
- Handle edge cases gracefully
- Clean up resources (files, connections) on errors

### Testing

- Write unit tests for new functions
- Test error conditions
- Use table-driven tests when appropriate
- Mock external dependencies (HTTP calls, file system)

## 📚 Documentation

### Code Documentation

- Document all exported functions and types
- Use clear, concise comments
- Include examples in documentation when helpful

### README Updates

- Update feature lists when adding new functionality
- Add new command examples
- Update installation instructions if needed

## 🔄 Pull Request Process

1. **Ensure your PR has a clear title and description**
2. **Reference any related issues** using `Fixes #123` or `Closes #123`
3. **Include tests** for new functionality
4. **Update documentation** as needed
5. **Ensure all checks pass** (tests, linting, etc.)
6. **Respond to review feedback** promptly

### PR Title Format

Use conventional commit format:
- `feat: add new feature`
- `fix: resolve bug in upload`
- `docs: update README`
- `refactor: improve error handling`
- `test: add unit tests for config`

## 🧪 Testing Guidelines

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./cmd
```

### Writing Tests

```go
func TestUploadFile(t *testing.T) {
    tests := []struct {
        name     string
        filePath string
        wantErr  bool
    }{
        {
            name:     "valid file",
            filePath: "testdata/sample.txt",
            wantErr:  false,
        },
        {
            name:     "non-existent file",
            filePath: "testdata/missing.txt",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := uploadFile(tt.filePath)
            if (err != nil) != tt.wantErr {
                t.Errorf("uploadFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 🏷️ Release Process

1. **Version bumping** follows semantic versioning (SemVer)
2. **Changelog** is updated with new features and fixes
3. **Tags** are created for releases
4. **Binaries** are built and attached to releases

## 💬 Community

- **Discussions**: Use GitHub Discussions for questions and ideas
- **Issues**: Use GitHub Issues for bugs and feature requests
- **Code of Conduct**: Be respectful and inclusive

## 📞 Getting Help

- **Documentation**: Check the README and Wiki first
- **Issues**: Search existing issues for similar problems
- **Discussions**: Ask questions in GitHub Discussions

## 🙏 Recognition

Contributors will be recognized in:
- The project README
- Release notes
- GitHub contributors page

Thank you for contributing to UploadThing CLI! 🎉 