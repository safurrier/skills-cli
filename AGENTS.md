# skills-cli

Go CLI for the [Agent Skills](https://agentskills.io) open standard. Discovers, validates, inspects, and invokes SKILL.md files in repositories.

## WHY

The Agent Skills spec defines a file format and directory structure for agent skills but has no CLI layer for humans. This tool bridges that gap. Key design choice: default `run` outputs a pointer prompt (skill location + params) instead of rendering content inline — agents read the file themselves.

**Done means**: all `mise run check` gates pass; intended behavior is covered by tests.

Correctness invariants live in [`SPEC.md`](SPEC.md). System design and decisions in [`docs/architecture.md`](docs/architecture.md).

## WHAT

```
skills-cli/
├── cmd/main.go             # Entry point (cobra)
├── internal/
│   ├── cmd/                # CLI commands: list, inspect, init, run
│   └── skill/              # Domain: model, parser, discover, render, pointer
├── testdata/sample-skills/ # Test fixtures (3 sample skills)
├── .mise.toml              # Task runner config
├── go.mod                  # Go module (github.com/safurrier/skills-cli)
├── .golangci.yml           # Linter config (17 linters)
└── Dockerfile              # Multi-stage distroless build
```

Key steering files:
- `AGENTS.md` — this file (steering index)
- `SPEC.md` — correctness envelope (requirements, contracts, invariants)
- `docs/architecture.md` — system description, principles, decisions
- `docs/decisions/` — architectural decision records (ADRs)
- `.agent/skills/` — optional opinionated workflow plugins

## HOW

```bash
mise run setup      # install tools and dependencies (one-time)
mise run check      # fast quality gate: fmt + lint + typecheck + test  ← before committing
mise run verify     # heavier validation: integration, security, docker  ← before merging
mise run dev        # start local development```

CI calls `mise run ci` (= `check`). Pre-commit hooks mirror CI exactly.
If `check` passes locally, CI passes.

## Stack: go

- Formatter: gofumpt (stricter than gofmt)
- Linter: golangci-lint (17 linters enabled)
- Type checker: go vet (integrated compiler checks)
- Tests: go test

## Starting Work

```bash
git checkout -b feat/<slug>
mise run plan -- <slug>    # creates .ai/plans/YYYY-MM-DD-HHmmSS-<slug>/
```

`mise run plan` refuses to run on the default branch. Slugs must be lowercase
kebab-case and unique within `.ai/plans/`.

The task scaffolds META.yaml, TODO.md, LEARNING_LOG.md, and VALIDATION.md.
Add SPEC.md or IMPLEMENTATION.md if the work is complex.

See `.ai/plans/AGENTS.md` for the full plan structure and `_example/` for a reference.

## Before Pushing

1. `mise run check` — must pass
2. `/plan-sync` — verify plan META, TODOs, and logs are current
3. `/spec-sync` — review changes against SPEC.md, capture any decisions as ADRs
4. `/context-engineering update` — update AGENTS.md if workflows or structure changed
5. `/docs-workflow update` — update docs/ if architecture or module docs affected

Steps 2-5 are fast when nothing changed. The agent reviews its own diff and skips
if no updates are needed.

## Skills

Optional workflow plugins live in `.agent/skills/`. Each `SKILL.md` describes when to
load it and what workflow it encodes. Skills are preference-based — canonical truth
stays in `docs/`.

`.claude/skills` is symlinked to `.agent/skills` for Claude Code discovery. Other
harnesses can create `<harness-config>/skills → ../.agent/skills` similarly.

## Further Reading

| Document | Purpose |
|---|---|
| `SPEC.md` | Correctness envelope — requirements, contracts, invariants |
| `.ai/plans/` | Plan directories for units of work (see `.ai/plans/AGENTS.md`) |
| `docs/architecture.md` | System description, principles, decisions |
| `docs/decisions/` | Architectural decision records |
| `README.md` | Human-oriented quick start |
