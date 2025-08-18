# Contributing to Gonzo

First off, thank you for considering contributing to Gonzo! It's people like you that make Gonzo such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps to reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed and expected**
* **Include logs and screenshots if possible**
* **Include your environment details** (OS, Go version, terminal emulator)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a detailed description of the suggested enhancement**
* **Provide specific examples to demonstrate the enhancement**
* **Describe the current behavior and expected behavior**
* **Explain why this enhancement would be useful**

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes (`make test`)
5. Make sure your code follows the existing style (`make fmt vet`)
6. Issue that pull request!

## Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/control-theory/gonzo.git
   cd gonzo
   ```

2. **Install Go 1.21 or higher**
   ```bash
   # Check your Go version
   go version
   ```

3. **Install dependencies**
   ```bash
   make deps
   ```

4. **Build the project**
   ```bash
   make build
   ```

5. **Run tests**
   ```bash
   make test
   ```

## Development Workflow

### Before committing:

1. **Format your code**
   ```bash
   make fmt
   ```

2. **Run the linter**
   ```bash
   make vet
   ```

3. **Run tests**
   ```bash
   make test
   ```

4. **Build to ensure it compiles**
   ```bash
   make build
   ```

Or run all at once:
```bash
make dev
```

### Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

### Testing

* Write tests for new functionality
* Ensure all tests pass before submitting PR
* Include both unit tests and integration tests where appropriate
* Test with different log formats (JSON, OTLP, plain text)

### Documentation

* Update the README.md if you change functionality
* Update the USAGE_GUIDE.md for new features
* Comment your code where necessary
* Update help text for new CLI flags

## Official Base colors:
LIGHT BLUE: #0F9EFC
BLACK: #000000
GREEN: #49E209
WHITE: #FFFFFF
DARK BLUE: #081C39
GRAY: #BCBEC0

## Project Structure

```
gonzo/
├── cmd/gonzo/                  # Main application entry point
│   ├── main.go                # CLI setup and initialization
│   ├── app.go                 # Application configuration and setup
│   ├── extractors.go          # Data extraction utilities
│   └── processing.go          # Log processing logic
├── internal/                   # Private application code
│   ├── tui/                   # Terminal UI components
│   │   ├── model.go           # Main Bubble Tea model
│   │   ├── view.go            # Rendering logic
│   │   ├── update.go          # Event handling
│   │   ├── components.go      # Reusable UI components
│   │   ├── charts.go          # Chart rendering
│   │   ├── tables.go          # Table components
│   │   ├── modals.go          # Modal dialogs
│   │   ├── navigation.go      # Navigation handling
│   │   ├── formatting.go      # Text formatting utilities
│   │   ├── severity.go        # Log severity handling
│   │   ├── styles.go          # UI styling definitions
│   │   ├── drain3_manager.go  # Drain3 integration
│   │   └── splash.txt         # Startup splash screen
│   ├── analyzer/              # Log analysis engine
│   │   ├── otlp.go           # OTLP log analysis
│   │   └── text.go           # Plain text analysis
│   ├── memory/                # Frequency tracking
│   │   └── frequency.go       # Frequency counting logic
│   ├── otlplog/              # OTLP format handling
│   │   ├── converter.go       # OTLP log conversion
│   │   └── detector.go        # OTLP format detection
│   ├── drain3/               # Drain3 log clustering
│   │   └── impl.go           # Drain3 implementation
│   ├── ai/                   # AI integration
│   │   └── openai.go         # OpenAI API integration
│   ├── output/               # Output handlers
│   │   └── stdout.go         # Standard output handler
│   └── reader/               # Input readers
│       └── stdin.go          # Standard input reader
├── docs/                      # Documentation assets
│   └── screenshot.png         # Project screenshot
├── examples/                  # Configuration examples
│   └── config.yml            # Example configuration file
├── build/                     # Build artifacts directory
├── Makefile                  # Build automation
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── test scripts              # Various testing scripts
```

## Style Guide

### Go Code Style

* Follow standard Go conventions
* Use `gofmt` for formatting
* Use meaningful variable names
* Keep functions small and focused
* Write clear comments for exported functions
* Handle errors explicitly

### TUI Guidelines

* Maintain consistent keyboard shortcuts
* Use color sparingly and meaningfully
* Ensure the UI is responsive
* Test on different terminal sizes
* Support both mouse and keyboard navigation

## Release Process

1. Update version numbers
2. Update CHANGELOG.md
3. Create a git tag
4. Push tag to trigger release build
5. GitHub Actions will create the release

## Questions?

Feel free to open an issue with the question label or reach out to the maintainers directly.

Thank you for contributing! 🎉