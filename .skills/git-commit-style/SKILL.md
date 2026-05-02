---
name: git-commit-style
description: Use when committing code changes in this repository.
metadata:
  short-description: Consistent commit messages and grouping
---

# Git Commit Style

Use this skill when the user wants to commit changes.

## Required Workflow

1. Inspect the full diff and current status.
2. Group files by semantic change, not by convenience.
3. Keep unrelated changes in separate commits.
4. If a change is already independent, commit it alone.
5. If a diff is ambiguous, pause and ask before mixing concerns.
6. Commit with the format:

   `type(scope): short imperative summary`

## Commit Message Rules

- Use lowercase `type`
- Use a short `scope` when it improves clarity
- Use imperative mood
- Keep the summary brief
- Do not end the summary with a period
- Prefer stable, descriptive wording over clever wording

## Recommended Types

- `feat` for new behavior
- `fix` for bug fixes
- `chore` for maintenance, editor config, tooling, dependency, or repo hygiene changes
- `test` for test-only changes
- `refactor` for code changes without behavior change
- `docs` for documentation only

## Validation Checklist

Before committing, confirm:

- each commit has a single purpose
- the working tree is clean after the last commit
- the commit message follows the chosen format
- the diff for each commit is readable on its own
