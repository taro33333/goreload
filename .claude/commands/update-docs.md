---
description: Update documentation after code changes
allowed-tools: Read, Write, Edit, Grep, Glob
---

After making code changes, update the relevant documentation.

## Steps

1. First, identify what changed by examining recent modifications
2. Determine which docs need updating based on CLAUDE.md rules:

| Change Type | Update Files |
|-------------|--------------|
| Config option added/changed | `docs/configuration.md` |
| CLI command/flag added | `docs/cli.md` |
| Interface changed | `docs/architecture.md` |
| Installation changed | `docs/getting-started.md`, `README.md` |

3. Read the relevant source files to understand the changes
4. Update the documentation to reflect the changes
5. Ensure consistency across all affected documents

## Important

- Keep documentation concise and accurate
- Include examples where helpful
- Update tables and code blocks as needed
- Maintain consistent formatting with existing docs
