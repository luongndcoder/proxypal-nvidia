# Contributing to ProxyPal NVIDIA Load Balancer

Thank you for your interest in contributing to ProxyPal! We welcome contributions from the community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/proxypal-nvidia.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes
5. Test your changes: `make test`
6. Commit your changes: `git commit -m 'Add some amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized testing)
- Make (optional, for build automation)

### Install Dependencies

```bash
go mod download
```

### Build

```bash
make build
# or
go build -o proxypal ./cmd/proxypal
```

### Run Tests

```bash
make test
# or
go test -v ./...
```

### Run Locally

```bash
# Copy and configure
cp config.example.yaml config.yaml
# Edit config.yaml with your API keys

# Run
make run
# or
./proxypal
```

## Code Style

- Follow Go best practices and conventions
- Run `go fmt ./...` before committing
- Run `go vet ./...` to check for common issues
- Write tests for new functionality
- Keep functions small and focused
- Use meaningful variable and function names

## Testing

- Write unit tests for new features
- Ensure all tests pass before submitting PR
- Aim for high test coverage
- Test both success and error cases

Example:
```bash
go test ./... -v -race -coverprofile=coverage.txt
```

## Pull Request Guidelines

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] All tests pass
- [ ] New features have tests
- [ ] Documentation is updated (README, comments, etc.)
- [ ] Commit messages are clear and descriptive
- [ ] No unnecessary dependencies added

### PR Description

Please include:
- What changes were made
- Why the changes were necessary
- Any breaking changes
- Screenshots (if applicable)
- Related issues

### Example PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Code formatted
- [ ] Documentation updated
```

## Reporting Bugs

### Before Reporting

- Check if the bug has already been reported
- Try the latest version
- Collect relevant information (OS, Go version, config, logs)

### Bug Report Template

```markdown
**Describe the bug**
A clear description of the bug

**To Reproduce**
Steps to reproduce:
1. Configure with '...'
2. Run command '...'
3. See error

**Expected behavior**
What you expected to happen

**Environment:**
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.21]
- ProxyPal version: [e.g., v1.0.0]

**Logs**
```
Paste relevant logs here
```

**Additional context**
Any other context about the problem
```

## Suggesting Features

We welcome feature suggestions! Please:

1. Check if the feature has already been requested
2. Clearly describe the feature and use case
3. Explain why it would be useful
4. Provide examples if possible

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of the problem

**Describe the solution you'd like**
A clear description of what you want to happen

**Describe alternatives you've considered**
Alternative solutions or features you've considered

**Additional context**
Any other context, screenshots, or examples
```

## Code Review Process

1. Maintainers will review your PR
2. Feedback will be provided if changes are needed
3. Once approved, PR will be merged
4. Changes will be included in next release

## Questions?

- Open an issue for questions
- Check existing issues and documentation first
- Be respectful and patient

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make ProxyPal better for everyone. We appreciate your time and effort!
