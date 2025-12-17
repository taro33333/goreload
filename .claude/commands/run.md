---
description: Run goreload in the current directory
allowed-tools: Bash(./tmp/goreload:*)
---

Run the goreload hot reload tool:

```bash
./tmp/goreload
```

This will:

1. Watch for file changes in the configured directories
2. Automatically rebuild when .go files change
3. Restart the application after successful builds

To stop, press Ctrl+C.
