---
id: skills-cli-adr-0001
title: ADR 0001 — Stack Choice for skills-cli
description: >
  Records the go stack selection decision for skills-cli,
  including rationale, trade-offs, and alternatives considered.
index:
  - id: decision
    keywords: [stack, choice, python, go, tools, rationale]
  - id: consequences
    keywords: [trade-offs, positive, negative, accepted]
  - id: alternatives-considered
    keywords: [alternatives, rejected, comparison]
---

# ADR 0001: Stack Choice for skills-cli

**Status**: Accepted
**Date**: <!-- YYYY-MM-DD -->
**Deciders**: <!-- names or team -->
**Generated from**: init

---

## Context

skills-cli requires a primary implementation stack for building, testing, and
deploying the application. The choice constrains tooling, CI configuration, and
contributor onboarding.

## Decision

**Stack**: go

The Go stack uses:
- **go** (1.23) as the compiler and module manager
- **gofumpt** for formatting (stricter than gofmt)
- **golangci-lint** for linting (17 linters)
- **go vet** for static analysis
- **go test** for testing

## Consequences

**Positive**:

- Standard tooling with strong ecosystem support.
- Consistent quality gates via `mise run check`.
- Reproducible builds via mise tool version pinning.

**Negative / Trade-offs**:

- <!-- list accepted trade-offs, e.g., "Go compilation adds ~30s to cold CI runs" -->

## Alternatives Considered

<!-- List stacks that were considered but not chosen, and why -->

| Alternative | Reason not chosen |
|---|---|
| <!-- alt --> | <!-- reason --> |
