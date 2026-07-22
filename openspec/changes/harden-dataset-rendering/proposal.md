---
change-id: harden-dataset-rendering
title: "Harden rendering against malformed optional dataset fields; document the schema"
status: implemented
created: 2026-07-22
---

# Proposal: Harden Rendering Against Malformed Optional Fields; Document the Schema

## Summary

A direct question — "do we have a strongly typed schema for the JSON metadata, e.g. `keyJourneys` empty/null shouldn't cause a problem" — surfaced two real bugs, confirmed with a synthetic malformed dataset and a real browser drive (not just code review): (1) a single `keyJourneys` entry missing `steps` crashed the *entire* app, not just that card, because `renderProjectSummary()` was called unguarded from `initializeUI()`; (2) `fpConfig.gscDefinitions` — a field both real datasets populate — was silently ignored by a fully hardcoded JS array, discarding the Provly dataset's actual authored GSC default (`5`, silently replaced with a hardcoded `4`).

## What Changes

- `journey.steps` and `journey.name` default safely instead of throwing; an empty/malformed journey renders a placeholder, not a crash.
- `techStack`/`targetUsers`/`valuePropositions` check `.length > 0`, not just truthiness (an empty array is truthy in JS).
- `initializeUI()`'s per-section render calls are wrapped in a `safeRender(label, fn)` helper — one section's failure no longer takes down every other section.
- `gscDefinitions` now actually reads the dataset's `fpConfig.gscDefinitions` when present and non-empty, falling back to the standard IFPUG 14-factor list only when absent.
- Added `data/schema/metadata.schema.json` and `data/schema/product.schema.json` — documentation/IDE-validation aids (not runtime-enforced), referenced via `"$schema"` in each dataset's `metadata.json`.
- Added `data/qa-malformed-fixture/` as a permanent regression fixture: a synthetic dataset exercising missing/null/empty optional fields across the board.

## Impact

- Affected specs: `fp-estimation-engine` (ADDED — resilient per-section rendering, dataset-driven GSC definitions).
- Affected code: `app/index.html` only.
- Breaking: none. Both real datasets (`ai-agents-provly`, `tripma`) regression-tested via headless browser after the change — identical rendered output except the GSC default fix, which now correctly reflects each dataset's own authored value.

## Related

`docs/adr/0007-json-schema-for-dataset-validation.md` records the reasoning for keeping the schema documentation-only rather than runtime-enforced.
