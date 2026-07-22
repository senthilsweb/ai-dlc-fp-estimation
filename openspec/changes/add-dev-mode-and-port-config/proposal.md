---
change-id: add-dev-mode-and-port-config
title: "Add --dev live-filesystem mode and a clear port-conflict error"
status: implemented
created: 2026-07-22
---

# Proposal: Add `--dev` Live-Filesystem Mode and a Clear Port-Conflict Error

## Summary

Two small operability gaps surfaced while actually running and testing the estimator: (1) every content edit — even a one-line dataset or HTML tweak — required a full `go build` because `app/` and `data/` are compiled into the binary, and (2) a port collision produced only the raw Go bind error with no hint of the fix. Both are addressed without touching the production (embedded) code path.

## What Changes

- **`--dev` / `FP_DEV=true`**: when set, `app/` and `data/` are served live from disk (`os.DirFS`) instead of the compiled-in `embed.FS`. Fails fast with a clear error if run somewhere without `app/index.html` on disk. Default remains `false` — no behavior change for existing deployments.
- **Clear port-conflict error**: if `r.Run(addr)` fails with "address already in use", the server now exits with a message naming the port and suggesting `--port`/`FP_PORT`, instead of a bare Go bind error.
- **`make dev`** and **`make run PORT=… APP=…`** Makefile targets, so both are reachable without remembering the flag names.

## Impact

- Affected specs: `dataset-partitioning` (ADDED — local dev mode, clear port-conflict error).
- Affected code: `main.go` only.
- Breaking: none — both changes are opt-in or purely additive to error messages.
