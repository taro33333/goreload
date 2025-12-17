# Development Guide

## Prerequisites

- Go 1.21 or later
- Git
- (Optional) staticcheck for linting

## Getting Started

### Clone Repository

```bash
git clone https://github.com/taro33333/goreload.git
cd goreload
```

### Install Dependencies

```bash
go mod download
```

### Build

```bash
go build -o ./tmp/goreload ./cmd/goreload
```

### Run Tests

```bash
# All tests
go test ./...

# With verbose output
go test -v ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Lint

```bash
# go vet
go vet ./...

# staticcheck (if installed)
staticcheck ./...
```

## Project Structure

```
goreload/
├── cmd/goreload/          # CLI entry point
├── internal/              # Private packages
│   ├── config/            # Configuration
│   ├── logger/            # Logging
│   ├── builder/           # Build execution
│   ├── runner/            # Process management
│   ├── watcher/           # File watching
│   └── engine/            # Orchestration
├── docs/                  # Documentation
├── .claude/               # Claude Code config
├── .github/workflows/     # CI/CD
├── goreload.yaml          # Sample config
├── .goreleaser.yaml       # Release config
└── go.mod                 # Go modules
```

## Coding Standards

### Go Idioms

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` or `goimports` for formatting

### Naming

| Type | Convention | Example |
|------|------------|---------|
| Package | lowercase, single word | `config`, `logger` |
| Exported | PascalCase | `LoadConfig`, `Builder` |
| Unexported | camelCase | `parseCommand`, `runLoop` |
| Interface | -er suffix | `Builder`, `Runner`, `Watcher` |
| Receiver | 1-2 letters | `func (w *Watcher)` |

### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("load config: %w", err)
}

// Define sentinel errors
var ErrNotFound = errors.New("not found")

// Check specific errors
if errors.Is(err, ErrNotFound) {
    // handle
}
```

### Documentation

```go
// Package config provides configuration loading and validation.
package config

// Config represents the complete goreload configuration.
type Config struct {
    // Root is the project root directory.
    Root string `yaml:"root"`
}

// LoadWithDefaults loads configuration from path and merges with defaults.
func LoadWithDefaults(path string) (*Config, error) {
    // ...
}
```

## Testing

### Table-Driven Tests

```go
func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        cfg     Config
        wantErr error
    }{
        {
            name:    "valid config",
            cfg:     validConfig(),
            wantErr: nil,
        },
        {
            name: "empty build cmd",
            cfg: Config{
                Build: BuildConfig{Cmd: ""},
            },
            wantErr: ErrEmptyBuildCmd,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.cfg.Validate()
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Helpers

```go
func TestBuilder_Build(t *testing.T) {
    // Use t.TempDir() for temporary files
    tmpDir := t.TempDir()

    // Create test files
    goFile := filepath.Join(tmpDir, "main.go")
    os.WriteFile(goFile, []byte("package main\nfunc main(){}"), 0644)

    // Use t.Cleanup() for cleanup
    t.Cleanup(func() {
        // cleanup code
    })
}
```

### Coverage Targets

| Package | Target |
|---------|--------|
| config | 90%+ |
| logger | 95%+ |
| builder | 85%+ |
| watcher | 80%+ |
| runner | 75%+ |
| engine | 60%+ |

## Adding New Features

### 1. Design

1. Define the interface
2. Document the behavior
3. Consider error cases
4. Plan for testability

### 2. Implement

1. Create the implementation
2. Add GoDoc comments
3. Handle errors properly
4. Use context for cancellation

### 3. Test

1. Write table-driven tests
2. Test error cases
3. Test edge cases
4. Achieve coverage target

### 4. Document

1. Update relevant docs in `docs/`
2. Update `CLAUDE.md` if architecture changes
3. Update README if user-facing

## Release Process

### Version Tagging

```bash
# Create and push tag
git tag v0.1.0
git push origin v0.1.0
```

### Automated Release

GitHub Actions automatically:

1. Runs tests and linting
2. Builds binaries for all platforms
3. Creates GitHub Release
4. Publishes to Homebrew tap

### Local Testing

```bash
# Check goreleaser config
goreleaser check

# Build snapshot (no publish)
goreleaser release --snapshot --clean
```

## Debugging

### Enable Debug Logging

```yaml
# goreload.yaml
log:
  level: "debug"
```

### Build with Debug Symbols

```bash
go build -gcflags="all=-N -l" -o ./tmp/goreload ./cmd/goreload
```

### Use Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/goreload -- -c goreload.yaml
```

## CI/CD

### GitHub Actions

| Workflow | Trigger | Actions |
|----------|---------|---------|
| `ci.yaml` | Push/PR to main | Test, lint, build |
| `release.yaml` | Tag push (v*) | Release, publish |

### Required Secrets

| Secret | Purpose |
|--------|---------|
| `GITHUB_TOKEN` | Release assets |
| `HOMEBREW_TAP_TOKEN` | Homebrew tap push |

## Troubleshooting

### Common Issues

**Build fails with module error:**

```bash
go mod tidy
```

**Tests fail with permission error:**

```bash
# Check file permissions in test temp dirs
chmod +x ./tmp/test-binary
```

**Watcher not detecting changes:**

- Check `watch.dirs` includes the directory
- Check `watch.exclude_dirs` doesn't exclude it
- Check file extension is in `watch.extensions`

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [fsnotify Documentation](https://github.com/fsnotify/fsnotify)
- [Cobra Documentation](https://cobra.dev/)
