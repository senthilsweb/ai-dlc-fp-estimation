# Dataset Partitioning

## ADDED Requirements

### Requirement: Local development mode without rebuilding
The server SHALL support a `--dev`/`FP_DEV=true` mode in which `app/` and `data/` are served directly from the local filesystem instead of the compiled-in embedded copies, so that editing a dataset's JSON or `app/index.html` is visible on the next request without a rebuild. This mode SHALL be opt-in (default `false`) and SHALL fail with a clear error if `app/index.html` is not found relative to the current directory.

#### Scenario: Editing a dataset in dev mode requires no rebuild
- **GIVEN** the server is running with `--dev` from the repo root
- **WHEN** a file under `data/<active-app>/` is edited on disk
- **THEN** the next `GET /api/data` reflects the edit, with no `go build` or restart required

#### Scenario: Dev mode run from the wrong directory
- **GIVEN** `--dev` is passed but the current directory has no `app/index.html`
- **WHEN** the server starts
- **THEN** it exits immediately with an error naming the missing path, rather than starting and 404ing on every request

### Requirement: Clear error on port conflict
If the configured port is already in use, the server SHALL exit with a message that names the port and suggests the `--port`/`FP_PORT` override, rather than surfacing only the underlying bind error.

#### Scenario: Starting a second instance on an occupied port
- **GIVEN** an instance is already listening on port 8080
- **WHEN** a second instance is started with the same port
- **THEN** it exits with a message naming port 8080 and suggesting `--port`/`FP_PORT` to pick another
