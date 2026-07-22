---
change-id: migrate-fp-estimator-standalone
title: "Extract the FP estimator into a standalone, data-partitioned Go app"
status: implemented
created: 2026-07-22
---

# Proposal: Extract the FP Estimator into a Standalone, Data-Partitioned App

## Summary

Extract the Function Point / WBS estimator out of `ctms-gtm-mono-repo` (where it existed as two diverging forks — `ctms-business-features-fp` and `ai-agents-provly/business-features-fp`) into this standalone repo, restructured so the app layer is fully generic and every project's data lives in its own `data/<app-name>/` partition selected at runtime by the Go server. No project name, branding, or copy is hardcoded in the app or server layer.

## Motivation

The two prior forks were byte-for-byte identical except for six lines (title, logo, two status-legend labels, footer, export filename) and their `data/*.json` — proving the app layer was already conceptually generic, but the repo structure didn't reflect that. Each new project meant a full fork-and-hand-edit instead of dropping in a new dataset. Three concrete hardcodes blocked treating the engine as truly generic:

1. Status legend labels ("Roadmap"/"Beta" vs "Build New"/"Integrate OSS") were literal HTML strings, never read from `metadata.json`'s `statusDefinitions.*.label`, even though that field existed and was populated per fork.
2. The Excel export filename prefix (`'CTMS-WBS-'` / `'Provly-WBS-'`) was a literal string, not derived from `metadata.brandPrefix` even though that value was already fetched and used for the page title two lines above.
3. `localStorage` keys were hardcoded with a literal `ctms-` prefix in both forks (including the Provly fork, which never renamed them) — harmless only because each fork ran on its own origin/port; unsafe once multiple datasets can be tried against the same running instance.

Separately, the prior architecture required a Node.js build step (`combine-wbs.js`) to pre-merge `data/*.json` into `wbs-data.json` before serving — an extra step this migration removes by moving that merge into the Go server itself, so the binary alone is sufficient to run any dataset.

## What Changes

- **New repo** (`ai-dlc-fp-estimation`), not a directory inside `ctms-gtm-mono-repo`.
- **App/data split**: `app/index.html` (the generic engine, stripped of all CTMS/Provly/Zynomi branding) + `data/<app-name>/` (one partition per project).
- **Two datasets seeded**: `data/ai-agents-provly/` (moved from `ai-agents-provly/business-features-fp/data/`) and `data/tripma/` (copied from `ctms-business-features-fp/data-tripma/`, as a second dataset to prove the partitioning works end-to-end).
- **Go-native data merge**: `server/handler/data.go` replaces `combine-wbs.js`, merging `metadata.json` + `tech-stack.json` + per-product files into the response shape the SPA expects, at request time, served at `GET /api/data`. `GET /api/apps` lists available dataset folders.
- **Config-driven dataset selection**: `FP_APP` env var / `--app` flag chooses the default dataset the binary serves; `?app=` query param overrides per-request (no access control gates this — none is in scope, see below).
- **All three hardcodes removed**: status legend labels now read from `wbsData.statusDefinitions[key].label`; export filename derives from `wbsData.metadata.brandPrefix`; `localStorage` keys are prefixed via `nsKey()`, namespaced on the active dataset's `appId`.
- **No access control**: the role-gated runtime-config feature from the prior repo (`fp-runtime-config`: admin mode, cost obfuscation, export/rate hiding) is intentionally **not** carried over — out of scope at this stage, per explicit instruction.
- **Spec-driven scaffolding**: `openspec/`, `CLAUDE.md`, and a `.claude/skills/add-fp-dataset` skill are initialized in this repo so future dataset additions and app-layer changes follow the same proposal → spec → implement workflow as the source repo.

## Impact

- Affected specs (this change): `fp-estimation-engine` (ADDED — generic WBS/FP engine, including capability-level inclusion carried forward from the source repo's `add-capability-level-inclusion` change), `dataset-partitioning` (ADDED — data/<app> partitioning, config-driven selection, request-time merge).
- Affected code: entire repo (new).
- Source repo impact: `ai-agents-provly/business-features-fp/` and the `add-capability-level-inclusion` openspec change are removed from `ctms-gtm-mono-repo` now that this repo is their successor. `ai-agents-provly/inputs/` and `ai-agents-provly/pitch-deck/` remain in the source repo — out of scope for an FP-estimation-focused repo.
- Breaking: N/A (new repo).

## Out of Scope

- Access control / roles / admin gating (explicitly deferred).
- The `templrpress`-style additional datasets mentioned as a future example — only `ai-agents-provly` and `tripma` are seeded now.
- The alternate Next.js/Postgres "project planner" architecture proposed in the source repo's `create-project-planner-app` — never built there, and not adopted here; this migration keeps the static-SPA-plus-Go-binary shape.
