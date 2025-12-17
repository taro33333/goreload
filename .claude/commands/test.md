---
description: Run all tests with coverage report
allowed-tools: Bash(go test:*)
---

Run all tests with verbose output and coverage:

```bash
go test ./... -v -cover
```

If tests fail, analyze the output and suggest fixes.
