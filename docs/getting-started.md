# Getting Started

[English](getting-started.md) | [日本語](getting-started_ja.md)

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap taro33333/tap
brew install --cask goreload
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

## Quick Start

### 1. Initialize Configuration

Navigate to your Go project directory and run:

```bash
goreload init
```

This creates a `goreload.yaml` configuration file with default settings.

### 2. Start Watching

```bash
goreload
```

goreload will:

1. Build your application
2. Start the compiled binary
3. Watch for file changes
4. Automatically rebuild and restart when changes are detected

### 3. Stop

Press `Ctrl+C` to stop goreload.

## Example Project Setup

```
myapp/
├── main.go
├── go.mod
├── go.sum
└── goreload.yaml
```

### Minimal `goreload.yaml`

```yaml
root: "."
tmp_dir: "tmp"

build:
  cmd: "go build -o ./tmp/main ."
  bin: "./tmp/main"

watch:
  extensions:
    - ".go"
  dirs:
    - "."
  exclude_dirs:
    - "tmp"
    - "vendor"

log:
  level: "info"
```

## Next Steps

- [Configuration Reference](./configuration.md) - Learn about all configuration options
- [CLI Reference](./cli.md) - Explore command-line options
