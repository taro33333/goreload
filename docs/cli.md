# CLI Reference

[English](cli.md) | [日本語](cli_ja.md)

## Synopsis

```
goreload [flags]
goreload [command]
```

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `goreload.yaml` | Path to configuration file |
| `--help` | `-h` | | Show help for command |

## Commands

### `goreload` (default)

Run the hot reload watcher.

```bash
goreload
goreload -c ./custom-config.yaml
```

**Behavior:**

1. Loads configuration from `goreload.yaml` (or specified file)
2. Creates temporary directory if needed
3. Performs initial build
4. Starts the compiled binary
5. Watches for file changes
6. On change: stops process → rebuilds → restarts
7. Continues until interrupted (Ctrl+C)

**Exit Codes:**

| Code | Description |
|------|-------------|
| 0 | Normal termination (Ctrl+C) |
| 1 | Configuration error or fatal error |

### `goreload init`

Generate a default configuration file.

```bash
goreload init
```

**Behavior:**

- Creates `goreload.yaml` in the current directory
- Fails if file already exists (prevents accidental overwrite)

**Output:**

```
Created goreload.yaml
```

### `goreload version`

Print version information.

```bash
goreload version
```

**Output:**

```
goreload v0.1.0
  commit: abc1234
  built:  2024-01-01T00:00:00Z
```

### `goreload help`

Show help information.

```bash
goreload help
goreload help init
goreload --help
```

## Usage Examples

### Basic Usage

```bash
# Use default configuration
goreload

# Use custom configuration
goreload -c ./config/dev.yaml
```

### Web Application

```yaml
# goreload.yaml
build:
  cmd: "go build -o ./tmp/server ./cmd/server"
  bin: "./tmp/server"
  args:
    - "-port=8080"
    - "-env=development"

watch:
  extensions:
    - ".go"
    - ".html"
    - ".css"
  dirs:
    - "."
    - "templates"
  exclude_dirs:
    - "tmp"
    - "node_modules"
```

### Microservice with Make

```yaml
# goreload.yaml
build:
  cmd: "make build"
  bin: "./bin/service"

watch:
  extensions:
    - ".go"
    - ".proto"
  exclude_files:
    - "*.pb.go"
```

### Debug Mode

```yaml
# goreload.yaml
build:
  cmd: "go build -gcflags='all=-N -l' -o ./tmp/debug ."
  bin: "./tmp/debug"

log:
  level: "debug"
```

## Signal Handling

goreload handles the following signals:

| Signal | Behavior |
|--------|----------|
| `SIGINT` (Ctrl+C) | Graceful shutdown |
| `SIGTERM` | Graceful shutdown |

**Graceful Shutdown Process:**

1. Stop watching for file changes
2. Send SIGINT to running process
3. Wait for `kill_delay` duration
4. Send SIGKILL if process still running
5. Exit

## Output Format

### Log Messages

```
[TIMESTAMP] [LEVEL] message
```

Example:

```
15:04:05 [INFO] watching: /path/to/project
15:04:05 [INFO] excluding: [tmp vendor .git]
15:04:05 [INFO] building...
15:04:07 [INFO] ✓ build completed (2.15s)
15:04:07 [INFO] ✓ running ./tmp/main
```

### Log Levels

| Level | Color | Description |
|-------|-------|-------------|
| `[DEBUG]` | Gray | Verbose debugging information |
| `[INFO]` | Cyan | Normal operation messages |
| `[WARN]` | Yellow | Warnings |
| `[ERROR]` | Red | Errors |

### Status Indicators

| Indicator | Meaning |
|-----------|---------|
| `✓` | Success (build completed, process started) |
| `✗` | Failure (build failed, process error) |
