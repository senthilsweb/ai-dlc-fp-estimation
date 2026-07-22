# 1. Record architecture decisions

Date: 2026-07-22

## Status

Accepted

## Context

This repo was extracted from `ctms-gtm-mono-repo`, where two forks of the same tool (`ctms-business-features-fp`, `ai-agents-provly/business-features-fp`) had diverged from a shared design without any record of why — the reasons for choices like "no build step" or "localStorage for state" existed only as tribal knowledge or had to be reverse-engineered from diffing the two forks. That cost real time during this migration (see `openspec/changes/migrate-fp-estimator-standalone/`).

## Decision

We will keep a set of Architecture Decision Records in `docs/adr/`, one per architecturally significant decision, using Michael Nygard's format (Title, Status, Context, Decision, Consequences). Numbered sequentially, never renumbered; a superseded decision gets a new ADR that supersedes the old one (old one's status updated to "Superseded by ADR-000N"), rather than editing history.

## Consequences

- Future contributors (human or AI) can answer "why is it built this way" without spelunking git history or asking someone who remembers.
- Adds a small amount of process: a genuinely structural decision should get an ADR before or alongside the code change, not as an afterthought.
- Decisions that are easily reversible with a single small edit should NOT get an ADR — this is a filter against process bloat, not a mandate to document everything.
