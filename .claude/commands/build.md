---
description: Build the goreload binary
allowed-tools: Bash(go build:*), Bash(go mod:*)
---

Build the goreload binary:

1. First, ensure dependencies are up to date:

```bash
go mod tidy
```

2. Build the binary:

```bash
go build -o ./tmp/goreload ./cmd/goreload
```

3. Verify the build:

```bash
./tmp/goreload version
```
