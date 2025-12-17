---
description: Sync documentation with code changes - checks for inconsistencies
allowed-tools: Read, Grep, Glob
---

Check if documentation is in sync with code. Analyze the following:

## 1. Configuration Options

Compare `internal/config/config.go` with `docs/configuration.md`:

- Check if all Config struct fields are documented
- Check if all default values match
- Check if all validation rules are documented

## 2. CLI Commands

Compare `cmd/goreload/main.go` with `docs/cli.md`:

- Check if all commands are documented
- Check if all flags are documented
- Check if usage examples are accurate

## 3. Interfaces

Compare `internal/*/` interfaces with `docs/architecture.md`:

- Check if all public interfaces are documented
- Check if method signatures match

## 4. Report

Generate a report listing:

1. Missing documentation
2. Outdated documentation
3. Inconsistencies between code and docs

If issues found, suggest specific updates needed.
