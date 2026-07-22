---
name: add-fp-dataset
description: Add a new project dataset (data/<app-name>/) to the FP estimator so it can be served by setting FP_APP, without touching app/ or server/ code. Use whenever the user wants to estimate a new project, add a WBS for another product, or set up a demo dataset for this repo.
---

# Add a new FP estimator dataset

This repo serves one generic app layer (`app/index.html`) against whichever dataset the server is configured to load (`FP_APP` env var / `--app` flag, or `?app=` query override). Adding a new project means adding a data partition — never editing `app/` or `server/`.

## Steps

1. **Create the folder**: `data/<app-name>/` (kebab-case, becomes the value of `FP_APP` and the `appId` in `/api/data`).

2. **Write `metadata.json`**. Model it on an existing dataset (`data/ai-agents-provly/metadata.json` or `data/tripma/metadata.json`). Required shape:
   - `projectConfig`: `title`, `mainProductName`, `brandPrefix`, `organization`, `defaultHourlyRate`, `currency`, `hoursPerDay`, `daysPerMonth`, `levels`, and `products` — an array of `{ shortCode, description, dataFile }` entries, one per product file you'll create in step 4.
   - `fpConfig`: `fpWeights`, `complexityGuidelines`, `gscDefinitions`, `gscRatingGuide`.
   - `effortConfig`: `techStackFactors`, `productivityFactors`, `sdlcPhases`.
   - `statusDefinitions`: label/description per status key (`completed`, `partial`, `in-progress`, `roadmap`, `beta`, ...) — this is what drives the status-legend text in the UI, so relabel here rather than touching HTML.
   - `projectSummary` (optional but recommended): powers the Project Summary tab — `name`, `tagline`, `description`, `techStack`, `businessContext`, `keyJourneys`.
   - `glossary` (optional).

3. **Write `tech-stack.json`** (optional): `{ "products": [...], "summary": {...} }` for the Tech Stack tab.

4. **Write one JSON file per product**, named exactly as referenced in `metadata.json`'s `products[].dataFile`. Each file is `{ name, features: [{ name, capabilities: [{ name, type, complexity, status, ... }] }] }`.

5. **Run it**: `FP_APP=<app-name> ./fp-estimator` (or `go run . --app <app-name>`), then hit `http://localhost:8080`. To sanity-check the merge without a browser: `curl "http://localhost:8080/api/data?app=<app-name>"`.

6. **Verify `/api/apps`** lists your new folder — it's derived from the directory listing, so no registration step is needed beyond creating the folder.

## Common mistakes

- Forgetting a `dataFile` referenced in `metadata.json` — the server silently skips missing product files (same behavior as the old `combine-wbs.js`), so a product silently vanishing from the UI usually means a typo'd filename.
- Reusing another dataset's `appId`/folder name — this becomes the `localStorage` namespace, so a collision would mix two datasets' saved inclusion state in the same browser profile.
- Adding project-specific copy to `app/index.html` instead of `projectSummary`/`statusDefinitions` in `metadata.json`. If it's about one project, it belongs in the dataset.
