---
description: Generate detailed test coverage report
allowed-tools: Bash(go test:*)
---

Generate a detailed test coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Optionally, generate an HTML report:

```bash
go tool cover -html=coverage.out -o coverage.html
```

Analyze coverage and identify areas that need more tests.
