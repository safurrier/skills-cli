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
**Date**: 2026-03-12
**Deciders**: Alex Furrier
**Generated from**: init

---

## Context

skills-cli needs a primary implementation stack. The prototype was in Python (~400 lines), but the production tool should be a single binary with no runtime dependencies for easy distribution.

## Decision

**Stack**: Go

- **go** as the compiler and module manager
- **gofumpt** for formatting (stricter than gofmt)
- **golangci-lint** for linting (17 linters)
- **go vet** for static analysis
- **go test** for testing
- **cobra** for CLI framework
- **gopkg.in/yaml.v3** for YAML parsing

## Consequences

**Positive**:

- Single binary distribution — no Python/Node runtime needed
- Fast startup — important for CLI tools invoked frequently
- Strong stdlib for file I/O and path handling
- cobra gives shell completion and help generation for free

**Negative / Trade-offs**:

- YAML parsing is less ergonomic than Python's pyyaml
- More verbose than the Python prototype (~1600 lines vs ~400)

## Alternatives Considered

| Alternative | Reason not chosen |
|---|---|
| Python (production) | Requires Python runtime on user's machine; prototype proved the concept but distribution is harder |
| Rust | Heavier build toolchain; Go is sufficient for this use case |
