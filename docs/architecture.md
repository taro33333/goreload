# Architecture

## Overview

goreload follows a modular architecture with clear separation of concerns. Each component has a single responsibility and communicates through well-defined interfaces.

## Directory Structure

```
goreload/
├── cmd/goreload/
│   └── main.go              # CLI entry point, signal handling
├── internal/
│   ├── config/
│   │   ├── config.go        # Configuration structures, validation
│   │   └── loader.go        # YAML loading, default values
│   ├── logger/
│   │   └── logger.go        # Structured logging, color output
│   ├── builder/
│   │   └── builder.go       # Build command execution
│   ├── runner/
│   │   └── runner.go        # Process lifecycle management
│   ├── watcher/
│   │   ├── watcher.go       # File system watching
│   │   └── filter.go        # Path/extension filtering
│   └── engine/
│       └── engine.go        # Orchestration, main loop
├── docs/                    # Documentation
├── .claude/                 # Claude Code configuration
└── .github/workflows/       # CI/CD
```

## Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI (main.go)                       │
│  - Parse flags                                              │
│  - Signal handling (SIGINT, SIGTERM)                        │
│  - Context management                                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Engine (engine.go)                     │
│  - Orchestrates watch-build-run cycle                       │
│  - Coordinates all components                               │
│  - Main event loop                                          │
└─────────────────────────────────────────────────────────────┘
          │              │              │              │
          ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    Config    │ │   Watcher    │ │   Builder    │ │    Runner    │
│              │ │              │ │              │ │              │
│ - Load YAML  │ │ - fsnotify   │ │ - Execute    │ │ - Start      │
│ - Validate   │ │ - Debounce   │ │   build cmd  │ │ - Stop       │
│ - Defaults   │ │ - Filter     │ │ - Capture    │ │ - Restart    │
│              │ │              │ │   output     │ │ - SIGINT/    │
│              │ │              │ │              │ │   SIGKILL    │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
                       │
                       ▼
               ┌──────────────┐
               │    Filter    │
               │              │
               │ - Extensions │
               │ - Glob match │
               │ - Exclude    │
               └──────────────┘
```

## Core Components

### Config (`internal/config/`)

**Responsibility:** Load, validate, and provide configuration.

**Key Types:**

```go
type Config struct {
    Root   string      // Project root directory
    TmpDir string      // Build artifacts directory
    Build  BuildConfig // Build settings
    Watch  WatchConfig // Watch settings
    Log    LogConfig   // Log settings
}

type BuildConfig struct {
    Cmd       string        // Build command
    Bin       string        // Binary path
    Args      []string      // Runtime arguments
    Delay     time.Duration // Debounce delay
    KillDelay time.Duration // Shutdown grace period
}
```

**Key Functions:**

| Function | Description |
|----------|-------------|
| `LoadWithDefaults(path)` | Load YAML and merge with defaults |
| `Default()` | Get default configuration |
| `Validate()` | Validate configuration values |
| `WriteDefault(path)` | Generate default config file |

### Logger (`internal/logger/`)

**Responsibility:** Structured logging with color and level support.

**Interface:**

```go
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    SetOutput(w io.Writer)
    SetLevel(level Level)
}
```

**Features:**

- Color-coded output (configurable)
- Timestamp prefix (configurable)
- Log level filtering
- Thread-safe output

### Builder (`internal/builder/`)

**Responsibility:** Execute build commands and manage artifacts.

**Interface:**

```go
type Builder interface {
    Build(ctx context.Context) Result
    Clean() error
}

type Result struct {
    Success  bool
    Output   string
    Duration time.Duration
    Error    error
}
```

**Features:**

- Context-aware (cancellable builds)
- Stdout/stderr capture
- Duration tracking
- Automatic tmp directory creation

### Runner (`internal/runner/`)

**Responsibility:** Manage application process lifecycle.

**Interface:**

```go
type Runner interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Restart(ctx context.Context) error
    Running() bool
}
```

**Graceful Shutdown:**

1. Send `SIGINT` to process group
2. Wait for `KillDelay`
3. Send `SIGKILL` if still running
4. Clean up process resources

**Process Group:**

Uses `Setpgid: true` to create a process group, ensuring child processes are also terminated.

### Watcher (`internal/watcher/`)

**Responsibility:** Watch file system for changes.

**Interface:**

```go
type Watcher interface {
    Start(ctx context.Context) error
    Events() <-chan Event
    Errors() <-chan error
    Close() error
}

