---
name: code-review
description: Review code changes and produce structured findings. Use when reviewing a path, diff, or branch for quality issues, bugs, style violations, or security concerns.
compatibility: Designed for Claude Code or similar agentic coding tools
metadata:
  skills-cli.render: template
  skills-cli.inputs: |
    - name: path
      type: string
      required: false
      description: File or directory path to review
    - name: base
      type: string
      required: false
      description: Base branch or commit to diff against
    - name: format
      type: string
      default: markdown
      description: Output format (markdown or json)
---

# Code Review

You are performing a structured code review. Follow these steps precisely.

## Process

1. **Understand scope**: Identify the files/changes being reviewed (from `path` or diff context)
2. **Read the code**: Load each relevant file and understand its purpose
3. **Analyze systematically**: Check for issues in these categories:
   - Correctness: Logic errors, off-by-one, null handling, race conditions
   - Security: Injection, auth bypasses, secrets in code, input validation
   - Performance: N+1 queries, unnecessary allocations, blocking I/O
   - Maintainability: Naming, complexity, missing tests, documentation
4. **Report findings**: Output findings grouped by severity (critical, warning, info)

## Output Format

```
## Code Review: <path or description>

### Critical
- [FILE:LINE] Description of issue

### Warnings
- [FILE:LINE] Description of issue

### Info / Style
- [FILE:LINE] Suggestion

### Summary
N issues found (X critical, Y warnings, Z info)
```

## Guidelines

- Be specific: always include file paths and line numbers
- Be constructive: explain *why* something is an issue
- Distinguish bugs from style preferences
- Note what looks good, not just problems
- If `base` is provided, focus on changed lines only
