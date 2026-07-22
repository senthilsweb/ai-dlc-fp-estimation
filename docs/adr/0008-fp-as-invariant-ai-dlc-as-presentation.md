# 8. Function Points stay the invariant; AI-DLC is presentation plus configurable weights

Date: 2026-07-22

## Status

Accepted

## Context

This repo is being adapted for AI-DLC — a delivery method where an agent performs roughly 90% of the work and the human acts as steerer, reviewer, decision maker, and approver. The question was how deep that adaptation should go, and specifically whether IFPUG Function Points survive it.

Two positions were considered seriously:

1. **FP is a category error for AI-DLC.** IFPUG classifies work as EI/EO/EQ/ILF/EIF — a user-facing functional view. AI-DLC's Intent → Bolt → Unit of Work is a delivery/process breakdown. Bolts like "Package skeleton + deterministic core" or "Telemetry backends" don't naturally fit any of the five FP types, so forcing a type per bolt produces numbers that look rigorous but encode a guess.
2. **FP is the only methodology-independent measure available**, and its independence is precisely the point.

Position 2 won, on the owner's reasoning: **a function's functional size does not change because an agent wrote it instead of a human.** That invariance is what makes it possible to compare agent-led delivery against human-led delivery at all. If the size measure itself moved with the methodology, there would be no baseline and no value story — only two incomparable numbers.

There is also a genuine gap being filled: no calibrated estimation model exists for agent-led delivery. Every industry benchmark (COCOMO, ISBSG, standard PDR tables) assumes a human writes the code. Those models need re-derivation, not a discount — and re-derivation needs a stable unit to re-derive *against*.

## Decision

**Keep Function Points as the invariant measurement. Express everything that varies by methodology as configurable, data-driven weights. Make AI-DLC a presentation layer over that core.**

Concretely:

- **Invariant (never methodology-dependent):** functional size, i.e. each unit of work's FP `type` and `complexity`.
- **The primary knob:** PDR — hours per Function Point — carried in `projectConfig.defaultPDR`. Working starting points: **1 FP ≈ 15 minutes** for a simple task under AI-DLC, versus **1 FP ≈ 8 hours** for human-driven Java enterprise work. Roughly a 32× spread, and that spread *is* the value story.
- **Presentation only:** the level vocabulary (Intent / Bolt / Unit of Work) via `projectConfig.levels`, and the status lifecycle via `statusDefinitions`. The UI must speak AI-DLC so practitioners recognise their own process — but nothing underneath changes.
- **Complexity is assigned, never inferred from child counts.** Many simple tasks can be faster than few complex ones; deriving a bolt's size by counting its tasks would rank work backwards.
- **Evals are units of work.** Eval design, dataset prep, and rubric prep are tasks inside a bolt, sized like any other task — not an additive effort line and not a phase-weight redistribution. This closes a design question that was open: no parallel estimation path for verification effort.

Rejected alternatives: eval-weighted sizing as the primary unit (no industry calibration exists, and the reference project only had eval counts for one of seven intents); and a full AI-DLC pivot abandoning FP (destroys the cross-methodology baseline that motivates the tool).

## Consequences

- The existing estimation engine survives intact. This adaptation is relabeling plus wiring up configuration, not a rewrite — consistent with ADR-0002's app-generic/data-driven split.
- Datasets of different methodologies remain comparable, because they share a size unit. A traditional and an AI-DLC estimate of the same scope differ in PDR, not in FP.
- **Three fields are defined-but-ignored and must be wired up for this to work at all**, all the same bug class:
  - `projectConfig.defaultPDR` — `app/index.html` reads it and defaults to `8`, but `server/handler/data.go` never passes it through, so a dataset setting it is silently discarded. This is the single most important knob in the model and it is currently non-functional.
  - `projectConfig.levels` — present in every dataset and passed through by Go, but never read by the app. This is most of what the AI-DLC relabeling needs.
  - `gscDefinitions` — same class, already fixed under ADR-0007.
- **Every effort number the tool currently produces is a hypothesis, not a prediction.** The PDR figures above are asserted, not measured. The mitigation is a deliberate calibration loop: capture pre-estimates and post-actuals for small projects and bolts, feed the delta back into PDR and the productivity factors. Until enough of those accumulate, outputs should be presented as tunable projections.
- Assigning FP `type` to process-shaped bolts remains genuinely awkward (the concern behind position 1 above). This decision accepts that friction as the price of a methodology-independent baseline, rather than resolving it. If bolt-level FP typing proves unworkable in practice, that warrants a new ADR revisiting this — not a silent drift to some other unit.

## Related

- [`docs/ai-dlc-estimation-model.md`](../ai-dlc-estimation-model.md) — the full model, workflow, and calibration plan
- [ADR-0002](0002-separate-generic-app-layer-from-per-project-data.md) — the app/data split this decision depends on
- [ADR-0007](0007-json-schema-for-dataset-validation.md) — the earlier defined-but-ignored field, and the defensive-rendering stance
