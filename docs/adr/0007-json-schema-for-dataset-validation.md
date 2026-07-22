# 7. JSON Schema for dataset documentation, defensive JS as the runtime safety net

Date: 2026-07-22

## Status

Accepted

## Context

A direct question surfaced a real gap: "do we have a strongly typed schema for the JSON metadata — e.g. if `keyJourneys` is empty or null, it should not cause a problem." Investigating found two concrete issues, not just the hypothetical one:

1. **A real crash bug.** `projectSummary.keyJourneys` itself was already guarded (`if (summary.keyJourneys && summary.keyJourneys.length > 0)`), but a single journey entry missing its `steps` array was not — `journey.steps.map(...)` threw a `TypeError`, and because `renderProjectSummary()` was called directly (unguarded) from `initializeUI()`, the exception aborted every subsequent render call (stats, WBS tree, tabs — everything), not just the one journey card. Reproduced with a synthetic dataset (`data/qa-malformed-fixture/`) before fixing: the WBS tree showed the generic "Failed to load data" error, and the whole app was unusable.
2. **A real silently-ignored field.** `fpConfig.gscDefinitions` is documented and populated in both real datasets, but `app/index.html` had a fully hardcoded 14-entry GSC array that never read `wbsData.gscDefinitions` at all — confirmed by the Provly dataset's authored default of `5` for "Data Communications" being silently discarded in favor of a different hardcoded `4`.

Neither Go structs nor a JSON-schema-validation library exists in this codebase (`server/handler/data.go` merges via untyped `map[string]interface{}`, per ADR-0003's deliberate choice to keep dataset shape flexible without a build/codegen step). So the honest answer to "do we have a strongly typed schema" was **no** — and the two bugs above are exactly the kind of thing that absence allows to go unnoticed.

## Decision

Two separate things, deliberately not one:

1. **Fix the actual defensive-coding gaps** (the runtime safety net, which is what actually prevents a bad dataset from breaking the app):
   - `journey.steps` defaults to `[]` and renders a "No steps defined for this journey" placeholder instead of throwing; `journey.name` falls back to "Untitled journey".
   - `techStack`, `targetUsers`, `valuePropositions` now check `.length > 0`, not just truthiness, so an empty array (which is truthy in JS) doesn't render an empty section header.
   - `initializeUI()`'s render calls are now each wrapped via a `safeRender(label, fn)` helper — one section's failure logs a console error naming which dataset field to check, and every *other* section still renders. A future oversight like the `keyJourneys` bug will degrade one card, not the whole app.
   - `gscDefinitions` now actually reads `wbsData.gscDefinitions` (dataset-provided) when present and non-empty, falling back to the standard IFPUG 14-factor list (renamed `defaultGscDefinitions`) only when the dataset omits it.

2. **Add JSON Schema files** (`data/schema/metadata.schema.json`, `data/schema/product.schema.json`) as **documentation and IDE-validation aids**, referenced via a `"$schema"` key in each dataset's `metadata.json`. This is deliberately *not* wired into the Go server or a build step — no new Go dependency, no new required tooling, consistent with ADR-0003's "no build step" stance and `openspec/AGENTS.md`'s "avoid frameworks without clear justification." Editors that support `$schema` (VS Code and others, natively) give real-time validation/autocomplete while authoring a dataset, without the runtime paying any cost for it.

## Consequences

- The two real bugs found are fixed and regression-tested (both real datasets rebuilt and browser-driven; `qa-malformed-fixture` kept permanently under `data/` as a regression fixture — every one of its deliberately-malformed fields is documented inline and cross-referenced from this ADR).
- "Strongly typed" here means *documented and IDE-checkable*, not compiler-enforced — a dataset author can still ignore schema warnings and ship a file that violates it. The runtime safety net (item 1 above) is what actually guarantees a violation degrades gracefully instead of crashing; the schema's job is to make violations rare and visible at authoring time, not to be the only defense.
- Every future optional field added to the dataset shape should follow the same pattern demonstrated here: guard against absent/null/empty at the render site, and document the fallback behavior in the schema's `description`, not just in a comment.
- If a stronger guarantee is ever needed (e.g. CI validating datasets before merge), the schema files are already in the right shape to hand to any standard JSON Schema validator (`ajv-cli` via `npx`, a Go schema library, etc.) — that would be a new decision (a new ADR), not a retroactive change to this one, since it would add a dependency this decision deliberately avoided.
