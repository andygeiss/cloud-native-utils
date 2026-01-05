# AGENTS.md

## Overview

This repository uses a set of specialized AI agents to maintain documentation, enforce conventions, and assist with development. Each agent has a specific responsibility and operates according to documented ground-truth files.

All agent definitions live in `.github/agents/`.

---

## Agent Index

| Agent | File | Role | When to Use |
|-------|------|------|-------------|
| **coding-assistant** | `.github/agents/coding-assistant.md` | Senior engineer implementing code changes | Feature development, bug fixes, refactoring |
| **CONTEXT-maintainer** | `.github/agents/CONTEXT-maintainer.md` | Architect maintaining `CONTEXT.md` | Architecture changes, convention updates |
| **README-maintainer** | `.github/agents/README-maintainer.md` | Documentation specialist for `README.md` | User-facing doc updates, feature descriptions |
| **VENDOR-maintainer** | `.github/agents/VENDOR-maintainer.md` | Vendor library curator for `VENDOR.md` | Dependency changes, integration patterns |
| **AGENTS-maintainer** | `.github/agents/AGENTS-maintainer.md` | Agent orchestrator maintaining this file | Agent additions, role changes |

---

## Agent Details

### coding-assistant

- **Path**: `.github/agents/coding-assistant.md`
- **Role**: Senior software engineer and documentation-oriented coding agent
- **Responsibilities**:
  - Implement features and fixes following repository conventions
  - Reuse existing patterns and vendor utilities before creating new ones
  - Enforce architecture rules from `CONTEXT.md`
- **Ground Truth**: `CONTEXT.md` → `README.md` → `VENDOR.md`
- **Use When**: Writing or modifying code, implementing features, fixing bugs

---

### CONTEXT-maintainer

- **Path**: `.github/agents/CONTEXT-maintainer.md`
- **Role**: Senior architect maintaining the authoritative project context
- **Responsibilities**:
  - Keep `CONTEXT.md` accurate and in sync with actual codebase
  - Document architecture, conventions, and directory contracts
  - Ensure high signal-per-token for AI agents and developers
- **Ground Truth**: Actual codebase structure and conventions
- **Use When**: Architecture changes, new packages added, conventions evolve

---

### README-maintainer

- **Path**: `.github/agents/README-maintainer.md`
- **Role**: Documentation specialist for human-facing README
- **Responsibilities**:
  - Keep `README.md` accurate and aligned with codebase
  - Document features, installation, and usage examples
  - Ensure consistency with `CONTEXT.md`
- **Ground Truth**: `CONTEXT.md` (for architecture) → actual codebase
- **Use When**: New features documented, installation steps change, usage examples needed

---

### VENDOR-maintainer

- **Path**: `.github/agents/VENDOR-maintainer.md`
- **Role**: Vendor library curator and integration specialist
- **Responsibilities**:
  - Document external dependencies and their usage patterns
  - Provide agent-friendly guidance on when/how to use vendor libraries
  - Discourage duplicate implementations of vendor functionality
- **Ground Truth**: `CONTEXT.md` → `README.md` → vendor official docs
- **Use When**: Dependencies added/removed, integration patterns change, version upgrades

---

### AGENTS-maintainer

- **Path**: `.github/agents/AGENTS-maintainer.md`
- **Role**: Agent orchestrator maintaining this index
- **Responsibilities**:
  - Keep `AGENTS.md` in sync with `.github/agents/*.md` files
  - Document agent roles, collaboration patterns, and ground-truth hierarchies
- **Ground Truth**: `.github/agents/*.md` files → `CONTEXT.md` → `README.md` → `VENDOR.md`
- **Use When**: New agents added, agent responsibilities change

---

## Ground Truth Hierarchy

When conflicts arise between documents:

| Concern | Authoritative Source |
|---------|---------------------|
| Architecture & conventions | `CONTEXT.md` |
| Human-facing description | `README.md` |
| Vendor usage & integration | `VENDOR.md` |
| Agent definitions | `.github/agents/*.md` |

---

## Agent Collaboration

```
┌─────────────────────┐
│  coding-assistant   │ ◀── implements code changes
└─────────────────────┘
          │
          ▼ triggers updates to
┌─────────────────────┐     ┌─────────────────────┐
│ CONTEXT-maintainer  │     │  README-maintainer  │
│  (architecture)     │     │  (user docs)        │
└─────────────────────┘     └─────────────────────┘
          │                           │
          ▼                           ▼
     CONTEXT.md                  README.md
          │                           │
          └───────────┬───────────────┘
                      ▼
              VENDOR-maintainer
              (dependency docs)
                      │
                      ▼
                 VENDOR.md
```

### Typical Workflow

1. **coding-assistant** implements a feature or adds a new package
2. **CONTEXT-maintainer** updates `CONTEXT.md` if architecture/conventions change
3. **README-maintainer** updates `README.md` if user-facing documentation needs changes
4. **VENDOR-maintainer** updates `VENDOR.md` if new dependencies are introduced
5. **AGENTS-maintainer** updates this file if agent roles or definitions change

---

## For External Tools

AI tools (Zed, MCP servers, VS Code Copilot, etc.) should:

1. **Locate agents** at `.github/agents/*.md`
2. **Select the appropriate agent** based on the task type (see Agent Index table)
3. **Respect ground-truth hierarchy** when resolving conflicts
4. **Follow the collaboration model** when changes span multiple concerns

---

## Current Dependencies

This repository (`cloud-native-utils`) is a Go module providing:

- Resilience patterns (`stability/`)
- Structured logging (`logging/`)
- Generic CRUD persistence (`resource/`)
- Message dispatching (`messaging/`)
- Security primitives (`security/`)
- HTTP middleware and utilities

See `CONTEXT.md` for full architecture and `VENDOR.md` for approved external libraries.
