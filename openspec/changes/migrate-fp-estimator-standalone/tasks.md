# Tasks: Extract the FP Estimator into a Standalone, Data-Partitioned App

## 1. Repo scaffold
- [x] 1.1 `git init` new repo at `ai-dlc-fp-estimation` (remote already existed: `github.com/senthilsweb/ai-dlc-fp-estimation`)
- [x] 1.2 Create `app/`, `data/`, `server/handler/`, `openspec/`, `.claude/skills/` directory structure

## 2. Data partitioning
- [x] 2.1 Copy `ai-agents-provly/business-features-fp/data/*.json` → `data/ai-agents-provly/`
- [x] 2.2 Copy `ctms-business-features-fp/data-tripma/*.json` → `data/tripma/`, renaming `metadata-tripma.json` → `metadata.json`
- [x] 2.3 Verify each dataset's `metadata.json` `products[].dataFile` entries match the copied filenames

## 3. App layer (generic engine)
- [x] 3.1 Copy base `index.html` into `app/`
- [x] 3.2 Remove hardcode: status legend labels — add `id="label-<status>"` spans, add `applyStatusLabels()` driven by `wbsData.statusDefinitions`
- [x] 3.3 Remove hardcode: export filename — derive from `wbsData.metadata.brandPrefix` instead of literal `'CTMS-WBS-'`/`'Provly-WBS-'`
- [x] 3.4 Remove hardcode: `localStorage` key prefix — add `dataNamespace`/`nsKey()`, set from `wbsData.appId` after data load, replace all `'ctms-*'` literal keys
- [x] 3.5 Change data fetch from static `wbs-data.json` to `GET /api/data`
- [x] 3.6 Remove all access-control plumbing: `window.__GTM_FP_CONFIG`, `/api/fp-config` fetch, `applyRuntimeConfig()`, `isAdminMode()`, cost-obfuscation CSS; simplify `isProductVisible()` to `return true`
- [x] 3.7 Remove all Zynomi/CTMS/Provly branding: title, footer copyright, fallback brand strings, stale GitHub doc link (updated to this repo)

## 4. Go server
- [x] 4.1 Port `server/handler/spa.go` unchanged (already generic)
- [x] 4.2 Write `server/handler/data.go` — Go-native port of `combine-wbs.js`'s merge logic, request-time, no Node dependency
- [x] 4.3 Write `main.go` — embeds `app/` and `data/`, `FP_PORT`/`FP_APP`/`FP_LOG_LEVEL`/`FP_LOG_FORMAT` config, `/api/data`, `/api/apps`, SPA fallback via `NoRoute` (no auth)
- [x] 4.4 `go.mod` module `github.com/senthilsweb/ai-dlc-fp-estimation`, no `zynomi`/`ctms` references anywhere in module path, binary name, or Docker image name

## 5. Verification
- [x] 5.1 `go build` succeeds
- [x] 5.2 `GET /api/data` (default app) returns merged JSON with correct `appId`, `summary` counts matching source data (87 capabilities for ai-agents-provly)
- [x] 5.3 `GET /api/data?app=tripma` returns the second dataset correctly (160 capabilities)
- [x] 5.4 `GET /api/data?app=<unknown>` returns a clean 404 with an error message, not a crash
- [x] 5.5 `GET /api/apps` lists both dataset folder names
- [x] 5.6 Served HTML confirmed to fetch `/api/data` and use `nsKey()` for every `localStorage` call (no bare `'ctms-*'` strings remain)
- [x] 5.7 Full-repo case-insensitive scan for `zynomi`/`ctms`/`provly` in `app/index.html` returns no matches

## 6. Supporting scaffolding
- [x] 6.1 `Makefile`, `Dockerfile`, `docker-compose.yml`, `.env.example`, `.gitignore`, `README.md` — all `FP_*` env vars, no `GTM_*`/`zynomi` naming
- [x] 6.2 `CLAUDE.md` — project non-negotiables (no project-specific branding in app/server, no access control, no Node build step, namespaced localStorage)
- [x] 6.3 `.claude/skills/add-fp-dataset/SKILL.md` — checklist for adding a new dataset
- [x] 6.4 `openspec/AGENTS.md` (copied, already generic), `openspec/project.md` (filled in for this repo)
- [x] 6.5 This change's spec deltas: `specs/fp-estimation-engine/spec.md`, `specs/dataset-partitioning/spec.md`

## 7. Source repo cleanup
- [x] 7.1 Remove `ai-agents-provly/business-features-fp/` from `ctms-gtm-mono-repo` (fully migrated)
- [x] 7.2 Remove `openspec/changes/add-capability-level-inclusion/` from `ctms-gtm-mono-repo` (superseded by `specs/fp-estimation-engine/spec.md` here)
- [x] 7.3 Leave `ai-agents-provly/inputs/` and `ai-agents-provly/pitch-deck/` in place in the source repo (out of scope — not FP-estimation content)
- [x] 7.4 Commit the removal in the source repo with a message explaining the migration and pointing at this repo

## 8. Architecture Decision Records
- [x] 8.1 Confirm no ADRs exist anywhere in `ctms-gtm-mono-repo` to copy/adapt (repo-wide search came up empty — only an unrelated deploy guide in `docs/`)
- [x] 8.2 Write `docs/adr/0001-record-architecture-decisions.md` (meta-ADR establishing the practice)
- [x] 8.3 Write `docs/adr/0002-separate-generic-app-layer-from-per-project-data.md`
- [x] 8.4 Write `docs/adr/0003-merge-datasets-in-go-at-request-time.md`
- [x] 8.5 Write `docs/adr/0004-defer-access-control.md`
- [x] 8.6 Write `docs/adr/0005-namespace-client-state-by-active-dataset.md`
- [x] 8.7 Cross-link ADRs from `README.md`, `CLAUDE.md`, and this change's `design.md`
