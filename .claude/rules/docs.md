---
globs: ["docs/**/*.md", "CLAUDE.md", "README.md"]
---

# Documentation Rules

## When to Update Docs

Documentation MUST be updated when:

1. **Configuration changes** - Any change to `internal/config/config.go`:
   - New config options → Add to `docs/configuration.md`
   - Changed defaults → Update `docs/configuration.md`
   - Removed options → Remove from `docs/configuration.md`

2. **CLI changes** - Any change to `cmd/goreload/main.go`:
   - New commands → Add to `docs/cli.md`
   - New flags → Add to `docs/cli.md`
   - Changed behavior → Update `docs/cli.md`

3. **Interface changes** - Any change to `internal/*/` interfaces:
   - New methods → Update `docs/architecture.md`
   - Changed signatures → Update `docs/architecture.md`

4. **Installation changes**:
   - Update `docs/getting-started.md` and `README.md`

## Documentation Standards

- Use consistent markdown formatting
- Include code examples for complex features
- Keep tables aligned and readable
- Use proper heading hierarchy (h1 → h2 → h3)
- Cross-reference related documentation

## File Responsibilities

| File | Content |
|------|---------|
| `docs/configuration.md` | All config options, defaults, validation |
| `docs/cli.md` | All commands, flags, examples |
| `docs/architecture.md` | Interfaces, data flow, components |
| `docs/getting-started.md` | Installation, quick start |
| `docs/development.md` | Contributing, testing, debugging |
| `README.md` | Overview, basic usage (user-facing) |
| `CLAUDE.md` | AI context, quick reference |
