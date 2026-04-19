# Contributing to Sonic-Screwdriver

Thank you for your interest in contributing to Sonic-Screwdriver! This document outlines the process for contributing to the project.

## 🎯 Ways to Contribute

There are many ways to contribute:
- **Bug Reports**: File issues for bugs you encounter
- **Feature Requests**: Suggest new features or improvements
- **Code Contributions**: Submit pull requests with fixes or features
- **Documentation**: Improve existing documentation or add new guides
- **Testing**: Help test new features and releases

## 📋 Getting Started

### Prerequisites

- Go 1.21+
- Docker (for container runtime)
- Git
- Make

### Setup

```bash
# Clone the repository
git clone https://github.com/sonic-family/sonic-screwdriver.git
cd sonic-screwdriver

# Build the project
make build

# Run tests
make test
```

## 🚀 Development Workflow

### 1. Fork the Repository

Fork the repository to your GitHub account and clone it locally.

### 2. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 3. Make Changes

- Follow the existing code style
- Add tests for new features
- Update documentation as needed

### 4. Commit Changes

```bash
git commit -m "feat: add your feature"
```

### 5. Push Changes

```bash
git push origin feature/your-feature-name
```

### 6. Open a Pull Request

Open a PR against the `main` branch with a clear description of your changes.

## 📝 Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Add comments for complex logic
- Write tests for new functionality

## 🧪 Testing

All contributions should include tests:

```bash
# Run unit tests
make test

# Run integration tests (requires Docker)
make test-integration
```

## 📚 Documentation

Update documentation when making changes:
- Add comments to code
- Update README files
- Add examples where helpful

## 🎉 Code Review

All contributions will be reviewed by the maintainers. We aim to respond within 72 hours.

## 🤝 Community

Join our community:
- **Discussions**: GitHub Discussions
- **Issues**: GitHub Issues
- **Pull Requests**: GitHub Pull Requests

## 📄 License

By contributing to Sonic-Screwdriver, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Sonic-Screwdriver! 🎊