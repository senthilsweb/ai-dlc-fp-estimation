# 5. Namespace client-side state by the active dataset

Date: 2026-07-22

## Status

Accepted

## Context

The SPA persists per-user scenario-planning state (which products/features/capabilities are included, status filter selections, GSC slider values) to `localStorage` so it survives page reloads. Both predecessor forks hardcoded the key prefix as the literal string `'ctms-'` — including the Provly fork, which forked the entire file but never renamed this prefix, because each fork ran on its own origin/port so the collision was invisible. This repo can serve multiple datasets from a single running instance (via `FP_APP`/`?app=`, see ADR-0002), so two datasets sharing a browser profile and origin would silently share — and corrupt — each other's saved inclusion/filter/GSC state under a fixed prefix.

## Decision

Derive the `localStorage` key prefix (`dataNamespace`, exposed via the `nsKey(name)` helper in `app/index.html`) from the loaded dataset's `appId` field (returned by `GET /api/data`), falling back to a slugified `metadata.brandPrefix` if `appId` is absent. This is set once, immediately after the data fetch resolves and before any `localStorage` read happens, so every subsequent read/write in the same page load is correctly namespaced.

## Consequences

- Two datasets (e.g., `ai-agents-provly` and `tripma`) can be tried in the same browser, on the same origin, without their saved state colliding — verified during migration by toggling capability exclusions under one dataset and confirming the other was unaffected.
- The namespace is derived, not configured — a dataset doesn't need to declare its own storage prefix, removing a class of copy-paste mistake (forgetting to change it, as both prior forks did).
- If a dataset's `appId`/folder name is ever renamed, users revisiting it will see fresh (default-included) state rather than their prior saved state — an acceptable trade-off since dataset renames are expected to be rare and the state is scenario-planning convenience, not data of record.
