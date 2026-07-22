# The AI-DLC Estimation Model

> **Status:** Foundational intent — recorded 2026-07-22, before implementation.
> This document captures the owner's reasoning for what this tool is becoming and why.
> Numbers here are **starting points to be calibrated**, not findings. See [Calibration](#8-calibration-how-this-stops-being-guesswork).

## 0. Positioning and vocabulary

This tool estimates software delivery along the **delivery continuum**:

| Stage | Who does the work | Modelled as |
|---|---|---|
| **Human-led** | People write the software | Nominal PDR, no factors — conventional benchmarks apply |
| **AI-Assisted** | Human-led build, agent assists (completion, drafting, review support); the human owns the work | A productivity factor |
| **AI-Native** (agent-led) | Agent produces the bulk; human steers, reviews, decides, approves | A productivity factor |

The stages are **values in the tool, not labels on it** — each is a productivity factor composing onto the nominal human baseline, so an estimate can sit anywhere along the continuum, including partway between stages. That is the direct benefit of composition (§3): the continuum is expressible rather than being three hard-coded modes.

**AI-DLC** is the reference *methodology* for the AI-Native end, and supplies the hierarchy vocabulary (Intent → Bolt → Unit of Work).

**On terminology choices.** These follow how the market already frames it, so the tool is legible to people outside this repo:

- GlobalLogic's **VelocityAI** describes guiding enterprises along *"a delivery continuum, from human-led to increasingly AI-native software development"* — hence *continuum* and *AI-native* rather than invented alternatives.
- Ascendion positions as *"AI-Native Software Engineering Services"*, with **AAVA** (Agentic Value Accelerator) as the agentic platform and *EngineeringAI* as the method.

Deliberately avoided: **"AI Engineering"**, which in industry means engineering *AI systems* (building ML/LLM products) — the inverse of engineering *with* AI, and a guaranteed source of confusion.

## 1. The gap this fills

There is no scientific estimation mechanism for AI-DLC — or for AI-assisted delivery generally — where **the agent does ~90% of the work and the human is the steerer, reviewer, decision maker, and approver of anything requiring approval.**

Traditional estimation assumes a human writes the code. Every calibrated model in the industry (COCOMO, ISBSG benchmarks, standard PDR tables) is built on that assumption. When an agent writes the code and a human steers, those models don't just need a discount — they need re-derivation. Nobody has published that work.

This tool is an attempt at a defensible starting point.

## 2. Core thesis: FP is the invariant, everything else is a knob

**Function Points are independent of who or what builds the software, and independent of the methodology used to build it.** A function's functional size does not change because an agent wrote it instead of a human. That independence is exactly why FP is the right measurement to keep.

What changes under AI-DLC is not the size — it's the **rate**.

> Under the hood it is FP, with knobs and whistles to tune.

So the design rule is:

| Layer | Varies by methodology? | Where it lives |
|---|---|---|
| Functional size (FP) | **No** — invariant | Capability/task `type` + `complexity` |
| Rate (hours per FP, "PDR") | **Yes** — this is the dial | `projectConfig.defaultPDR` |
| Productivity factors | Yes | `effortConfig.productivityFactors` |
| Vocabulary (Intent/Bolt/Task) | Yes — presentation only | `projectConfig.levels` |
| Status lifecycle | Yes — presentation only | `statusDefinitions` |

Everything below the first row is configurable. The first row is the anchor that makes cross-methodology comparison meaningful at all.

## 3. The productivity dial (PDR)

PDR — hours per Function Point — is the primary knob:

| Delivery mode | Illustrative PDR |
|---|---|
| AI-DLC, simple task within a bolt | **1 FP ≈ 15 minutes** (0.25 h) |
| Human-driven, e.g. Java enterprise development | **1 FP ≈ 8 hours** |

That is roughly a **32× spread**, and it is the honest centre of the value story: same functional size, radically different delivery rate.

These two figures are the owner's working starting points. They are **not** derived from measured data yet — see [Calibration](#8-calibration-how-this-stops-being-guesswork).

### How the rate is composed (implemented 2026-07-22)

**Terminology.** *PDR* is **Project Delivery Rate** — hours to deliver one Function Point, the standard bridge from size to effort (`effort = size × PDR`). ISBSG publishes PDR benchmarks by language and platform, which is what makes an FP estimate externally defensible. **Nominal PDR** is the rate *before* any modifiers; **effective PDR** is what remains after them. The nominal/adjusted split mirrors COCOMO's nominal effort plus effort multipliers.

PDR is **not** set directly to an AI-adjusted number. It composes, so each discount stays separately visible and separately calibratable:

```
effective PDR = nominal PDR × Π(active productivity factors)
hours         = AFP × effective PDR
cost          = hours × rate
```

Nominal PDR is therefore the **human baseline before AI**, and the AI speedup lives *only* in the productivity factors. That gives the discount one source of truth.

**The double-counting trap this avoids.** Nominal PDR and the productivity factors both express "how fast do we deliver":

| Setup | Nominal PDR | Factors | Effective | |
|---|---|---|---|---|
| A | 0.25 | none | 0.25 | correct |
| B | 8 | AI 0.65 × Reuse 0.7 | 3.64 | correct |
| C | **0.25** | **AI 0.65** | **0.1625** | **AI counted twice** |

Case C is easy to reach because both inputs are individually true statements ("we deliver ~15 min/FP", "we use AI agents") that encode the *same* speedup. The convention — nominal is pre-AI — prevents it.

**Technology stack is documentation only.** It is recorded for optics and carried into generated proposals/SOWs, and deliberately does **not** drive hours or cost.

The Effort & Cost tab renders the full chain so any number is auditable, e.g.:

```
8 (nominal) × 0.65 (AI Agent-Assisted) × 0.7 (Component Reuse) = 3.64 h/FP
```

**Naming:** productivity factors must be vendor-neutral — "AI Agent-Assisted", "AI Agent-Led" — never product names.

Why composition rather than a single blended PDR: when actuals arrive, a blended figure can't tell you *which* term was wrong. Composition can. That matters directly for §8.

**All knobs are live inputs**, seeded from the dataset, persisted per dataset in `localStorage`, with a reset-to-dataset-defaults control: base PDR, hours/day, days/month, tech stack (single-select — one dominant stack), productivity factors (multi-select, they compound), rate/hr, and the VAF constants. Per point 3 of the owner's direction, granularity is the estimator's choice: one config applies to every bolt in a dataset, and per-bolt assumptions are achieved by splitting into separate datasets and switching with `?app=` — no per-bolt override machinery.

**What was broken before this (2026-07-22):** three separate defects meant almost none of this worked.
- `projectConfig.defaultPDR` was read by the app but never passed through by `server/handler/data.go`, so any dataset value was silently discarded and every estimate used `8`.
- `productivityFactors` and `techStackFactors` were rendered with bold `0.65×`-style badges but **never multiplied into anything** — `renderProductivityFactors()` and `renderTechFactors()` only built HTML. They were pictures of dials.
- The VAF constants, hours/day, days/month and a literal `160` "working hours/month" were hardcoded in the markup and math.

Of the four things that should tune an AI-DLC estimate — PDR, productivity factors, tech factors, GSC — only GSC actually worked. This is the same defined-but-ignored bug class as `gscDefinitions` (ADR-0007) and `levels` (still outstanding).

## 4. Bolt complexity is rated, never inferred from task count

**Task count is not a proxy for complexity.** Both of these are true and common:

- A bolt with **many tasks** that are each simple → completed quickly.
- A bolt with **few tasks** that are each complex → slow and risky.

So complexity is an **assigned rating** on the work itself, not a number derived by counting children. Any future feature that tries to auto-derive a bolt's size from how many tasks it contains is working against this principle.

(The reference project bears this out: its Bolt 3 had *more* tasks than its Bolt 1 but roughly a quarter of the verification burden. Counting tasks would have ranked those two backwards.)

## 5. Evals are units of work, not a separate effort model

Within each bolt, **eval design, dataset preparation, and rubric preparation are tasks/activities** — performed by the agent with human support, exactly like any other task in the bolt.

This is deliberate and it settles a question that was open during design: eval effort is **not** a separate additive line item, and **not** a phase-weight redistribution. Evals are simply units of work that get sized like every other unit of work. No parallel estimation path.

This also means eval-heavy work is naturally priced: if a bolt needs 37 evals, those evals appear as tasks, carry their own size, and roll up. Nothing special is required.

## 6. A bolt is the work packet, and it is 100% deployable when done

The bolt is the unit that gets estimated and the unit that gets delivered. When a bolt is done it is deployable — it carries its own full slice of the lifecycle (elaboration → construction → evals → deploy) rather than handing off to a downstream test or release phase.

## 7. AI-DLC vocabulary is a presentation layer

The UI should speak AI-DLC — **Intent → Bolt → Unit of Work** — because that is what makes the tool legible to people practising the method:

> Only then will the people going to use this tool, or adapt it, understand.

But the machinery underneath stays Function Points. The vocabulary swap is presentation; the estimation core is unchanged. This is a relabeling plus configurable weights, not a new estimation paradigm.

The hierarchy mapping:

| Generic (IFPUG framing) | AI-DLC framing |
|---|---|
| Product | **Intent** |
| Feature | **Bolt** |
| Capability | **Unit of Work / Task** |

**Implementation note:** `projectConfig.levels` already exists in every dataset (`{"1":"Product","2":"Feature","3":"Capability"}`) and is already passed through by the Go merge layer — but `app/index.html` never reads it. Wiring it up is most of what this relabeling requires.

## 8. Calibration: how this stops being guesswork

Every number in this document is a hypothesis today. The plan to fix that:

1. Capture **pre-estimates** for small projects and individual bolts.
2. Record **post actuals** for the same work.
3. Feed the delta back into PDR and the productivity factors.

Over time this produces the thing the industry currently lacks: a **measured** hours-per-FP figure for agent-led delivery, rather than an asserted one. Until then, the tool's outputs should be read as structured, tunable projections — not predictions.

## 9. Intended workflow

This tool is aimed at **greenfield projects and new proposals**, where the process is:

1. Requirements understanding
2. Q&A to clarify
3. Work breakdown — expressed AI-DLC style (Intent → Bolt → Unit of Work), because the engineering and the process are both AI-DLC
4. Feed bolts, tasks, and evals into this tool
5. Get tangible **effort hours and dollar value** out

The mob (PM / Architect / the collaborating team) identifies the bolts, tasks, and evals during elaboration. This tool turns that breakdown into hours and cost.

## 10. On the `agent-job-matcher` reference

`agent-job-matcher` is a **reference only** — it is the one project attempted end-to-end following full AI-DLC, so it is a useful source of realistic structure and a sanity check on the model. It is explicitly **not** the authoritative reference, and the model must not be overfitted to it.

Its role: the first worked example used to decorate a dataset JSON and validate that the UI and the resulting hours/dollars are sensible.

For the record, its measured structure (parsed 2026-07-22): **7 intents / 28 bolts / 190 tasks** (177 complete). Two of the seven intents have no bolt layer at all — tasks hang directly off the intent — so a strict three-level tree needs a synthetic wrapper bolt for those. It contains **no estimation data of any kind**: no story points, no hours, no complexity ratings. Effort figures for it would have to be assigned, not read.

## 11. Status

**This is a starting point.** It is expected to change as calibration data arrives. Nothing here should be treated as settled beyond the core thesis in §2 — that FP is the invariant and everything else is a configurable weight.

---

## Related

- [ADR-0008](adr/0008-fp-as-invariant-ai-dlc-as-presentation.md) — the architectural decision recording §2 and §7
- [ADR-0007](adr/0007-json-schema-for-dataset-validation.md) — why dataset fields are documented but not runtime-enforced
- [ADR-0002](adr/0002-separate-generic-app-layer-from-per-project-data.md) — the app/data split that makes methodology-as-configuration possible
