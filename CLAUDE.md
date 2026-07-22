# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

## Commands

```bash
go build -o fp-estimator .        # or: make build
FP_APP=ai-agents-provly ./fp-estimator   # or: FP_APP=tripma ./fp-estimator; or go run . --app <name>
make build-linux / make build-darwin     # cross-compile
make docker                              # build the Docker image
docker compose up --build                # run via compose (reads .env — cp .env.example .env first)

# Manual smoke test (there is no automated test suite — see openspec/project.md's Testing Strategy):
curl "http://localhost:8080/api/data"              # active dataset
curl "http://localhost:8080/api/data?app=tripma"   # override for a specific dataset
curl "http://localhost:8080/api/apps"              # list dataset folders under data/
```

There is no linter or formatter configured beyond `gofmt` conventions; there is no JS build step by design (see ADR-0003) — `app/index.html` is served as-is.

## Architecture

Request flow: `main.go` embeds `app/` and `data/` via `//go:embed`, reads `FP_APP`/`FP_PORT`/`--app`/`--port` at startup, and fails fast if the configured dataset directory doesn't exist under `data/`. It registers exactly three routes: `GET /api/data` and `GET /api/apps` (`server/handler/data.go`), and everything else falls through `r.NoRoute` to `server/handler/spa.go`'s embedded-FS SPA handler (serves the requested static file, or `index.html` as the SPA fallback).

`server/handler/data.go` is the one piece of real logic on the server side: `buildAppData()` reads `data/<appId>/metadata.json`, resolves its `projectConfig`/`fpConfig`/`effortConfig` sub-blocks (falling back to top-level fields for older flat-schema datasets), joins in each product's JSON file via the `dataFile` reference in `metadata.json`'s product list, optionally folds in `tech-stack.json`, tallies capability status counts, and returns one merged JSON object — the same shape the predecessor forks used to produce via a `combine-wbs.js` Node script, but computed fresh on every request instead of as a build artifact.

`app/index.html` is a single-file SPA (vanilla JS, Tailwind/Iconify/SheetJS via CDN, no bundler). On load it fetches `/api/data`, sets `dataNamespace` from the response's `appId` (see `nsKey()`), and renders everything — title, WBS tree, FP/GSC calculations, status legend, exports — from that one payload. It has three layers of inclusion state (product → feature → capability), each cascading from its parent, persisted to `localStorage` under the dataset-namespaced keys.

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
