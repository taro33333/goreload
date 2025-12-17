# Configuration Reference

[English](configuration.md) | [日本語](configuration_ja.md)

goreload uses a YAML configuration file (default: `goreload.yaml`) to control its behavior.

## Configuration File Location

goreload looks for configuration in the following order:

1. Path specified by `-c` or `--config` flag
2. `goreload.yaml` in current directory

## Complete Configuration Example

```yaml
# Project root directory
root: "."

# Temporary directory for build artifacts
tmp_dir: "tmp"

# Build settings
build:
  cmd: "go build -o ./tmp/main ."
  bin: "./tmp/main"
  args: []
  delay: "200ms"
  kill_delay: "500ms"

# File watching settings
watch:
  extensions:
    - ".go"
  dirs:
    - "."
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  exclude_files:
    - "*_test.go"

# Logging settings
log:
  color: true
  time: true
  level: "info"
```

## Configuration Options

### Root Level

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `root` | string | `"."` | Project root directory. All relative paths are resolved from here. |
| `tmp_dir` | string | `"tmp"` | Directory for build artifacts. Created automatically if not exists. |

### Build Settings (`build`)

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `cmd` | string | `"go build -o ./tmp/main ."` | Build command to execute. Supports shell-style quoting. |
| `bin` | string | `"./tmp/main"` | Path to the compiled binary to execute. |
| `args` | []string | `[]` | Arguments to pass to the binary when running. |
| `delay` | duration | `"200ms"` | Debounce delay before triggering build after file change. |
| `kill_delay` | duration | `"500ms"` | Grace period for process termination before SIGKILL. |

#### Duration Format

Duration values support Go's duration format:

- `"100ms"` - 100 milliseconds
- `"1s"` - 1 second
- `"1m30s"` - 1 minute 30 seconds

#### Build Command Examples

```yaml
# Standard Go build
build:
  cmd: "go build -o ./tmp/main ."

# With build tags
build:
  cmd: "go build -tags=dev -o ./tmp/main ."

# With ldflags
build:
  cmd: "go build -ldflags='-s -w' -o ./tmp/main ."

# Using make
build:
  cmd: "make build"
  bin: "./bin/app"
```

### Watch Settings (`watch`)

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `extensions` | []string | `[".go"]` | File extensions to watch. At least one required. |
| `dirs` | []string | `["."]` | Directories to watch recursively. At least one required. |
| `exclude_dirs` | []string | `["tmp", "vendor", ".git", "node_modules"]` | Directories to exclude from watching. |
| `exclude_files` | []string | `[]` | File patterns to exclude (glob patterns supported). |

#### Extension Format

Extensions should include the leading dot:

```yaml
watch:
  extensions:
    - ".go"
    - ".html"
    - ".tmpl"
```

#### Glob Patterns

The `exclude_files` option supports glob patterns:

| Pattern | Description |
|---------|-------------|
| `*_test.go` | All test files |
| `*.gen.go` | All generated files |
| `mock_*.go` | All mock files |
| `*.pb.go` | All protobuf generated files |

### Log Settings (`log`)

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `color` | bool | `true` | Enable colored output. |
| `time` | bool | `true` | Show timestamps in log output. |
| `level` | string | `"info"` | Log level: `debug`, `info`, `warn`, `error`. |

#### Log Levels

| Level | Description |
|-------|-------------|
| `debug` | Verbose output for debugging |
| `info` | Normal operation messages |
| `warn` | Warnings that don't stop operation |
| `error` | Errors only |

## Environment Variables

goreload respects standard Go environment variables:

| Variable | Description |
|----------|-------------|
| `GO111MODULE` | Go modules mode |
| `GOOS` | Target operating system |
| `GOARCH` | Target architecture |
| `CGO_ENABLED` | Enable/disable CGO |

## Validation Rules

The configuration is validated on load. The following rules apply:

1. `build.cmd` - Must not be empty
2. `build.bin` - Must not be empty
3. `build.delay` - Must be non-negative
4. `build.kill_delay` - Must be non-negative
5. `watch.extensions` - Must have at least one extension
6. `watch.dirs` - Must have at least one directory
7. `log.level` - Must be one of: `debug`, `info`, `warn`, `error`

## Default Configuration

If no configuration file exists, goreload uses these defaults:

```yaml
root: "."
tmp_dir: "tmp"
build:
  cmd: "go build -o ./tmp/main ."
  bin: "./tmp/main"
  args: []
  delay: "200ms"
  kill_delay: "500ms"
watch:
  extensions:
    - ".go"
  dirs:
    - "."
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  exclude_files: []
log:
  color: true
  time: true
  level: "info"
```
