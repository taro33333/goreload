---
globs: ["**/*.go"]
---

# Go Development Rules

## Code Style

- Follow Effective Go and Go Code Review Comments
- Use `gofmt` or `goimports` for formatting
- Keep functions short and focused (< 50 lines preferred)
- Use early returns to reduce nesting

## Naming Conventions

- Package names: lowercase, single word, no underscores
- Exported names: PascalCase
- Unexported names: camelCase
- Receiver names: 1-2 letters (e.g., `func (w *Watcher)`)
- Interface names: -er suffix (e.g., `Reader`, `Builder`)
- Avoid stuttering (e.g., `config.Config` is OK, `config.ConfigStruct` is not)

## Error Handling

- Always check errors immediately after function calls
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Define sentinel errors: `var ErrNotFound = errors.New("not found")`
- Return errors to callers; log only at the top level
- Never use panic except in init() for unrecoverable errors

## Documentation

- All exported types, functions, and methods must have GoDoc comments
- Comments should start with the name being documented
- Example: `// Builder executes build commands.`

## Concurrency

- Always pass `context.Context` as the first parameter
- Use `select` with context cancellation
- Document goroutine lifecycle and ownership
- Close channels from the sender side
- Use `sync.Once` for one-time initialization

## Testing

- Table-driven tests for multiple cases
- Use `t.Run()` for subtests
- Test files in the same package (white-box testing)
- Use `t.TempDir()` for temporary files
- Mock external dependencies via interfaces

## Project Structure

- `cmd/` for main packages
- `internal/` for private packages
- Keep package dependencies acyclic
- Small interfaces (1-3 methods)
