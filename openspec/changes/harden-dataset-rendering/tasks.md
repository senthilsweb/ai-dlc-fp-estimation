# Tasks: Harden Rendering Against Malformed Optional Fields

- [x] 1.1 Create `data/qa-malformed-fixture/` (missing/null/empty `keyJourneys[].steps`, null `gscDefinitions`, null `sdlcPhases`) and confirm it reproduces a real crash before fixing anything
- [x] 1.2 Confirm via headless browser: `TypeError: Cannot read properties of undefined (reading 'map')` at `journey.steps.map`, aborting `initializeUI()` — WBS tree shows the generic "Failed to load data" error, nothing renders
- [x] 1.3 Fix `journey.steps` (default `[]`, placeholder text) and `journey.name` (fallback label)
- [x] 1.4 Add `.length > 0` checks to `techStack`/`targetUsers`/`valuePropositions` ternaries
- [x] 1.5 Add `safeRender(label, fn)` helper; wrap every per-section render call in `initializeUI()`
- [x] 1.6 Fix `gscDefinitions` to read `wbsData.gscDefinitions` (rename hardcoded array to `defaultGscDefinitions`, used only as fallback)
- [x] 1.7 Rebuild, re-run the malformed fixture through the browser: confirm zero console errors, WBS tree renders, each malformed journey shows "No steps defined for this journey"
- [x] 1.8 Regression-test both real datasets via browser: confirm identical output, plus confirm the GSC fix (Provly's first GSC slider now reads `5`, its actual authored default, not the previously-hardcoded `4`)
- [x] 1.9 Write `data/schema/metadata.schema.json` and `data/schema/product.schema.json`; add `"$schema"` to all three datasets' `metadata.json`
- [x] 1.10 Rebuild and re-verify all three datasets still merge/render correctly after the `$schema` key addition
- [x] 1.11 Write `docs/adr/0007-json-schema-for-dataset-validation.md`, add to `docs/adr/README.md` index
- [x] 1.12 Update `openspec/specs/fp-estimation-engine/spec.md` with the new requirements
- [x] 1.13 Update `.claude/skills/add-fp-dataset/SKILL.md` and `openspec/project.md` to reference the schema files and the QA fixture
