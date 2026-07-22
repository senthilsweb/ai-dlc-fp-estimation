# 3. Merge dataset JSON in Go at request time, not via a Node build step

Date: 2026-07-22

## Status

Accepted

## Context

Both predecessor forks required running `node combine-wbs.js` to pre-merge `data/*.json` (metadata, tech-stack, per-product files) into a single `wbs-data.json` before the SPA could load it — a manual step easy to forget after editing data, and a second language/runtime (Node) alongside the Go server that eventually served the files. This repo's explicit goal is to be "Go-binary driven" — the binary should be the only thing you need to run it.

## Decision

Port `combine-wbs.js`'s merge logic (resolve `projectConfig`/`fpConfig`/`effortConfig` fallbacks in `metadata.json`, join each product's JSON file by its `dataFile` reference, tally status counts, fold in `tech-stack.json`) into `server/handler/data.go`, operating on `map[string]interface{}` since the schema is entirely data-defined. The Go server performs this merge on every `GET /api/data` request and returns the combined JSON directly — there is no `wbs-data.json` file and no Node dependency anywhere in this repo.

## Consequences

- Editing a dataset's JSON takes effect on the next request (or the next binary rebuild, if the data is embedded via `//go:embed` for a production build) — no separate combine step to remember.
- Node.js is not a dependency of this repo at all, at build time or runtime.
- The merge is untyped in Go (`map[string]interface{}`) rather than backed by structs, trading compile-time schema safety for the flexibility of letting each dataset define its own shape without a Go code change. A malformed dataset fails at request time with a descriptive error (see ADR-0002's consequences) rather than at a build step.
- Per-request merge cost was measured negligible at current data sizes (87 capabilities for the `ai-agents-provly` dataset, 160 for `tripma`); this decision should be revisited if a dataset grows large enough that merge latency becomes noticeable, e.g. by caching the merged result and invalidating on file change.
