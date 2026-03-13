---
name: write-spec
description: Write a technical specification document for a proposed feature, API, or system component. Use when asked to spec out, design, or document a new capability.
metadata:
  skills-cli.render: template
  skills-cli.inputs: |
    - name: topic
      type: string
      required: true
      description: The feature or system to specify
    - name: audience
      type: string
      required: false
      default: engineers
      description: Target audience (engineers, product, executives)
---

# Technical Specification Writer

You are writing a technical specification. Produce a well-structured document.

## Spec Template

```markdown
# <Topic> — Technical Specification

## Status
Draft | Review | Approved

## Overview
One paragraph summary of what this is and why it exists.

## Background & Motivation
What problem does this solve? Why now?

## Goals
- Goal 1
- Goal 2

## Non-Goals
- Explicitly out of scope

## Design
### Architecture
### API / Interface
### Data Model
### Error Handling

## Implementation Plan
Phases and milestones

## Open Questions
Unresolved issues

## Alternatives Considered
What else was considered and why it was rejected
```

## Guidelines

- Be precise about interfaces and contracts
- Acknowledge uncertainty with explicit open questions
- Keep the Overview readable by non-engineers
- Tailor depth to the `audience` parameter
