# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

"Терем" is a Go utility designed for simplifying work with utilities on routers with entware/openwrt. It's structured as a command-line application with internationalization support (Russian, English, Belarusian, Kazakh) and custom logging capabilities. It is a TUI application that allows you to manage utilities on routers with entware/openwrt. In a TUI application, the user interface is displayed in the terminal, and the user can interact with it using keyboard shortcuts. User can navigate through the interface using arrow keys and perform actions using function keys ans setting up utilities which are available and popular on routers with entware/openwrt.
The TUI interface is implemented using the termos library (https://github.com/qzeleza/termos), which provides a simple and easy-to-use interface for creating TUI applications. It's woks with argument through Cobra CLI. Arguments are used to pass values to the application to run the concrete application or to run the application with the concrete configuration.

## Build Commands

The project uses a Makefile in the `app/` directory for common development tasks:

```bash
cd app

# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Install dependencies
make deps

# Clean build artifacts
make clean

# Install to GOPATH
make install
```

## Project Architecture

### Directory Structure
- `app/` - Main Go application directory
  - `main.go` - Application entry point
  - `cmd/args/` - Cobra CLI command definitions
  - `cmd/terem/` - Core application logic and configuration
  - `internal/config/` - Configuration management with environment variable support
  - `internal/lang/` - Internationalization module with lightweight translation system
- `docs/` - Mintlify documentation site
- `builder/` - Build output directory

### Key Components

**Application Configuration (`app/cmd/terem/start.go`)**
- `AppConfig` struct centralizes logging, configuration, context, and language settings
- Supports graceful shutdown with context cancellation
- Custom logger initialization with configurable levels and file rotation

**CLI Framework (`app/cmd/args/root.go`)**
- Uses Cobra for command-line interface
- Currently contains placeholder Hugo references that should be updated to Terem-specific content

**Configuration System (`app/internal/config/config.go`)**
- Environment-driven configuration with defaults
- Supports `ARCH`, `DEV_MODE`, and `DEBUG` environment variables
- Architecture-aware (defaults to ARM for router environments)

**Internationalization (`app/internal/lang/ternslate.go`)**
- Lightweight translation system without external dependencies
- Enum-based language selection (RUS, ENG, BEL, KAZ)
- Russian as primary language with dictionary-based translations
- Fallback to original Russian text when translations are missing

### Dependencies

The project uses minimal external dependencies:
- `github.com/spf13/cobra` - CLI framework
- `github.com/qzeleza/zlogger` - Custom logging library
- Go 1.25.0+ required

## Development Notes

- The main entry point is `app/main.go`, not the deleted `app/cmd/terem/main.go`
- Build output goes to `../builder/` directory relative to the app folder
- Logging files default to `/tmp/terem.log` and `/tmp/terem.sock`
- The CLI command definitions in `cmd/args/root.go` need updating from Hugo placeholders to Terem-specific content
- Documentation is managed through Mintlify and available at https://terem.zeleza.com

## Testing

Run tests from the `app/` directory:
```bash
make test           # Run all tests
make test-coverage  # Run tests with HTML coverage report
```