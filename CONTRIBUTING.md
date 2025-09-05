# Contributing to jscal

Thank you for your interest in contributing to the jscal library! We welcome contributions from the community.

## How to Contribute

### Reporting Issues

- Check if the issue has already been reported
- Include Go version and OS information
- Provide a minimal reproducible example
- Include actual vs expected behavior

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Run linting (`go vet ./...`)
7. Commit with descriptive messages
8. Push to your fork
9. Open a pull request

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add godoc comments for exported functions
- Keep functions focused and small
- Use meaningful variable names

### Testing

- Write tests for all new functionality
- Maintain or increase test coverage
- Include both positive and negative test cases
- Use table-driven tests where appropriate

### Documentation

- Update README.md if adding features
- Add godoc comments for all exported types and functions
- Include examples for complex functionality
- Keep documentation clear and concise

## Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/jscal.git
cd jscal

# Install dependencies
go mod download

# Run tests
go test ./...

# Build CLI
cd cmd/jscal && go build

# Run linting
go vet ./...
```

## Areas for Contribution

- Additional calendar format converters (Google Calendar, Microsoft Graph)
- Performance optimizations
- Additional validation rules
- Timezone handling improvements
- Documentation and examples
- Bug fixes

## Questions?

Feel free to open an issue for any questions about contributing.