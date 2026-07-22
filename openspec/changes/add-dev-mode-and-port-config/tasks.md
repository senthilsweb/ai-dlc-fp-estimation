# Tasks: Add --dev Live-Filesystem Mode and a Clear Port-Conflict Error

- [x] 1.1 Add `--dev`/`FP_DEV` flag to `main.go`; when true, resolve `app`/`data` via `os.DirFS` instead of the embedded FS
- [x] 1.2 Fail fast with a clear error if `--dev` is used without `app/index.html` on disk (wrong working directory)
- [x] 1.3 Verify: run with `--dev`, edit `app/index.html` on disk, confirm the running (already-started) binary serves the edited content on the next request with no rebuild
- [x] 1.4 Wrap the port-bind error: detect "address already in use" and print a message naming the port and the `--port`/`FP_PORT` fix
- [x] 1.5 Verify: start two instances on the same port, confirm the second prints the friendly message
- [x] 1.6 Add `make dev` and `PORT`/`APP` overrides to `make run`
- [x] 1.7 Document in `README.md` (Quick Start: port-conflict + `--dev` subsections, Configuration table), `CLAUDE.md` (Commands section), `.env.example` (`FP_DEV`)
- [x] 1.8 Write `docs/adr/0006-dev-mode-serves-live-filesystem.md`, add to `docs/adr/README.md` index
- [x] 1.9 Update `openspec/specs/dataset-partitioning/spec.md` with the two new requirements
