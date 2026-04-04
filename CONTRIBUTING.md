# Contributing to sb

Thanks for your interest in contributing to `sb`! Here's how you can help.

## Getting Started

1. Fork and clone the repository
2. Make sure you have Go installed (see `go.mod` for the required version)
3. Build and run tests:
   ```sh
   go build ./...
   go test ./...
   ```

## Making Changes

1. Create a branch for your change
2. Write your code and tests
3. Make sure all tests pass: `go test ./...`
4. Make sure the code builds cleanly: `go build ./...`
5. Commit with a clear message explaining *why* the change was made

## Pull Requests

- Keep PRs focused on a single change
- Include tests for new functionality
- Update documentation if your change affects user-facing behavior
- Make sure CI passes before requesting review

## Reporting Bugs

Open an issue with:
- What you expected to happen
- What actually happened
- Steps to reproduce
- macOS version and `sb` version (or commit)

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use `github.com/stretchr/testify/assert` for test assertions
- Keep test output readable — prefer test helpers over deeply nested table-driven tests

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
