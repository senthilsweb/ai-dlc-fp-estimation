# 2. Separate the generic app layer from per-project data partitions

Date: 2026-07-22

## Status

Accepted

## Context

The FP estimator existed as two forks in `ctms-gtm-mono-repo` — `ctms-business-features-fp` (CTMS) and `ai-agents-provly/business-features-fp` (Provly). Diffing the two showed they were identical except for six lines (page title, header logo, two status-legend labels, footer text, export filename) and their `data/*.json` content. Every new project meant copying the whole app directory and hand-editing those few lines — a fork-per-project model that doesn't scale past two or three projects and drifts further apart with every fork (the two forks had already independently diverged on those six lines).

## Decision

Split the tool into one generic `app/` (the SPA — HTML/CSS/JS, containing zero project-specific branding, copy, or logic) and N data partitions under `data/<app-name>/` (one directory per project: `metadata.json`, `tech-stack.json`, per-product JSON files). The Go server is told which partition to serve via config (`FP_APP` env var / `--app` flag), with an optional `?app=` query override. Adding a new project means adding a new `data/<app-name>/` directory — no code change, no fork.

Everything that previously lived only in the six hardcoded lines now comes from the dataset: page title and export filename from `metadata.brandPrefix`, status legend labels from `metadata.statusDefinitions[key].label`.

## Consequences

- New projects are a data-only change (see `.claude/skills/add-fp-dataset/SKILL.md`), reviewable and testable without touching `app/` or `server/`.
- The app layer must never assume a specific project again — this is enforced by convention and code review, not by a type system, so a future contributor could still accidentally hardcode something into `app/index.html`. `CLAUDE.md` calls this out explicitly as a non-negotiable.
- A dataset with a malformed `metadata.json` (missing a required field, wrong `dataFile` reference) fails at request time rather than at a build/compile step — traded off deliberately against removing the Node build step (see ADR-0003).
- This does not (yet) provide per-project UI customization beyond what `metadata.json` exposes (branding, labels, GSC/FP weights, products). A project needing a genuinely different UI, not just different data, would need a new ADR revisiting this decision.
