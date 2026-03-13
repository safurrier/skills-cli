---
id: skills-cli-architecture
title: skills-cli Architecture
description: >
  Durable system record for skills-cli: invariants, principles, cross-cutting
  workflows, and architectural decisions. Read before implementing.
index:
  - id: system-overview
    keywords: [components, data-flows, trust-boundaries, modules]
  - id: invariants-boundaries
    keywords: [invariants, security, migrations, worktree-safety, observability, correlation-id, idempotent]
  - id: principles-patterns
    keywords: [thin-handlers, timeouts, structured-logs, typed-utilities, principles]
  - id: cross-cutting-workflows
    keywords: [validation-loop, check, verify, artifact-debugging, ci-failure]
  - id: decisions
    keywords: [adrs, decisions, truth-hierarchy, decisions-dir]
  - id: where-human-thought-goes
    keywords: [human-ownership, agent-ownership, promotion-rules]
---

# skills-cli — Architecture

> This document is the durable system of record for invariants, principles, and decisions.
> Keep it updated as the system evolves. Agents: read this before implementing.

## 1. System Overview

**Stack**: go

### Components

| Component | Purpose |
|---|---|
| `cmd/main.go` | Entry point — delegates to cobra root command |
| `internal/cmd/` | CLI layer — cobra subcommands (list, inspect, init, run) |
| `internal/skill/model.go` | Domain types: Skill, Properties, RenderInput |
| `internal/skill/parser.go` | YAML frontmatter parsing + spec validation |
| `internal/skill/discover.go` | Directory scanning across spec-defined paths |
| `internal/skill/render.go` | Full inline rendering (template + passthrough modes) |
| `internal/skill/pointer.go` | Pointer output mode (location + params, no content) |

### Primary Data Flows

1. **Discovery**: Scan `.agents/skills/`, `.claude/skills/`, `skills/`, and root children for SKILL.md files
2. **Parse**: Split frontmatter from body, validate against spec constraints, extract `skills-cli.*` extensions from metadata
3. **Invoke**: Either pointer mode (default — outputs location + params) or render mode (fills templates, builds parameter blocks)

### Trust Boundaries

- Filesystem input: SKILL.md files are untrusted user content — YAML parsing uses safe loader, field lengths are bounded

---

## 2. Goals / Non-Goals

**Goals** — what this system optimizes for:

- Spec compliance: match agentskills.io validation rules exactly
- CLI ergonomics for humans invoking skills
- Pointer-first: agents read the skill file themselves rather than receiving rendered content

**Non-Goals** — explicit exclusions:

- No agent runtime or execution engine
- No remote skill registry or fetching
- No complex template engine (Jinja2, Go templates)

---

## 3. Invariants & Boundaries

> These rules prevent catastrophic mistakes. Violating them requires an ADR.

### Spec Compliance

- SKILL.md parsing must match the [agentskills.io specification](https://agentskills.io/specification) exactly.
- `skills-cli.*` metadata keys are the only extension mechanism — no custom frontmatter fields outside `metadata`.
- Discovery paths follow the spec: `.agents/skills/`, `.claude/skills/`, plus `skills/` as a pragmatic addition.

### Worktree Safety

- All tests must pass from a **clean checkout or Git worktree**.
- No reliance on absolute paths, mutable global state, or undeclared local artifacts.
- Build artifacts must be deterministic or listed in `.gitignore`.

> Security & Auth, Data & Migrations, and Traceability sections are not applicable — this is a pure CLI tool with no network, database, or service components.

---

## 4. Principles & Preferred Patterns

> Durable preferences. If a rule is critical and objective, enforce it in CI instead.

- **Domain logic in `internal/skill/`, CLI wiring in `internal/cmd/`** — cobra commands are thin wrappers.
- **Pointer-first** — default output points agents to the skill file rather than inlining content.
- **Spec-compliant extensions** — all CLI-specific behavior uses `metadata` with `skills-cli.*` keys, never custom frontmatter fields.
- **Simple template substitution** — `{{key}}` only, no template engine dependency.

---

## 5. Cross-Cutting Workflows

### Validation Loop

CI defines "passing." Local commands mirror CI exactly:

```bash
mise run check    # fast: fmt + lint + typecheck + unit tests  (on push)
mise run verify   # heavy: integration, security, docker        (on PR)
```

CI failures emit test artifacts to `test-results/` (uploaded by `actions/upload-artifact`).

### Artifact-First Debugging

When CI fails:

1. Open the GitHub Actions run summary → **Artifacts** → `test-results`.
2. **Python**: `test-results/junit.xml` — structured test output with failure details.
3. **Go**: `test-results/go-test.txt` — verbose test log.
4. Integration/docker failures: check logs in `test-results/`.

---

## 6. Decisions (ADRs)

Global architectural decisions live in `docs/decisions/`. Include only decisions that
materially constrain work.

**Truth hierarchy** (highest to lowest authority):

1. CI/tooling enforcement
2. ADRs in `docs/decisions/`
3. This document (`docs/architecture.md`)
4. Module-level docs and `AGENTS.md`

**Index**:

| ADR | Decision |
|---|---|
| [0001-stack-choice](decisions/0001-stack-choice.md) | Go stack for single-binary CLI distribution |
| [0002-pointer-default](decisions/0002-pointer-default.md) | Default `run` outputs pointer prompt, not rendered content |

---

## 7. Module Map

<!-- For single-project repos, list major internal modules.
     For apps workspace, list apps and packages with links. -->

| Module | Purpose | Docs |
|---|---|---|
| `internal/skill` | Domain logic: types, parsing, discovery, rendering | — |
| `internal/cmd` | CLI commands: cobra subcommands and flag parsing | — |
| `testdata/sample-skills` | Test fixtures: 3 sample skills from the RFC prototype | — |

---

## 8. Where Human Thought Goes

**Humans define**:

- Invariants and boundaries (section 3)
- Decision records (ADRs)
- What gets enforced in CI/tooling
- Observability minimums
- Worktree compatibility rules

**Agents own**:

- Implementation following these constraints
- Completing the validation loop (`mise run check` must pass before committing)

**Promotion rules**:

- Repeated ambiguity → document here or in module docs
- Repeated objective failure → encode in CI/tooling
- Repeated preference debate → optional Skill plugin (`.agent/skills/`)
