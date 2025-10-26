# Contributing to gosh

Thank you for your interest in contributing to gosh! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/gosh.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -am 'Add some feature'`
7. Push to the branch: `git push origin feature/your-feature-name`
8. Create a Pull Request

## Development Setup

### Prerequisites

- Go 1.19 or higher
- Git

### Building

```bash
go build -o gosh .
```

### Testing

Run the shell to test your changes:

```bash
./gosh
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Write clear, descriptive commit messages
- Add comments for complex logic

## Project Structure

```
gosh/
├── main.go              # Entry point
├── pkg/
│   ├── shell/          # Shell REPL implementation
│   └── workspace/      # Session and file management
└── examples/           # Example scripts
```

## Adding Features

When adding new features:

1. Ensure cross-platform compatibility (Linux, Windows, macOS)
2. Update documentation in README.md
3. Add examples if applicable
4. Consider backward compatibility

## Reporting Bugs

When reporting bugs, please include:

- Go version
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages (if any)

## Feature Requests

We welcome feature requests! Please:

- Check existing issues first
- Clearly describe the feature
- Explain the use case
- Consider implementation complexity

## Code Review Process

All submissions require review. We use GitHub pull requests for this purpose.

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

## Questions?

Feel free to open an issue for any questions or concerns.
