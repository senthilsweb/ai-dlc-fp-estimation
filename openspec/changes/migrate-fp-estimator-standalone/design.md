## Context

The FP estimator existed as two forks in `ctms-gtm-mono-repo`, diverging only in six lines of `index.html` and their `data/*.json`. Growing this to a third project (Provly already proved a second fork was needed) meant either forking again or fixing the underlying design. Separately, a Node.js build step (`combine-wbs.js`) stood between editing data and seeing it rendered, which added friction and a second runtime dependency (Node) alongside the Go binary.

## Goals / Non-Goals

- Goals:
  - One app layer, N datasets, selected by config, no fork-per-project.
  - Single Go binary as the only runtime dependency — no Node build step.
  - Remove every place where the app assumed a specific project (title, export filename, legend labels, localStorage keys).
- Non-Goals:
  - Access control / multi-tenant role gating — explicitly deferred.
  - A database-backed rewrite (the source repo's `create-project-planner-app` proposal) — out of scope; this stays a static-SPA-plus-embedded-JSON design.
  - Building a UI dataset switcher — dataset selection is config/query-driven only for now.

## Decisions

- **Decision: merge JSON server-side in Go, not via a Node script.** The prior `combine-wbs.js` logic (read `metadata.json`, resolve `projectConfig`/`fpConfig`/`effortConfig` fallbacks, join product files by `dataFile`, tally status counts) was ported directly into `server/handler/data.go`, operating on `map[string]interface{}` since the schema is data-defined and varies per dataset. This removes Node as a dependency and means editing a dataset's JSON takes effect on the next request, not the next build.
  - Alternatives considered: keep `combine-wbs.js` and have the Go server just serve the pre-built `wbs-data.json` per app — rejected because it reintroduces a build step the "go-binary-driven" requirement was meant to eliminate, and because per-request merging costs negligible latency at this data size (confirmed by testing: ai-agents-provly at 87 capabilities and tripma at 160 capabilities both merge well under load time budgets).

- **Decision: dataset identity flows through the API response as `appId`, and the client derives its `localStorage` namespace from it.** Rather than hardcoding a prefix in `app/index.html`, `dataNamespace` is set from `wbsData.appId` (falling back to a slugified `brandPrefix`) immediately after the data fetch, before any `localStorage` read/write happens.
  - Alternatives considered: pass the namespace via a URL path segment (`/app/<name>/`) — rejected as unnecessary complexity for a single-SPA-at-root design; query-param-driven dataset switching plus response-derived namespacing achieves the same isolation with less routing surface.

- **Decision: status legend labels read from `statusDefinitions[key].label` at render time.** Previously two forks had two different hand-edited copies of the same five `<span>` tags. Now `applyStatusLabels()` runs once per data load and sets each legend span's text from the dataset, falling back to a generic English default if a dataset omits `statusDefinitions`.

- **Decision: `isProductVisible()` and the runtime-config plumbing (`window.__GTM_FP_CONFIG`, `/api/fp-config`, admin/hideCost/hideExports/hideRate) are deleted, not stubbed with dead branches.** No access control is in scope, so keeping the branching logic around unexercised would just be dead code inviting drift. `isProductVisible()` is kept as a single `return true` function — a no-op hook, not a feature — so future per-app product scoping (unrelated to access control) has an obvious seam without resurrecting the old role-gating design.

## Risks / Trade-offs

- **Schema drift between datasets**: since the merge is untyped (`map[string]interface{}`), a dataset with a malformed `metadata.json` fails at request time with a 404 + error message rather than at a build step. Mitigated by `buildAppData` returning a clear "unknown dataset" error including the underlying JSON error, and by the `add-fp-dataset` skill documenting the required shape.
- **`?app=` query override has no guardrail**: any caller can request any dataset folder name. Acceptable because there's no sensitive data distinction between datasets at this stage (explicitly no access control) — revisit if a dataset ever contains information not meant for all viewers.

## Migration Plan

1. Scaffold this repo; copy `ai-agents-provly/business-features-fp/data/*.json` → `data/ai-agents-provly/`, and `ctms-business-features-fp/data-tripma/*.json` → `data/tripma/` (renaming `metadata-tripma.json` → `metadata.json` for convention consistency).
2. Port `index.html`, removing the three hardcodes and all access-control plumbing.
3. Port `server/handler/spa.go` unchanged; add `server/handler/data.go` new.
4. Verify: `go build`, then curl `/api/data`, `/api/data?app=tripma`, `/api/apps`, and an unknown-dataset 404 case.
5. In the source repo: remove `ai-agents-provly/business-features-fp/` and the `add-capability-level-inclusion` openspec change (fully superseded here); leave `ai-agents-provly/inputs/` and `ai-agents-provly/pitch-deck/` in place (out of scope for this repo).
6. Rollback: the source repo's removal is a plain `git rm` + commit on `main`, recoverable via `git revert` if needed; nothing was force-pushed or rewritten.

## Open Questions

- Should a future `data/templrpress/` (or other) dataset get its own openspec change, or is adding a dataset covered entirely by the `add-fp-dataset` skill without a proposal? Current expectation: dataset additions are content changes (skip-proposal per `openspec/AGENTS.md`'s triage rules) unless they require an app/server code change.
