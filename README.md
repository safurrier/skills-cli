# skills-cli

A Go CLI for the [Agent Skills](https://agentskills.io) open standard. Discover, validate, inspect, and invoke SKILL.md files in your repositories.

## Install

```bash
go install github.com/safurrier/skills-cli/cmd@latest
```

Or build from source:

```bash
git clone https://github.com/safurrier/skills-cli.git
cd skills-cli
go build -o bin/skills ./cmd/...
```

## Usage

```bash
# List all skills in the current directory
skills list

# List skills in a specific directory
skills list --dir /path/to/repo

# Inspect a skill's metadata, inputs, and resources
skills inspect code-review --dir /path/to/repo

# Validate all skills
skills init --dir /path/to/repo

# Invoke a skill (default: pointer prompt)
skills run code-review --path src/ --base main

# Invoke with full rendered content
skills run code-review --path src/ --render

# Pass freeform context after --
skills run code-review --path src/ -- "Focus on security"
```

### Pointer vs Render

By default, `skills run` outputs a **pointer prompt** that tells the agent where the skill is located:

```
# Skill: code-review

**Location**: /path/to/code-review/SKILL.md
**Directory**: /path/to/code-review/

> Review code changes and produce structured findings.

Load and follow the instructions in the skill file above.

## Parameters
- **path**: src/
- **base**: main
```

This is ideal for agents that can read files themselves. Use `--render` for the full inline content.

## Agent Skills Standard

Skills live in `.agents/skills/`, `.claude/skills/`, or `skills/` directories. Each skill is a directory with a `SKILL.md` containing YAML frontmatter.

### CLI Extensions

skills-cli extends the spec via `metadata` keys (spec-compliant):

```yaml
---
name: my-skill
description: A skill that does something
metadata:
  skills-cli.render: template
  skills-cli.inputs: |
    - name: path
      type: string
      required: true
---
```

## Development

```bash
mise run setup      # install tools
mise run check      # fmt + lint + typecheck + test
```

## Project Structure

```
skills-cli/
├── cmd/main.go             # Entry point (cobra)
├── internal/
│   ├── cmd/                # CLI commands: list, inspect, init, run
│   └── skill/              # Domain: model, parser, discover, render, pointer
├── testdata/sample-skills/ # Test fixtures
└── go.mod                  # github.com/safurrier/skills-cli
```
