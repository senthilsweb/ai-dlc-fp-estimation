# 4. Defer access control

Date: 2026-07-22

## Status

Accepted

## Context

The source repo (`ctms-gtm-mono-repo`) had grown a `fp-runtime-config` feature on top of the FP estimator: gin-based admin/viewer auth, a `/api/fp-config` endpoint, and role-gated UI (hiding cost data, export buttons, and the rate input from non-admin viewers via a `window.__GTM_FP_CONFIG` object and a `.fp-cost-obfuscated` CSS class). That feature existed because the source repo's FP tool was one app among several in a shared, authenticated portal (`main.go`'s multi-app router with `GTM_ADMIN_USERNAME`/`GTM_USER_ALLOWED_APPS`, etc.) serving external stakeholders who shouldn't see cost data.

This repo was explicitly scoped to skip that: "No access control required at this time" was a direct instruction when this repo was extracted.

## Decision

Do not carry over any part of the auth/role-gating feature. Concretely: no login, no user/role model, no `/api/fp-config` endpoint, no `window.__GTM_FP_CONFIG`, no cost-obfuscation CSS, no per-viewer hiding of export/rate controls. `isProductVisible()` is kept in `app/index.html` as a bare `return true` — a no-op seam for a *possible* future per-app product filter, but explicitly not a security boundary and not wired to anything today. The `?app=` query override on `/api/data` is unguarded by design (see ADR-0002/0003): any caller can request any dataset folder name, because there is currently no concept of a viewer who shouldn't see a given dataset.

## Consequences

- Simpler code: no auth middleware, no session/cookie handling, no role branching in the SPA.
- Anyone who can reach the server can see every dataset in full (all cost/rate data, all products, all exports). This is fine as long as datasets contain nothing that needs to be kept from some viewers of the running instance.
- If a future requirement needs to hide some data from some viewers (e.g., sharing a running instance with an external stakeholder who shouldn't see cost data), that is a new decision requiring a new ADR — it should not be silently reintroduced by copying the source repo's `fp-runtime-config` design without re-evaluating whether that design still fits (e.g., the source design was role-based; a simpler per-link-token or per-deployment approach might now fit better given the single-app, config-driven shape of this repo).
