---
description: Run static analysis and linting
allowed-tools: Bash(go vet:*), Bash(golangci-lint:*)
---

Run static analysis on the codebase:

1. Run go vet:

```bash
go vet ./...
```

2. Run golangci-lint (if available):

```bash
golangci-lint run ./...
```

Report any issues found and suggest fixes.
