---
description: Run static analysis and linting
allowed-tools: Bash(go vet:*), Bash(staticcheck:*)
---

Run static analysis on the codebase:

1. Run go vet:

```bash
go vet ./...
```

2. Run staticcheck:

```bash
staticcheck ./...
```

Report any issues found and suggest fixes.
