# Architecture Decision Records

This directory records the architecturally significant decisions for this repo, using [Michael Nygard's ADR format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions): Title, Status, Context, Decision, Consequences.

## Index

| # | Title | Status |
|---|---|---|
| [0001](0001-record-architecture-decisions.md) | Record architecture decisions | Accepted |
| [0002](0002-separate-generic-app-layer-from-per-project-data.md) | Separate the generic app layer from per-project data partitions | Accepted |
| [0003](0003-merge-datasets-in-go-at-request-time.md) | Merge dataset JSON in Go at request time, not via a Node build step | Accepted |
| [0004](0004-defer-access-control.md) | Defer access control | Accepted |
| [0005](0005-namespace-client-state-by-active-dataset.md) | Namespace client-side state by the active dataset | Accepted |
| [0006](0006-dev-mode-serves-live-filesystem.md) | Dev mode serves app/ and data/ live from disk | Accepted |
| [0007](0007-json-schema-for-dataset-validation.md) | JSON Schema for dataset docs, defensive JS as the runtime safety net | Accepted |
| [0008](0008-fp-as-invariant-ai-dlc-as-presentation.md) | FP stays the invariant; AI-DLC is presentation plus configurable weights | Accepted |

For the estimation model itself — the reasoning, the PDR dial, and the calibration plan — see [`../ai-dlc-estimation-model.md`](../ai-dlc-estimation-model.md).

## When to add a new ADR

Add one when a decision would be expensive to reverse, affects more than one file, or a future contributor (human or AI) would otherwise have to reverse-engineer *why* from the code. Skip it for anything reversible with a single small edit — that belongs in a commit message, not here.

## Relationship to `openspec/`

This repo also follows OpenSpec (`openspec/AGENTS.md`) for proposing and tracking behavior changes. The two serve different purposes:
- **ADRs** (here) record *why* a structural/technical decision was made — permanent, rarely revisited once accepted.
- **OpenSpec `design.md`** files record the technical approach *for one specific change* — scoped to that change's proposal, folded into the ADR record only if the decision turns out to be durable.

The ADRs below were extracted from `openspec/changes/migrate-fp-estimator-standalone/design.md`, which covers the same decisions in the context of that specific migration.
