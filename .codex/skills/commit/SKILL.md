---
name: commit
description: Use when asked to write a git commit for this repo. The repo-specific guidance is to use Conventional Commit format and keep each commit focused.
---

# Commit

Write commits in this format:

```text
type(scope): description
```

## Rules

- Do not add AI attribution, co-author lines, or generated-by footers.
- Use the `ai` scope only for changes that are actually about AI behavior, AI tooling, or AI-specific repo features.
- Prefer a single purpose per commit.
- Keep unrelated changes in separate commits.

## Common types

- `feat` for new behavior
- `fix` for bug fixes
- `chore` for maintenance, tooling, or repo hygiene
- `test` for test-only changes
- `refactor` for code changes without behavior change
- `docs` for documentation only

## Common scopes

- `renderer`
- `httpclient`
- `urlparser`
- `ai`
