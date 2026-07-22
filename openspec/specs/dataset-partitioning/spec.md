# Dataset Partitioning

## Requirements

### Requirement: Datasets are partitioned by app name under data/
Every project's WBS data SHALL live under its own `data/<app-name>/` directory containing `metadata.json` (required), `tech-stack.json` (optional), and one JSON file per product referenced from `metadata.json`'s product list. Adding a new project SHALL require only adding a new `data/<app-name>/` directory — no changes to `app/` or `server/`.

#### Scenario: New dataset requires no code change
- **GIVEN** a new `data/acme-widgets/` directory is created with a valid `metadata.json` and product files
- **WHEN** the server is started with `FP_APP=acme-widgets`
- **THEN** the app renders `acme-widgets`' data correctly
- **AND** no file under `app/` or `server/` was modified to make this work

### Requirement: Active dataset is selected via server configuration
The server SHALL determine the default active dataset from the `FP_APP` environment variable or `--app` flag at startup, and SHALL fail to start if the configured dataset does not exist under `data/`.

#### Scenario: Valid configured dataset
- **GIVEN** `FP_APP=ai-agents-provly` and `data/ai-agents-provly/metadata.json` exists
- **WHEN** the server starts
- **THEN** it starts successfully and serves that dataset by default at `GET /api/data`

#### Scenario: Invalid configured dataset
- **GIVEN** `FP_APP=does-not-exist`
- **WHEN** the server starts
- **THEN** it exits with a clear error naming the missing dataset, rather than starting in a broken state

### Requirement: Per-request dataset override
`GET /api/data` SHALL accept an optional `app` query parameter that overrides the server's configured default dataset for that request only, without requiring a restart. This capability is not access-controlled.

#### Scenario: Query override serves a different dataset
- **GIVEN** the server's configured default is `ai-agents-provly`
- **WHEN** a client requests `GET /api/data?app=tripma`
- **THEN** the response contains `tripma`'s data, with `appId: "tripma"`

#### Scenario: Unknown dataset in override
- **GIVEN** a client requests `GET /api/data?app=does-not-exist`
- **WHEN** the server processes the request
- **THEN** it responds with HTTP 404 and an error message naming the missing dataset, without crashing or affecting the server's default-dataset configuration

### Requirement: Request-time JSON merge, no build step
The server SHALL merge a dataset's `metadata.json`, `tech-stack.json`, and per-product JSON files into the API response shape at request time. No Node.js or other external build step SHALL be required between editing a dataset's JSON and seeing the change reflected.

#### Scenario: Editing a dataset takes effect without a build step
- **GIVEN** the server is running with a dataset loaded from disk (not embedded)
- **WHEN** a product JSON file is edited and the server is restarted (embedded builds require a rebuild; a `go run` dev loop does not)
- **THEN** the next `GET /api/data` reflects the edited content, with no intermediate combine/build command required

### Requirement: Dataset discovery
`GET /api/apps` SHALL list the names of all dataset directories present under `data/`, derived from the directory listing rather than a maintained registry.

#### Scenario: New dataset appears automatically
- **GIVEN** a new `data/acme-widgets/` directory is added
- **WHEN** a client calls `GET /api/apps`
- **THEN** `"acme-widgets"` appears in the returned list, with no separate registration step

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
