# Contributing

Thank you for your interest in contributing to the Service Registry project.
Please follow these guidelines to ensure consistency and quality.

## Development Setup

1. Ensure you have Go 1.22+ installed
2. Install golangci-lint: https://golangci-lint.run/usage/install/
3. Clone the repo and run `make test` to verify your setup

## Code Style

- Follow standard Go conventions and idioms
- Run `make fmt` before committing to ensure consistent formatting
- Run `make lint` to catch common issues
- Keep functions focused and reasonably sized

## Testing

- All new features should include tests
- All bug fixes should include a test that reproduces the issue
- Run `make test` before submitting changes
- Aim to maintain or improve test coverage

## Commit Messages

Write clear, descriptive commit messages:

- Use the imperative mood ("Add feature" not "Added feature")
- Keep the first line under 72 characters
- Reference issue numbers when applicable

**Good examples:**
- `Add bulk service registration endpoint`
- `Fix health check timeout handling`
- `Update README with new API endpoints`

**Bad examples:**
- `Fixed stuff`
- `WIP`
- `Updates`

## Pull Requests

1. Create a feature branch from `main`
2. Make your changes with clear, atomic commits
3. Ensure all tests pass (`make test`)
4. Ensure linting passes (`make lint`)
5. Update documentation if needed
6. Submit a pull request with a clear description

## Project Structure

```
.
├── cmd/server/         # Application entrypoint
├── internal/
│   ├── api/            # HTTP handlers
│   ├── models/         # Data models
│   └── store/          # Data storage
├── Makefile            # Build and development commands
└── README.md           # Project documentation
```

## Questions?

If you have questions about contributing, please open an issue for discussion.
