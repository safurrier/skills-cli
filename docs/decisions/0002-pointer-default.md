---
id: skills-cli-adr-0002
title: ADR 0002 — Pointer Output as Default Run Mode
description: >
  Default `skills run` outputs a location-based pointer prompt instead of
  rendering full skill content inline. Agents read the file themselves.
index:
  - id: decision
    keywords: [pointer, render, default, run, output, mode]
  - id: consequences
    keywords: [trade-offs, agent, context-window, source-of-truth]
---

# ADR 0002: Pointer Output as Default Run Mode

**Status**: Accepted
**Date**: 2026-03-12
**Deciders**: Alex Furrier
**Generated from**: agent-session

---

## Context

The Python prototype's `skills run` command rendered the full skill content as a prompt — header, body with template substitution, parameter block, freeform context, and resource references. This works but has drawbacks:

1. **Context bloat**: Full skill content is dumped into the agent's context window, consuming tokens on content the agent could read on-demand
2. **Stale renders**: The rendered output is a snapshot — if the skill file is updated, previously rendered prompts are outdated
3. **Agent capability mismatch**: Modern agents (Claude Code, Codex, Cursor) can read files directly; they don't need content pre-rendered into their prompt

## Decision

Default `skills run <name>` outputs a **pointer prompt** that tells the agent where the skill is located, includes the description and parameters, and instructs the agent to read the file:

```markdown
# Skill: code-review

**Location**: /path/to/code-review/SKILL.md
**Directory**: /path/to/code-review/

> Review code changes and produce structured findings.

Load and follow the instructions in the skill file above.

## Parameters
- **path**: src/
- **base**: main
```

Full inline rendering is still available via `--render` flag for non-agent use cases (piping to `claude --pipe`, embedding in prompts, etc.).

## Consequences

**Positive:**

- Smaller context footprint — pointer is ~10 lines vs full render which can be 50+
- Source of truth stays at the file — no stale snapshots
- Agents read the skill when they need it, leveraging progressive disclosure
- Parameters and freeform context are still passed inline (they're invocation-specific)

**Negative / Trade-offs:**

- Agents that cannot read files (pure API clients) get less useful output by default — they need `--render`
- Extra file read by the agent adds a tool call vs having content pre-loaded

## Alternatives Considered

| Alternative | Reason not chosen |
|---|---|
| Full render as default | Original prototype behavior; works but bloats context and creates stale snapshots |
| XML `<available_skills>` format (skills-ref style) | Designed for skill discovery, not invocation; doesn't carry parameters |
| Separate `pointer` subcommand | Adds CLI surface area; a flag (`--render`) is simpler |
