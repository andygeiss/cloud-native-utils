# Specification

## Job

A Go developer imports one package from this module — a circuit breaker, a CRUD store,
a session middleware — and gets that one part, not a framework.

## Why

Every service rebuilds the same small parts: retry a call, hash a password, put a value
somewhere and read it back. A framework hands you all of them and takes the shape of
your project in return. These packages are the parts on their own, so a service takes
what it needs and stays plain Go.

## Guardrails

- **One concern per package.** A new concern is a new package. Widening an existing one
  so it covers two things is the thing to avoid.
- **Imports run one way.** `service` is the base; `web` is the top and nothing imports
  it. Today the module-internal edges are `efficiency`, `mcp` and `stability` onto
  `service`; `resource` onto `efficiency`; `messaging` onto `env`, `service` and
  `stability`; and `web` onto `efficiency`, `env`, `mcp`, `resource` and `security`.
  A new edge is a design decision, not a convenience.
- **The dependencies and the bent rules are listed once**, under *Baseline deviations*
  in [README.md](README.md). Adding either means adding an entry there.
- **The named decisions live in [CLAUDE.md](CLAUDE.md)** — read it before changing a
  pattern it records.

## Done means

- The boxes of the baseline's `checklists/library.md` are walked.
- `make check` is green before every commit; `make ci` is green before every push.
