# goreload

[![CI](https://github.com/taro33333/goreload/actions/workflows/ci.yaml/badge.svg)](https://github.com/taro33333/goreload/actions/workflows/ci.yaml)
[![Release](https://github.com/taro33333/goreload/actions/workflows/release.yaml/badge.svg)](https://github.com/taro33333/goreload/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/taro33333/goreload)](https://goreportcard.com/report/github.com/taro33333/goreload)

A hot reload tool for Go applications. Watches your source files and automatically rebuilds and restarts your application when changes are detected.

## Features

- File watching with configurable extensions
- Debounced builds to avoid rapid rebuilds
- Graceful process shutdown with configurable timeout
- Colored log output
- Glob pattern support for file exclusion
- Recursive directory watching
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Homebrew (macOS/Linux)

```bash
brew install taro33333/tap/goreload
```

### Go Install

```bash
go install github.com/taro33333/goreload/cmd/goreload@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/taro33333/goreload/releases).

### Build from Source

```bash
git clone https://github.com/taro33333/goreload.git
cd goreload
go build -o goreload ./cmd/goreload
```

## Usage

### Quick Start

```bash
# Initialize configuration file
goreload init

# Run with default configuration
goreload
```

### Specify Configuration File

```bash
goreload -c ./config.yaml
```

### Show Version

```bash
goreload version
```

## Configuration

Create a `goreload.yaml` file in your project root:

```yaml
# Project root directory
root: "."

# Temporary directory for build artifacts
tmp_dir: "tmp"

# Build settings
build:
  # Build command
  cmd: "go build -o ./tmp/main ."
  # Binary to execute
  bin: "./tmp/main"
  # Arguments to pass to the binary
  args: []
  # Delay before building after file change (debounce)
  delay: "200ms"
  # Grace period for process termination
  kill_delay: "500ms"

# File watching settings
watch:
  # File extensions to watch
  extensions:
    - ".go"
  # Directories to watch
  dirs:
    - "."
  # Directories to exclude
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  # Files to exclude (glob patterns)
  exclude_files:
    - "*_test.go"

# Logging settings
log:
  # Enable colored output
  color: true
  # Show timestamps
  time: true
  # Log level: debug, info, warn, error
  level: "info"
```

## Output Example

```
   __ _  ___  _ __ ___| | ___   __ _  __| |
  / _` |/ _ \| '__/ _ \ |/ _ \ / _` |/ _` |
 | (_| | (_) | | |  __/ | (_) | (_| | (_| |
  \__, |\___/|_|  \___|_|\___/ \__,_|\__,_|
  |___/                            v0.1.0

[INFO] watching: .
[INFO] excluding: [tmp vendor .git node_modules]
[INFO] building...
[INFO] ✓ build completed (1.23s)
[INFO] ✓ running ./tmp/main

[INFO] main.go changed
[INFO] building...
[INFO] ✓ build completed (0.45s)
[INFO] ✓ running ./tmp/main
```

## Architecture

```
goreload/
├── cmd/goreload/main.go     # CLI entry point
├── internal/
│   ├── config/              # Configuration loading and validation
│   ├── watcher/             # File system watching
│   ├── builder/             # Build command execution
│   ├── runner/              # Process management
│   ├── engine/              # Orchestrator
│   └── logger/              # Structured logging
├── goreload.yaml            # Sample configuration
└── README.md
```

## Dependencies

- [fsnotify](https://github.com/fsnotify/fsnotify) - File system notifications
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [color](https://github.com/fatih/color) - Colored terminal output

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT
