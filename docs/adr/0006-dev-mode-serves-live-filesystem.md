# 6. Dev mode serves app/ and data/ live from disk

Date: 2026-07-22

## Status

Accepted

## Context

`app/` and `data/` are compiled into the binary via `//go:embed` (ADR-0002, ADR-0003), which is the right choice for a distributable single-binary artifact. But it means that during development, every edit to `app/index.html` or a dataset's JSON — even a one-line label tweak — required a full `go build` before the change was visible, because the embedded copy is fixed at compile time. This was raised directly: "how to check it without go during dev time" — the friction wasn't building the binary once, it was rebuilding it for every content change, most of which (per ADR-0002) are supposed to be plain data/HTML edits with no Go involved at all.

## Decision

Add a `--dev` flag / `FP_DEV=true` env var that, when set, resolves `app/` and `data/` as `os.DirFS("app")` / `os.DirFS("data")` (read live from the current directory) instead of the compiled-in `embed.FS`. Everything else about the server is unchanged — same routes, same merge logic, same no-access-control posture. `--dev` fails fast with a clear error if `app/index.html` isn't found relative to the current directory, since the flag only makes sense run from the repo root.

This does not eliminate Go as a dependency — a binary still has to be built once (`go build` or `make dev`, which depends on `build`). What it eliminates is rebuilding *for every subsequent content edit*, which is the actual dev-loop cost.

## Consequences

- Editing `app/index.html` or any `data/<app>/*.json` while running with `--dev` is visible on the next browser refresh, no rebuild, no restart.
- A change to `main.go` or `server/handler/*.go` still requires a rebuild (and a restart) — `--dev` only affects how static/data assets are resolved, not the compiled server logic itself.
- `--dev` must not be used as a production flag: it reads from whatever `app/`/`data/` happen to be on disk relative to the current working directory, with no guarantee they match what was last committed or embedded. Default remains `false`; nothing changes for existing deployments.
- Distributing just the binary (e.g. a released artifact with no source tree alongside it) only works in the default embedded mode — `--dev` requires the source directories to be present.
