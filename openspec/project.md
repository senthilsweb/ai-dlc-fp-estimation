# Project Context

## Purpose

A standalone, single-Go-binary Function Point (IFPUG) / WBS estimator. One generic app layer serves any number of project datasets, each partitioned under `data/<app-name>/` and selected at runtime via config — no rebuild, no per-project fork of the app code.

This project is the generic successor to two prior forks that lived in `ctms-gtm-mono-repo`: `ctms-business-features-fp` (CTMS) and `ai-agents-provly/business-features-fp` (Provly). Those forks duplicated the entire app layer per project and diverged only in `data/`; this repo collapses that back into one engine plus N datasets.

## Tech Stack
- Go 1.22, Gin (routing), logrus (logging) — single static binary, no external runtime dependencies
- Front end: single-file HTML/CSS/vanilla JS SPA (Tailwind + Iconify via CDN, no build step)
- Data: per-app JSON files under `data/<app-name>/`, merged at request time by the Go server (no Node build step)

## Project Conventions

### Code Style
- Go: standard `gofmt` formatting, no external linters configured yet.
- Front end: no bundler — `app/index.html` is a single file; keep additions dependency-free (CDN scripts only) unless a design.md justifies otherwise.

### Architecture Patterns
- **App/data separation**: `app/` never contains project-specific content. Anything that varies per project (branding, products, FP weights, GSC defaults, status labels) lives in `data/<app-name>/metadata.json` and related files, and is read at runtime via `GET /api/data`.
- **Config-driven dataset selection**: the server picks which dataset to serve via `FP_APP` env var / `--app` flag (default), with an optional `?app=` query override for convenience. No access control gates this — it is out of scope at this stage.
- **Request-time merge, not a build step**: `server/handler/data.go` merges `metadata.json` + `tech-stack.json` + per-product files into the shape the SPA expects, replacing the old `combine-wbs.js` Node script from the predecessor forks.
- **Namespaced client state**: `localStorage` keys are prefixed with the active dataset's `appId` (see `nsKey()` in `app/index.html`) so multiple datasets can be exercised in one browser without clobbering each other's saved state.

### Testing Strategy
- No automated test suite yet. Verify changes by running `go build` and curling `/api/data`, `/api/data?app=<other>`, and `/api/apps` against at least two datasets, plus a manual browser check of the WBS tree, FP totals, and exports.

### Git Workflow
- `main` is the working branch. Conventional, verb-led commit messages (`feat:`, `fix:`, `refactor:`).

## Domain Context

IFPUG Function Point Analysis: products → features → capabilities, each capability typed and rated for complexity, contributing weighted function points; General System Characteristics (GSC) adjust the total via a Value Adjustment Factor. This domain vocabulary is fixed across datasets — what varies is which products/features/capabilities exist for a given project.

## Important Constraints
- No authentication/authorization in this repo at this stage — this was a deliberate scope cut when extracting this repo from a larger multi-tenant portal. Revisit only if explicitly requested.
- Keep the app layer free of any single project's name, logo, or copy — that is the entire point of this repo's separation from its predecessor forks.

## External Dependencies
- None beyond the Go module dependencies in `go.mod` (Gin, logrus) and CDN-hosted front-end libraries (Tailwind, Iconify, SheetJS/XLSX) loaded directly in `app/index.html`.
