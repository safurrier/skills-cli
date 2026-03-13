---
id: skills-cli-spec
title: skills-cli Specification
description: >
  Correctness envelope — requirements, contracts, and invariants
  for skills-cli.
index:
  - id: requirements
    keywords: [must, should, may, behavioral, requirements]
  - id: interfaces
    keywords: [api, contracts, boundaries, modules, io]
  - id: invariants
    keywords: [state, safety, data, rules, boundaries, always-true]
---

# skills-cli — Specification

> This document defines the correctness envelope: what must be true about any
> valid implementation. For how the system works now, see `docs/architecture.md`.
> For how to work in this repo, see `AGENTS.md`.

## Summary

A Go CLI for the [Agent Skills](https://agentskills.io) open standard. Discovers, validates, inspects, and invokes SKILL.md files in repositories. Default invocation outputs a pointer prompt (skill location + parameters) rather than rendering content inline — agents read the file themselves.

## Goals / Non-Goals

**Goals:**

- Provide a human-facing CLI layer for the Agent Skills spec
- Support all spec-defined discovery locations (`.agents/skills/`, `.claude/skills/`, `skills/`)
- Validate SKILL.md files against the spec's field constraints
- Default to pointer output (location-based) so agents use the skill at its source of truth
- Support full inline rendering when explicitly requested

**Non-Goals:**

- No agent runtime — this is a CLI tool, not an agent framework
- No skill registry or remote fetching (local filesystem only)
- No SKILL.md authoring/scaffolding
- No Jinja2 or complex template engines — simple `{{key}}` substitution only

## Requirements

### MUST

- `skills list` discovers skills across all spec-defined paths
- `skills inspect <name>` shows full metadata, render inputs, resources, and body preview
- `skills init` validates and reports on all discovered skills
- `skills run <name>` outputs pointer prompt by default
- `skills run <name> --render` outputs full rendered content
- `skills run` validates required inputs and fails with clear errors when missing
- SKILL.md name validation matches spec: lowercase, max 64 chars, no leading/trailing hyphens, directory name must match
- `mise run check` passes
- `mise run ci` produces identical results to `mise run check`

### SHOULD

- `skills list --json` outputs structured JSON for programmatic use
- Dynamic `--key value` pairs on `run` are passed through as skill parameters
- Freeform text after `--` is included in output as additional context
- Default render inputs are applied when not explicitly provided

### MAY

- Support `NO_COLOR` environment variable for disabling ANSI output
- Support `skills validate <name>` as a per-skill validation command
- Support user-level skill scanning (e.g., `~/.config/skills/`)

## Interfaces & Contracts

**CLI surface:**

| Command | Purpose |
|---------|---------|
| `skills list [--dir DIR] [--json]` | List discovered skills |
| `skills inspect <name> [--dir DIR]` | Show skill details |
| `skills init [--dir DIR]` | Validate/scan skills |
| `skills run <name> [--render] [--key val] [-- freeform]` | Invoke a skill |

**Internal packages:**

| Package | Responsibility |
|---------|---------------|
| `internal/skill` | Domain: model, parser, discovery, render, pointer |
| `internal/cmd` | CLI: cobra commands, flag parsing |

**Task contract:**

| Command | Purpose |
|---------|---------|
| `mise run check` | Fast quality gate (fmt + lint + typecheck + test) |
| `mise run verify` | Heavy validation (integration, docker, security) |
| `mise run ci` | CI entrypoint (delegates to `check`) |

## Invariants

- All tests must pass from a **clean checkout or Git worktree**.
- CI and local checks use the **same entrypoints** — no divergence.
- SKILL.md parsing follows the [agentskills.io specification](https://agentskills.io/specification) field constraints exactly.
- `skills-cli.*` metadata keys are the extension mechanism — no custom frontmatter fields outside `metadata`.

## Acceptance

```bash
mise run check    # fast: fmt + lint + typecheck + test
mise run verify   # heavy: integration, e2e, docker, security
```