type Event struct {
    Path string
    Op   Op
    Time time.Time
}
```

**Features:**

- Recursive directory watching
- Debouncing (coalesce rapid changes)
- Automatic new directory detection
- Exclusion patterns

### Filter (`internal/watcher/`)

**Responsibility:** Determine which files trigger rebuilds.

**Interface:**

```go
type Filter interface {
    Match(path string) bool
}
```

**Filtering Logic:**

1. Check if path is in excluded directory
2. Check if filename matches excluded patterns (glob)
3. Check if extension is in watch list
4. Return true if all checks pass

### Engine (`internal/engine/`)

**Responsibility:** Orchestrate the watch-build-run cycle.

**Main Loop:**

```
┌─────────────────────────────────────────────┐
│                   Start                      │
└─────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────┐
│            Initial Build & Run               │
└─────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────┐
│              Wait for Event                  │◄────┐
│  - File change event                         │     │
│  - Error event                               │     │
│  - Context cancellation                      │     │
└─────────────────────────────────────────────┘     │
                      │                              │
          ┌──────────┴──────────┐                   │
          ▼                     ▼                   │
   ┌────────────┐        ┌────────────┐            │
   │ File Event │        │  Shutdown  │            │
   └────────────┘        └────────────┘            │
          │                     │                   │
          ▼                     ▼                   │
   ┌────────────┐        ┌────────────┐            │
   │ Stop Proc  │        │ Stop Proc  │            │
   │ Build      │        │ Cleanup    │            │
   │ Start Proc │        │ Exit       │            │
   └────────────┘        └────────────┘            │
          │                                         │
          └─────────────────────────────────────────┘
```

## Data Flow

### Configuration Flow

```
goreload.yaml → LoadWithDefaults() → Validate() → Config struct
                      ↑
               Default values
```

### Event Flow

```
File System Change
       │
       ▼
   fsnotify
       │
       ▼
   Filter.Match()
       │
   ┌───┴───┐
   │ false │ → (ignored)
   └───────┘
       │ true
       ▼
   Debounce Timer
       │
       ▼
   Engine Event Channel
       │
       ▼
   Stop → Build → Start
```

## Concurrency Model

### Goroutines

| Goroutine | Owner | Purpose |
|-----------|-------|---------|
| Main | CLI | Signal handling, context |
| Watcher loop | Watcher | fsnotify event processing |
| Process wait | Runner | Wait for process exit |

### Synchronization

| Component | Mechanism | Purpose |
|-----------|-----------|---------|
| Logger | `sync.Mutex` | Thread-safe output |
| Runner | `sync.Mutex` | Process state protection |
| Watcher | Channels | Event communication |
| Engine | `sync.Mutex` | Running state |

### Context Usage

All long-running operations accept `context.Context` for cancellation:

```go
func (e *Engine) Run(ctx context.Context) error
func (b *Builder) Build(ctx context.Context) Result
func (r *Runner) Start(ctx context.Context) error
func (r *Runner) Stop(ctx context.Context) error
func (w *Watcher) Start(ctx context.Context) error
```

## Error Handling

### Error Types

| Type | Example | Handling |
|------|---------|----------|
| Configuration | Invalid YAML | Fatal, exit |
| Build | Compile error | Log, continue watching |
| Runtime | Process crash | Log, wait for next change |
| System | File permission | Log error, continue |

### Sentinel Errors

```go
var (
    ErrEmptyBuildCmd    = errors.New("build command cannot be empty")
    ErrEmptyBin         = errors.New("binary path cannot be empty")
    ErrInvalidDelay     = errors.New("delay must be positive")
    ErrInvalidKillDelay = errors.New("kill_delay must be positive")
    ErrInvalidLogLevel  = errors.New("invalid log level")
    ErrNoExtensions     = errors.New("at least one extension required")
    ErrNoDirs           = errors.New("at least one directory required")
)
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/fsnotify/fsnotify` | Cross-platform file watching |
| `gopkg.in/yaml.v3` | YAML configuration parsing |
| `github.com/spf13/cobra` | CLI framework |
| `github.com/fatih/color` | Terminal color output |
