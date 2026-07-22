<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

## Project

A single Go binary that serves a generic Function Point / WBS estimator SPA (`app/index.html`) plus one or more project datasets (`data/<app-name>/*.json`). The active dataset is chosen at runtime via config (`FP_APP` env var / `--app` flag), not baked into the binary or the HTML.

## Non-negotiables

- **No project-specific branding or logic in `app/index.html` or `server/`.** Everything project-specific (title, brand prefix, products, FP weights, GSC defaults, status labels) lives in `data/<app-name>/metadata.json` and is read at runtime via `/api/data`. If you catch yourself hardcoding a project name anywhere in the app or server layer, it belongs in a dataset instead.
- **No access control in this repo at this stage.** Don't add auth, roles, or per-viewer gating unless explicitly asked — this was a deliberate simplification from the multi-tenant portal this project was extracted from.
- **No Node/build step.** The old `combine-wbs.js` script that merged `data/*.json` into `wbs-data.json` was ported into `server/handler/data.go` — the Go binary does the merge at request time. Don't reintroduce a JS build step for this.
- **localStorage keys must stay namespaced.** The app derives a per-dataset key prefix from `wbsData.appId` (see `nsKey()` in `app/index.html`) so two datasets loaded in the same browser don't clobber each other's saved inclusion/GSC state. Never hardcode a bare localStorage key.

## Adding a dataset

See `.claude/skills/add-fp-dataset/SKILL.md` and the README's "Adding a new dataset" section.

## Spec-driven workflow

Follow `openspec/AGENTS.md` for any new capability or architecture change: scaffold a proposal under `openspec/changes/<change-id>/` (proposal.md, tasks.md, spec deltas), get it reviewed, implement, then archive. `openspec/specs/` holds the specs for what's actually built.

## Architecture Decision Records

The four non-negotiables above are each backed by an ADR in `docs/adr/` — read those before reversing one of them, and add a new ADR (see `docs/adr/0001-record-architecture-decisions.md` for when one is warranted) rather than silently reintroducing something a prior ADR deliberately ruled out (e.g. don't resurrect role-gating without a new ADR revisiting ADR-0004).
