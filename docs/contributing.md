# Contributing to WaqfWise Community Edition

Thank you for your interest in contributing to WaqfWise! This document provides guidelines for contributing to the Community Edition.

## Code of Conduct

Please be respectful and professional in all interactions with the community.

## What Can I Contribute?

The Community Edition (licensed under AGPL v3) welcomes contributions to:

- Core authentication features
- Campaign management
- Payment processing (basic features)
- Donor management
- Nazir management
- Basic reporting
- Bug fixes
- Documentation improvements
- Test coverage

**Note:** Enterprise features (multitenancy, advanced analytics, white-labeling, etc.) are not open for community contributions as they are under a commercial license.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/waqfwise.git`
3. Create a feature branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Write/update tests
6. Ensure tests pass: `go test -tags=community ./internal/core/...`
7. Commit your changes: `git commit -m "feat: add your feature"`
8. Push to your fork: `git push origin feature/your-feature-name`
9. Open a Pull Request

## Commit Message Guidelines

We follow conventional commits:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions or changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

## Code Style

- Follow standard Go conventions
- Run `gofmt` before committing
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small

## Testing

- Write unit tests for new features
- Ensure existing tests pass
- Aim for good test coverage
- Test edge cases

## Pull Request Process

1. Update documentation if needed
2. Add/update tests
3. Ensure CI passes
4. Request review from maintainers
5. Address review feedback
6. Squash commits if requested

## License

By contributing to WaqfWise Community Edition, you agree that your contributions will be licensed under the AGPL v3 license.

## Questions?

If you have questions about contributing, please open an issue or reach out to the maintainers.

Thank you for helping make WaqfWise better!
