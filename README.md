# skills-cli

A skills-cli project

## Quick Start

```bash
# Install dependencies
mise run setup

# Run quality checks
mise run check

# Start development
mise run dev```

## Starting New Work

```bash
git checkout -b feat/<slug>
mise run plan -- <slug>
```

## Task Reference

| Command | Purpose |
|---------|---------|
| `mise run setup` | Install dependencies and prepare the environment |
| `mise run fmt` | Auto-format code |
| `mise run lint` | Run lint checks (non-modifying) |
| `mise run typecheck` | Run static type analysis |
| `mise run test` | Run unit tests |
| `mise run build` | Build artifacts |
| `mise run check` | Fast quality gate (fmt + lint + typecheck + test) |
| `mise run dev` | Start local development || `mise run ci` | CI entrypoint (= check) |
| `mise run plan -- <slug>` | Create a plan directory for a unit of work |
| `mise run verify` | Heavy validation (integration, docker, security) |

## Project Structure

```
skills-cli/
├── cmd/                    # Entry points
├── internal/               # Private packages
│   └── app/                # Core application logic
├── .mise.toml              # Task runner config
├── go.mod                  # Go module config
├── .golangci.yml           # Linter config
├── Dockerfile              # Multi-stage distroless build
└── README.md
```

## Development

This project uses [mise](https://mise.jdx.dev/) as the task runner. All quality gates
are accessible via `mise run <task>`.

CI calls a single entrypoint: `mise run ci`.
