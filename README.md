# ai-dlc-fp-estimation

A standalone Function Point / WBS estimator, served as a single Go binary. One generic app layer, many datasets — each project's WBS lives in its own `data/<app-name>/` partition and the binary is told which one to serve.

## Quick Start

```bash
go build -o fp-estimator .
FP_APP=ai-agents-provly ./fp-estimator      # or FP_APP=tripma
open http://localhost:8080
```

Or with Docker:

```bash
cp .env.example .env
docker compose up --build
```

### Port already in use?

```bash
FP_PORT=8081 ./fp-estimator          # or: ./fp-estimator --port 8081
make run PORT=8081 APP=tripma        # Makefile shortcut, also lets you pick the dataset
```

A busy port fails fast with a message telling you to pick another, instead of a bare bind error.

### Iterating without rebuilding (`--dev`)

`app/` and `data/` are normally compiled into the binary (`//go:embed`), so by default even a one-line HTML or dataset tweak needs a `go build` to see it. Go is only required **once**, to produce the binary — but you can skip rebuilding on every content edit with `--dev` / `FP_DEV=true`, which serves `app/` and `data/` live from disk instead:

```bash
go build -o fp-estimator .           # once
FP_DEV=true ./fp-estimator           # or: make dev
# now edit app/index.html or any data/<app>/*.json and just refresh the browser —
# no rebuild needed until you change main.go or server/handler/*.go
```

`--dev` must be run from the repo root (it resolves `app/` and `data/` relative to the current directory). See [ADR-0006](docs/adr/0006-dev-mode-serves-live-filesystem.md) for why this isn't the default.

## Architecture

```
ai-dlc-fp-estimation/
├── main.go              # embeds app/ and data/, serves both via one binary
├── server/handler/
│   ├── spa.go            # generic embedded-FS SPA handler
│   └── data.go           # merges data/<app>/*.json at request time (replaces the old Node combine-wbs.js step)
├── app/
│   └── index.html        # the entire front end — one generic engine, no per-project branding baked in
└── data/
    ├── ai-agents-provly/  # Provly's WBS dataset
    │   ├── metadata.json  # project summary, FP weights, GSC defaults, status labels, product list
    │   ├── tech-stack.json
    │   └── p1-*.json ... p7-*.json
    └── tripma/             # demo dataset for testing the multi-app partitioning
        ├── metadata.json
        ├── tech-stack.json
        └── p1-*.json ... p4-*.json
```

The app layer reads everything it needs — title, branding, status labels, FP weights, GSC config, products — from whichever dataset the server hands it at `/api/data`. There is nothing in `app/index.html` that assumes a specific project.

## Adding a new dataset

1. Create `data/<your-app-name>/` with a `metadata.json` (see an existing dataset, and start with `"$schema": "../schema/metadata.schema.json"` for IDE validation/autocomplete against `data/schema/metadata.schema.json` — its field descriptions also document what happens if a field is omitted/null/empty), a `tech-stack.json`, and one JSON file per product listed in `metadata.json`'s `products` array (see `data/schema/product.schema.json`).
2. Run with `FP_APP=<your-app-name>` (or `--app <your-app-name>`).
3. No code changes, no rebuild — the Go server merges the files at request time.

See `.claude/skills/add-fp-dataset/SKILL.md` for the step-by-step checklist. There's no schema *enforcement* at runtime by design (see [ADR-0007](docs/adr/0007-json-schema-for-dataset-validation.md)) — the app is written to degrade gracefully when optional fields are missing. `data/qa-malformed-fixture/` is a permanent fixture for checking that after touching any renderer: `FP_APP=qa-malformed-fixture ./fp-estimator` should show zero browser console errors.

## Configuration

| Env var | Flag | Default | Purpose |
|---|---|---|---|
| `FP_PORT` | `--port` | `8080` | Listen port |
| `FP_APP` | `--app` | `ai-agents-provly` | Which `data/<name>/` partition to serve by default |
| `FP_DEV` | `--dev` | `false` | Serve `app/` and `data/` live from disk instead of the embedded copies — no rebuild needed to see edits |
| `FP_LOG_LEVEL` | `--log-level` | `info` | Log verbosity |
| `FP_LOG_FORMAT` | `--log-format` | `text` | `text` or `json` |

`GET /api/data?app=<name>` can override the active dataset per-request — there's no access control on this (none is required at this stage), it's a convenience for trying datasets side by side. `GET /api/apps` lists what's available under `data/`.

## Spec-driven development

This project follows the OpenSpec workflow — see `openspec/AGENTS.md` for the process and `openspec/project.md` for conventions. Active proposals live in `openspec/changes/`; built capabilities are recorded in `openspec/specs/`.

## Architecture Decision Records

The *why* behind this repo's structure is recorded in `docs/adr/` — see [`docs/adr/README.md`](docs/adr/README.md) for the index. Start with [ADR-0002](docs/adr/0002-separate-generic-app-layer-from-per-project-data.md) (app/data split) and [ADR-0004](docs/adr/0004-defer-access-control.md) (why there's no auth) if you're new here.

## Provenance

The app layer and estimation engine originated as `ctms-business-features-fp` in the `ctms-gtm-mono-repo`, and was forked once for Provly (`ai-agents-provly/business-features-fp`). This repo carries that engine forward as the generic, data-partitioned successor to both forks — the multi-tenant Go server, request-time JSON merge, and per-dataset `localStorage` namespacing didn't exist in either predecessor.
