---
name: commit
description: Use when asked to write a git commit for this repo. The repo-specific guidance is to use Conventional Commit format and keep each commit focused.
---

# Commit

When the user asks to commit, inspect the entire working tree and commit all changed files that belong in the repo history. Do not limit yourself to only the files touched by the current conversation.

Split the work into multiple commits grouped by the main behavior or purpose of the change. Do not collapse unrelated work into one commit.

Write commits in this format:

```text
type(scope): description
```

## Rules

- Do not add AI attribution, co-author lines, or generated-by footers.
- Use the `ai` scope only for changes that are actually about AI behavior, AI tooling, or AI-specific repo features.
- Prefer a single purpose per commit.
- Keep unrelated changes in separate commits.
- If there are multiple logical changes in the working tree, stage and commit each group separately.
- If the diff cannot be cleanly grouped, pause and ask before mixing concerns.

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
