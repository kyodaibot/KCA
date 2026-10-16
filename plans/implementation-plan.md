# KCA Implementation Plan

Build the project one verifiable increment at a time. This roadmap defines sequence and completion; [requirements](../docs/requirements.md) define product scope, [technical design](../docs/technical-design.md) defines shared implementation rules, and the [source-access record](../docs/discovery/source-access.md) holds acquisition evidence and restrictions.

## Chosen architecture

Go + Gin in `backend/` serves the API and business logic using pgx + sqlc. Python jobs in `rag/` handle offline data preparation and retrieval evaluation; React/TypeScript/Vite stays in `frontend/`. PostgreSQL + pgvector is shared. The authoritative boundaries, migration owner and technology rationale are in [design §1](../docs/technical-design.md#1-architecture-and-ownership).

Keep the existing root layout and Compose entry point. Introduce the local Go/Docker environment early; add academic schema, generated queries and cross-language compatibility only when their stages need them. A permanent Python service is unnecessary.

## Scope and milestones

The MVP includes accounts, onboarding and editable profiles/credit summaries, strict course search, reviewed curriculum predicates, semantic ranking, multilingual chat, comparisons, citations, feedback and personal-data deletion. Full scope and exclusions remain in [requirements §§5–7](../docs/requirements.md#5-functional-requirements).

Proposed pilot: Engineering / 地球工学科 / AY2026, one term initially, 30–100 permitted offerings and one applicable official curriculum document. Confirm this scope before collecting real data; a historical export must be labeled with its actual year. Affiliation/profile storage supports more scopes than the acquired course corpus.

- **First learning milestone:** Tasks 1–5, each completed separately; no course corpus or model provider required.
- **First data milestone:** a 3–5-page permitted sample passes Tasks 7–11, producing one source-backed internal Go result or an honest empty result.
- **First personalized milestone:** Task 17 serves authenticated, profile-aware structured recommendations using real session/onboarding guards.
- **Complete MVP:** all tasks and [design §9 release gates](../docs/technical-design.md#9-verification-and-release-gates), including real-source quality evaluation. Synthetic fixtures prove engineering behavior only.

## Shared development workflow

1. Begin implementation with a fresh Git status/diff and file inspection. Preserve unrelated work; never assume the scaffold, database contents or working tree still match earlier observations.
2. Take one task at a time. At its start, choose focused files/interfaces, verify relevant current documentation and pin only dependencies needed for that task. Write meaningful failing checks, implement the smallest useful result, then run the task's checks and relevant regressions. Task completion requires observed results, not expected results written in a plan.
3. Use Go unit/httptest and PostgreSQL integration checks for API/business logic, pytest for Python artifacts/jobs, and TypeScript/component/E2E checks for UI. Database tests use disposable named resources, never the normal service database; provider/mail tests use fakes without bypassing authentication.
4. Keep test setup proportional to the task. Add race, cross-runtime, E2E, load and restore checks when their owning stages arrive. Authentication/profile fixtures must precede personalized HTTP checks; browser bootstrap/CSRF/proxy fixtures must precede authenticated E2E.
5. Record deliverables, actual verification and remaining blockers before advancing. Review changes in small increments; commits and deployment follow the owner's separate instructions. No unavailable execution skill is required.

## Development sequence

Tasks 1–5 are deliberately small. Task 6 is an independent research track. Later tasks may use clearly labeled synthetic fixtures while source access is blocked; those fixtures cannot satisfy the real-data milestone or quality gates. Dependencies below describe engineering prerequisites; Task 6 additionally gates every use of real source material.

### Task 1: Run a local Go health endpoint

**Depends:** none.
- **Deliver:** minimal Go module and Gin process in `backend/`, listening on port 8000 with `GET /health` returning `{"status":"ok"}`; graceful shutdown and only the configuration needed to run it.
- **Verify:** run locally and check HTTP status/body with a focused handler test. No database, authentication, sqlc or Python setup is needed; health remains live when later dependencies are unavailable.

### Task 2: Build and run the minimal Go image

**Depends:** 1.
- **Deliver:** small multi-stage backend Docker build and non-root runtime containing the Go binary, without Python dependencies or application secrets.
- **Verify:** build/run it locally and exercise the same health contract through the published port. Inspect runtime user, shutdown and image exclusions; preserve unrelated scaffold files pending later cleanup.

### Task 3: Start PostgreSQL + pgvector through Compose

**Depends:** 2.
- **Deliver:** minimal local Compose wiring for the Go backend image and PostgreSQL 16/pgvector, retaining the existing database volume and localhost development ports.
- **Verify:** validate Compose configuration without exposing secrets; start a separate disposable project, check database health and extension availability, and confirm data survives restart. Do not reset an existing volume to make startup succeed.

### Task 4: Connect with pgx and add readiness

**Depends:** 3.
- **Deliver:** pgx pool, environment-based database configuration, bounded connection checks and `GET /ready`.
- **Verify:** readiness is 200 with a reachable database and 503 with invalid/unavailable database connectivity; health still returns 200. Close connections on shutdown and keep credentials out of responses/logs.

### Task 5: Introduce one versioned SQL migration

**Depends:** 4.
- **Deliver:** explicit migration command, first small migration establishing the schema-contract registry, and readiness validation for the schema version now required by the application. Follow [migration ownership](../docs/technical-design.md#1-architecture-and-ownership).
- **Verify:** apply to a fresh disposable database, repeat without duplicate changes, and exercise an upgrade/version mismatch. Runtime startup does not migrate automatically; keep academic tables and sqlc out of this increment.

### Task 6: Establish a permitted source sample

**Depends:** none; proceed independently of Tasks 1–5.
- **Deliver:** update the [source-access record](../docs/discovery/source-access.md) with allowed acquisition/use, select a 3–5-page sample or export, label published/missing fields and locate an applicable official curriculum document and period/year evidence.
- **Verify:** every real artifact has original URL, year evidence, acquisition time, encoding/hash and explicit storage/display/embedding permissions. Human comparison preserves Japanese text/layout. Missing authorization or rule applicability is recorded as blocked/unverified; do not automate access to the disallowed host.

### Task 7: Introduce academic schema and artifact contracts

**Depends:** 5; 6 informs real-source fields but is not needed for synthetic schema work.
- **Deliver:** focused migrations for catalog/provenance/meetings/assessment and reviewed rule structures in [design §2](../docs/technical-design.md#2-data-model-and-artifact-contracts), plus versioned normalized artifact schemas. Establish runtime roles and compatibility checks as tables are introduced.
- **Verify:** constraints reject invalid identities, broken evidence links and incompatible contract versions; unknown fields remain nullable. Test course/offering distinction, multiple meetings, immutable revisions and restricted RAG grants without adding identity tables prematurely.

### Task 8: Add sqlc queries and a Go catalog reader

**Depends:** 7.
- **Deliver:** sqlc configuration reading canonical migrations, generated pgx query code and a small Go reader for synthetic catalog/evidence fixtures.
- **Verify:** generation is reproducible; queries round-trip UUIDs, UTF-8, UTC timestamps and nullable fields. Keep generated rows separate from API models; use parameterized queries and no second schema authority.

### Task 9: Parse syllabus artifacts and produce Markdown

**Depends:** 7 contracts; 6 for real artifacts.
- **Deliver:** Python package/CLI under `rag/`, manifest validation, HTML parsing, source-needed PDF adapter, normalized records, assessment extraction candidates and Markdown with evidence locators.
- **Verify:** hand-labeled CP932/HTML and representative PDF tables retain text, all meetings and page references. Corrupt, undecodable, image-only or unsupported layouts require review. Apply the [candidate/accepted assessment lifecycle](../docs/technical-design.md#3-assessment-evidence-lifecycle); missing credits/year or mixed grading must not become confirmed facts.

### Task 10: Publish revisions and synchronize incrementally

**Depends:** 7 and 9; 6 for real publication.
- **Deliver:** transactional Python publication, review acceptance, stable identity mapping, audit/coverage reporting and operator-triggered incremental refresh.
- **Verify:** repeat import is idempotent, partial sync never withdraws records, failed publication leaves the prior revision intact, and a changed/conflicting assessment invalidates the current accepted answer. Exercise missing-year quarantine, duplicate titles and the [withdrawal policy](../docs/technical-design.md#2-data-model-and-artifact-contracts).

### Task 11: Read Python data in Go and enforce hard filters

**Depends:** 8 and 10.
- **Deliver:** Python-write/Go-read compatibility fixtures and internal structured search over accepted current academic revisions. No anonymous personalized route is needed for this increment.
- **Verify:** source/evidence/decimal/null values round-trip; day and period belong to the same meeting; multi-meeting fit, period ranges, unknown restrictions, irregular courses, year/term scope and strict no-exam constraints produce exact expected sets. Empty results never relax filters. Follow [design §4](../docs/technical-design.md#4-search-rules-and-response-behavior).

### Task 12: Evaluate reviewed curriculum predicates in Go

**Depends:** 7 and 11; 6 plus applicable reviewed documents for real eligibility claims.
- **Deliver:** operator import of reviewed rules, version/cohort/track resolution, deterministic predicate evaluation and cited category explanations.
- **Verify:** effective-date boundaries, exclusions, exact course mappings, missing history/track and overlapping versions produce the defined eligible/ineligible/unverified outcomes. Aggregate self-reported credits cannot prove completion. Python/LLM retrieval does not evaluate these predicates.

### Task 13: Add accounts and revocable sessions

**Depends:** 5; may proceed independently of course acquisition.
- **Deliver:** account migration, registration/login/logout/session APIs, password hashing, session/CSRF/Origin controls, database-backed limits and genuine isolated two-user test clients using fake mail.
- **Verify:** [design §5 security defaults](../docs/technical-design.md#5-authentication-and-personal-data-security), generic failures, cookie rotation, logout replay, expiry, restart persistence, multi-worker rate limits and secret-free telemetry. Never disable guards to create fixtures.

### Task 14: Complete email verification and recovery

**Depends:** 13.
- **Deliver:** expiring single-use verification/reset actions, recent-password reauthentication, fake local mail flow and configurable real mail adapter.
- **Verify:** used/expired tokens fail, reset revokes all sessions, reauthentication rotates session/CSRF, and unverified accounts cannot reach personalized access. Local tests send no real mail; public deployment cannot bypass mail/HTTPS prerequisites.

### Task 15: Persist profiles and independent credit summaries

**Depends:** 7 and 13.
- **Deliver:** owner-scoped profile/taxonomy/options APIs, required-field onboarding state, exact summary decimals, optional versioned category reports, optimistic revision checks and genuinely onboarded two-user API fixtures.
- **Verify:** [design §6](../docs/technical-design.md#6-profile-and-credit-semantics), including save/reload/edit of every summary with no curriculum rules, null/zero/omission distinctions, category scope, atomic invalid-patch rollback, ownership tampering and concurrent edits. Login/restart retains separate profiles; optional fields may all be skipped.

### Task 16: Interpret intent with minimized provider input

**Depends:** 11 and 15; fake provider is sufficient initially.
- **Deliver:** typed multilingual intent extraction, alias resolution, provider adapter/deadlines and server-side projection/redaction per [design §8](../docs/technical-design.md#8-provider-privacy-and-operations).
- **Verify:** Japanese/Chinese/English ambiguity and follow-up cases, invalid enums/extra keys, malformed model output, timeout and prompt injection. Inspect structured and free-text payloads; unsafe minimization yields no provider call and credits are omitted unless explicitly needed.

### Task 17: Serve profile-aware structured recommendations

**Depends:** 11–16.
- **Deliver:** authenticated search/chat endpoints, fresh profile context, confirmed temporary overrides, deterministic reason/citation cards and source-linked comparisons. Use the [request contract](../docs/technical-design.md#7-api-and-temporary-profile-overrides).
- **Verify:** real session/onboarding/CSRF/revision guards on both endpoints; no repeated questions for stored fields; unconfirmed/model-generated overrides cannot take effect or write profiles. Test account switches, stale context, explicit clears, comparison IDs, coverage warnings, unknown evidence IDs and provider failure. This milestone can return structured recommendations before semantic ranking is added.

### Task 18: Produce chunks and corpus embeddings in Python

**Depends:** 10; selected provider/model and processing rights for real embeddings.
- **Deliver:** versioned chunk/vector/model schema additions, section/evidence-bound chunks, shared model manifest, corpus embedding job and changed-chunk rebuild support.
- **Verify:** Japanese boundaries/overlap, source revision ownership, model fingerprints, dimension/preprocessing and partial-coverage reporting. Corpus vectors never contain profile/credit data; jobs terminate and do not introduce a Python API.

### Task 19: Add Go query embeddings and vector ranking

**Depends:** 8, 11, 12 and 18; integrates with 17 when available.
- **Deliver:** query embedding adapter matching the corpus model, pgvector driver/sqlc compatibility check, filtered exact ranking and scoped rule-passage retrieval.
- **Verify:** Python-written vectors are read by Go with documented numerical tolerance; wrong model/revision/dimension fails clearly; ranked IDs remain a subset of hard candidates, repeated chunks do not duplicate offerings, stale vectors/rule passages are excluded, and empty topics skip embedding calls. Provider outage cannot substitute an old model or bypass guards.

### Task 20: Establish authenticated browser infrastructure

**Depends:** 17 and 14.
- **Deliver:** React/Vite base, generated API types, same-origin proxy, cookie/CSRF-aware client, session/profile bootstrap, onboarding gate and isolated API-onboarded browser fixtures.
- **Verify:** personalized calls wait for bootstrap/completion, rotation refreshes CSRF, 401 clears personal state, 403/409 are explicit and two browsers never share context. These checks precede authenticated feature E2E; no bearer tokens in browser storage.

### Task 21: Build the Japanese-first discovery interface

**Depends:** 20; 19 adds semantic behavior without blocking structured UI development.
- **Deliver:** chat, filter summary, course/evidence cards, comparison and hypothetical-value review; optional Chinese/English responses, accessible keyboard interactions, loading/error/empty/clarification states.
- **Verify:** authenticated clarification → cards → multilingual follow-up → comparison journeys; preserve unrelated filters, render excerpts as text, show all meetings/unknowns/snapshot coverage, submit exact reviewed overrides and label hypothetical results without changing saved profiles.

### Task 22: Add registration, onboarding and profile-management forms

**Depends:** 15, 20 and 21.
- **Deliver:** register/verify/login/recovery screens, short required onboarding, optional progress Skip, profile editor and account settings using the existing client/gates.
- **Verify:** form-driven register → verify → login → onboard → recommend → hypothetical exploration → explicit saved edit → refreshed recommendation. Exercise no-rule summaries, invalid dependent affiliation choices, revision conflicts retaining unsaved edits, reload, expiry and account switching.

### Task 23: Implement personal-data deletion

**Depends:** 14, 15 and 22.
- **Deliver:** profile/account deletion, recent-auth UI, transactional scoped deletion markers and operator cleanup/replay support under [design §5](../docs/technical-design.md#5-authentication-and-personal-data-security).
- **Verify:** profile deletion resets onboarding and removes credits; account deletion cascades owned personal rows and revokes sessions; stale/recent-auth failures, concurrent edits and another user's/catalog rows are safe. Test replay logic on disposable snapshots; full packaged restoration rehearsal belongs to Task 25.

### Task 24: Evaluate quality, regressions, feedback and performance

**Depends:** 12, 17–19 and 21–23; permitted labeled snapshot for real recommendation-quality claims.
- **Deliver:** Python gold/retrieval evaluator consuming actual Go responses, sanitized Go timing/coverage counters and opt-in aggregate feedback. Follow [design §9](../docs/technical-design.md#9-verification-and-release-gates).
- **Verify:** all specified extraction/filter/rule/retrieval/grounding/conversation/latency gates with counts and per-language results. Regression tests cover source/model changes, privacy and ownership. Synthetic/fake results are reported separately; provider quality and real-source precision cannot be inferred from them.

### Task 25: Finish packaging, restoration and release acceptance

**Depends:** 23 and 24.
- **Deliver:** reproducible Go/frontend images, opt-in finite Python job image, explicit migration/maintenance steps, environment separation, documented startup/import/rebuild/backup/restore/rollback and reviewed scaffold retirement.
- **Verify:** packaged auth/profile/search/job flows, readiness/version mismatches, non-root/images without secrets, RAG role isolation, persistent volumes and a full encrypted-backup/deletion-ledger restore drill using a separate disposable Compose project. Apply all [release gates](../docs/technical-design.md#9-verification-and-release-gates) and [safe cutover rules](../docs/technical-design.md#10-local-environment-and-safe-cutover). Hosted deployment requires a concrete proposal and explicit approval.

## Blockers and decisions

| Decision or blocker | Needed before | Current position |
|---|---|---|
| Permitted syllabus/export acquisition, storage, excerpt display, embedding and redistribution | Real artifacts and quality evaluation | Prior robots restriction blocks automated syllabus crawling; [research record](../docs/discovery/source-access.md) defines the outstanding checks. |
| Pilot department/year/term, actual current period definitions and admission-cohort/track rules | Reliable real recommendations/category claims | Engineering/地球工学科/AY2026 is proposed; authoritative applicability is not established. |
| Rule reviewer, affiliation taxonomy and category mappings | Reviewed rule publication | Missing rules leave category outcomes unverified, without blocking profile or summary storage. |
| LLM/embedding provider, model/revision, data terms and budget | Live provider use and semantic release gate | Adapters/fakes allow independent progress; default proposals and budget measurement are in the design. |
| Mail sender, public origin/HTTPS and backup-window approval | Public launch | Fake mail/local development first; security defaults remain binding until explicitly revised. |
| Supported toolchain/library/image pins and vector/sqlc compatibility | Their owning engineering stages | Verify when introduced; a documented typed pgx fallback is allowed only for demonstrated vector compiler limitations. |

No source access, provider or deployment decision is needed to begin Task 1. Completing an early milestone does not waive remaining MVP acceptance criteria.
