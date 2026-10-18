# KCA Implementation Plan

Build one small verifiable increment at a time for three contributors. [Requirements](../docs/requirements.md) define scope; [technical design](../docs/technical-design.md) is the single implementation authority; [source access](../docs/discovery/source-access.md) records permissions and acquisition evidence. This roadmap does not authorize implementation during a documentation-only task.

## Goal, architecture and current progress

**Goal:** v0 register/sign in → save required Profile → one independent natural-language query → up to 3–5 course cards for Informatics master's students, AY2026 second semester. Validate a small permitted real sample before full target ingestion. v0.5 owns follow-ups/explanations/comparisons; expanded search, progress, course plans and timetables come later.

**Architecture:** React/TypeScript/Vite → Go/Gin (Auth, Profile, Course, minimal Curriculum, Chat; pgx/sqlc) → Python internal HTTP parsing/retrieval/ranking under `rag/`. Separate offline Python jobs publish academic metadata/documents/indexes. PostgreSQL/pgvector permissions, Go-owned migrations and state are defined in [design §§1–4](../docs/technical-design.md#1-architecture-and-ownership).

Observed on `update/docs`, 2026-10-18, before this documentation synchronization: clean working tree; local main points to the same commit, while the local `origin/main` reference is older. No fetch or assumption about current remote state. Preserve these existing artifacts:

- Go module/Gin dependency and `newRouter` scaffold exist. `TestHealthReturnsOK` already expresses the desired health contract; the router is empty and there is no runnable Go entry point. Test execution is unverified in this documentation task. Continue Task 1 from these files; do not recreate them or claim the endpoint is finished.
- Dockerfile/Compose scaffolding with PostgreSQL/pgvector, healthcheck and named `pgdata` volume exists. Backend image/command still target the removed legacy Python application; verified Go image/startup/DB availability remain pending. Existing configuration work is retained, not reset.
- Frontend and `rag/` are placeholders. No account/Profile/recommendation implementation or real course/handbook acquisition is evidenced.

Task numbers below are the revised sequence, not completion claims. Former tasks for credit summaries, full curriculum predicates, Go vector retrieval and conversations are superseded. No completed application milestone is inferred from expected commands in the previous plan.

## Milestones and contributor integration

- **Local foundation:** Tasks 1–5 separately; no course source or model provider needed.
- **Account/Profile flow:** Tasks 7–9 can proceed independently of academic acquisition.
- **Real sample pipeline:** Task 6 permission gate → Tasks 10–12 and 15, initially a few permitted records. Missing handbooks remain an explicit coverage state.
- **First end-to-end v0:** Tasks 13–19 integrate the sample with real auth/Profile guards and a live approved provider; complete target ingestion is not required. Fake-provider/synthetic demonstrations are labeled separately.
- **Full target collection:** Task 20 expands only after sample extraction/evidence checks pass.
- **Public release:** Task 21 adds the existing public gates; it does not block the first local/controlled milestone. No deployment is automatic.

| Contributor | Main work | Integration points |
|---|---|---|
| Frontend contributor | Task 9 auth/Profile, then 18 query/cards; contract mocks can begin from the design while Go is pending. | Configuration flags/parent codes, session/CSRF bootstrap, Profile revisions, response/coverage/error states. |
| Go contributor | Tasks 1–5, 7–8, 10, 12–13 and 17; define migrations/grants with Python input. | Canonical SQL, Python publication compatibility, candidate IDs/revisions, internal HTTP, database-filled cards. |
| Python contributor | Task 6, 11, 14–16; assist handbook mappings and 19–20. | Artifact/manifest format with Go, restricted writer/reader roles, parse/recommend DTOs, scoped evidence/IDs and metrics. |

Tasks 19 and 21 are shared integration checks, not separate services. Contributors can work independently against concise fixtures/contracts; each learner still follows the teaching checkpoints when that workflow is invoked.

## Shared workflow and review focus

Start each task with fresh branch/status/diff/file inspection; preserve unrelated work. Choose precise files/interfaces/dependency pins at task start and fetch current library documentation when applicable. Add meaningful failing checks for behavior, implement minimally, then run relevant checks/regressions. Completion requires observed output. Do not install the whole stack upfront or create another planning/handoff document.

Use Go unit/httptest and disposable PostgreSQL integration, Python pytest, frontend component/E2E checks. Real auth/Profile fixtures precede personalized HTTP checks; same-origin/CSRF/session/isolation fixtures precede authenticated E2E. Fakes prove engineering only. Test resources are disposable; never reset the existing volume. Commit/push/deploy require separate instructions.

Pay special attention to bypassed disabled controls, missing handbook coverage, personalization before budget cutoff, explicit exams versus unknown assessment, and provider/ID/reference failures. The tasks below own these checks. Do not introduce v0 conversation, progress, assessment-review state machines or full curriculum-engine tests.

## Tasks

### Task 1: Finish the local Go health endpoint

**Owner:** Go. **Depends:** none.
- **Deliver:** complete the existing router and add a runnable Gin process on port 8000, health JSON `{"status":"ok"}` and graceful shutdown. Reuse the existing test.
- **Verify:** run the focused handler test and local HTTP check; status/body/content type match. No DB/auth/sqlc/Python needed. Record actual commands/results; current test file alone is not a pass.

### Task 2: Build the minimal Go image

**Owner:** Go. **Depends:** 1.
- **Deliver:** Go multi-stage build/non-root runtime without Python or secrets; inspect legacy files before replacing the backend image.
- **Verify:** build/run in isolation, check health and shutdown/user/image exclusions. Preserve unrelated scaffold and the existing DB volume.

### Task 3: Wire Go and PostgreSQL/pgvector in Compose

**Owner:** Go. **Depends:** 2.
- **Deliver:** adapt existing Compose backend wiring while retaining DB configuration/volume and local ports. This does not yet add all later services.
- **Verify:** safe config validation, separate disposable project startup, DB health/extension availability and persistence after restart; never expose resolved secrets or use normal-volume resets.

### Task 4: Add pgx readiness

**Owner:** Go. **Depends:** 3.
- **Deliver:** bounded pgx pool/configuration and `/ready`, resource shutdown and sanitized errors.
- **Verify:** ready 200 with DB and 503 without it; health remains 200. No schema migration on startup.

### Task 5: Establish explicit migration execution

**Owner:** Go. **Depends:** 4.
- **Deliver:** single `backend/db/migrations/` sequence, minimal first migration/version check and explicit operator command. Add domain tables in their owning tasks.
- **Verify:** fresh disposable apply, repeat/upgrade/version-mismatch behavior; runtime never performs DDL. sqlc is generation only; no Python schema authority.

### Task 6: Validate a permitted Informatics sample

**Owner:** Python. **Depends:** none; parallel with foundation/accounts.
- **Deliver:** dated source-specific permission/use evidence and 3–5 representative course pages/export sections for AY2026 second semester; inspect formats, identifiers, encoding, meetings and assessment text. Locate applicable Informatics timetable/handbook sources separately.
- **Verify:** originals carry URL/year/time/hash and permitted storage/display/embedding decisions; human-check extracted field labels. Historical Engineering/ILAS timing is not new-scope proof. Missing source access blocks real acquisition only; missing handbook does not block basic sample discovery. No crawler before permission, no fabricated no-exam example.

### Task 7: Accounts, sessions and isolated clients

**Owner:** Go. **Depends:** 5.
- **Deliver:** account/session migration and register/login/logout/session APIs with hashing, cookies, CSRF/Origin, limits, genuine two-user fixtures and explicit local fake-mail mode under [design §7](../docs/technical-design.md#7-security-privacy-and-state).
- **Verify:** generic failures, ownership tampering, rotation/logout replay/expiry, restart persistence, limits and no personal logs. Public mail/restore checks belong to 21; essential security stays here.

### Task 8: Configuration and persistent Profile

**Owner:** Go. **Depends:** 7.
- **Deliver:** Go-maintained initialization for [design §2 configuration](../docs/technical-design.md#2-affiliation-configuration-and-profile), owner-scoped Profile/options APIs, revision checks, completion gate and reauthenticated Profile/account deletion. Keep earned credits/history absent.
- **Verify:** exact configuration, disabled-affiliation rejection, invalid parent/コース rejection even on direct requests, fixed master's level and independent years. Two users save/reload/edit across sessions/devices; next query loads edits; stale writes fail atomically. Deletion resets onboarding/revokes sessions without touching other users/catalog. No handbook dependency.

### Task 9: Frontend account and Profile journey

**Owner:** Frontend. **Depends:** 8 for live integration; mocks may start from design now.
- **Deliver:** Vite/TypeScript base, same-origin cookie/CSRF client, session/Profile bootstrap, separate registration/setup/edit forms and reauthenticated deletion controls, grouped disabled affiliation options and isolated browser fixtures.
- **Verify:** register → sign in → required Profile save; omitted fields block personalization. Disabled controls have labels and backend errors; seven コース are named correctly. Test reload/edit/revision conflict/logout/account switch and personal-state clearing. No credit inputs.

### Task 10: Minimal academic schema and shared artifact contract

**Owner:** Go with Python. **Depends:** 5; 6 informs real-source fields.
- **Deliver:** focused course/offering/source/meeting and document/index/mapping-link migrations plus small versioned artifact contract. Define separate Go, configuration-admin, Python-online and offline-writer grants from [design §§3–4](../docs/technical-design.md#3-database-permissions-and-migrations).
- **Verify:** stable source/course/offering identity, source/section/year/revision links, exact credits/nulls and multiple meetings; forbidden personal-table/DDL access denied for both Python roles and academic writes denied online. No assessment classification table, category engine or duplicate schema.

### Task 11: Parse and publish the first sample

**Owner:** Python. **Depends:** 10; 6 for real data.
- **Deliver:** focused manifest/acquisition/parsing/normalization/publication modules under `rag/`; preserve original assessment text, Markdown/locators and file/object keys. Choose HTML/PDF adapters from observed sample, not assumptions.
- **Verify:** human-labeled encoding/tables/meetings/credits/source fields; corrupt/unknown cases stay unknown. Repeat publication is idempotent; revisions publish atomically, missing year blocks publication, partial sync does not withdraw data, changes invalidate derivatives. Use restricted offline writes, never migrations.

### Task 12: Go Course reader and parameterized SQL

**Owner:** Go. **Depends:** 10; 11 for real compatibility.
- **Deliver:** sqlc generation/pgx reader, compact candidate DTO and same-meeting/term/year SQL with must/prefer separation. No public standalone search UI/API required.
- **Verify:** Python-write/Go-read UUID/UTF-8/timestamp/decimal/null/reference compatibility on PostgreSQL; meeting ranges, additional/unknown meetings, empty scope and no silent relaxation. Offering unit is not automatically restricted to Profile affiliation; no model SQL or mandatory exam SQL.

### Task 13: Minimal Curriculum scope lookup

**Owner:** Go with Python source review. **Depends:** 10; 6 only for real reviewed mappings.
- **Deliver:** manually confirmed school/コース/degree/admission-year → handbook document/section resolution, with source/effective scope. Implement `missing`/`ambiguous` outcomes before real documents exist.
- **Verify:** cohort boundaries, overlap/missing mappings and source versions; no newest-wins guess or full predicate engine. Missing handbook allows discovery and emits the alignment warning. Real handbook mapping is optional for the first basic discovery flow, never faked.

### Task 14: Internal Profile-aware parsing API

**Owner:** Python. **Depends:** design contract; 8 supplies real projected context at integration.
- **Deliver:** `/internal/parse-intent`, private service authentication, provider adapter/Instructor schemas, filters/semantic/must/prefer/unresolved response, sanitized context and bounded failures.
- **Verify:** Profile is in the first call; broad requests need no forced topic. Japanese/Chinese/English negation/must/prefer, invalid model keys/output, prompt injection and unavailable provider. Unresolved timetable wording prompts clarification via Go; no model Profile writes. Fake provider first, approved live provider later.

### Task 15: Academic documents, summaries and vector baseline

**Owner:** Python. **Depends:** 11; approved model/processing rights for real embeddings.
- **Deliver:** offline document/section chunks, compact course-summary indexes and embeddings with model/version/dimension/preprocessing manifests and explicit offering/revision links. Local files suffice for development.
- **Verify:** source/section/semester isolation, changed-source rebuild, partial coverage/model mismatch and Python online read compatibility. No student context in offline artifacts; no new LLM calls assumed necessary for document reading or baseline summaries. Generated summaries, if later used, have separately measured costs.

### Task 16: Scoped Python recommendation API

**Owner:** Python. **Depends:** 12 candidate contract, 13 coverage contract, 14 provider base and 15 retrieval data; reviewed real handbooks optional.
- **Deliver:** `/internal/recommend` selects by Profile/handbook/intent within candidate IDs before applying a configurable inspection budget, then vector retrieval/direct section reading and one generative evidence-ranking call. Return IDs/reasons/references/uncertainty and inspected scope/counts. Add a minimal local Python online image/Compose service with read-only credentials; offline execution stays separate, with final packaging at 21.
- **Verify:** relevant Profile candidate survives generic-ranking adversarial fixtures; no arbitrary first-N cutoff, cross-candidate/semester reference or relational SQL refiltering. Missing handbooks work; unknown exam evidence is separate, explicit exams excluded for strict no-exam. Check batch-empty/partial, private credential and provider failures. Hybrid, if justified later, scopes both branches.

### Task 17: Go Chat orchestration and validated cards

**Owner:** Go. **Depends:** 8, 12–14 and 16 (fake internal responses first).
- **Deliver:** independent `/api/chat` request: fresh Profile/target, minimized parsing call, validated SQL, handbook scope, compact RAG request, response checks and database-filled cards. Add internal HTTP deadlines, safe errors and aggregate metrics.
- **Verify:** real auth/completion/CSRF guards, saved edits used next request, no identity/secrets/full Profile sent, no generated SQL, verified timetable or clarification. Validate candidate/course pairs, active evidence and facts; invalid response fails safely. Distinguish SQL-empty from inspected-batch-empty. Two generative calls normally, separate retries/embeddings; no durable chat or expansion. Schema checks do not prove all claim semantics.

### Task 18: Single-query discovery UI

**Owner:** Frontend. **Depends:** 9 and 17; contract mocks can precede live RAG.
- **Deliver:** explicit AY2026 second-semester scope, one query input, at most five database-fact cards with short reason/source, separate “Needs confirmation,” coverage/partial/empty/clarification/error states and accessible Japanese-first interaction.
- **Verify:** form-based auth → saved Profile → independent query → cards on the permitted sample; incomplete/updated Profile and cross-user switching. Clarification requires a complete resubmission; refresh may discard question/results. No follow-up/comparison/show-more/credit interface or remembered filters.

### Task 19: Fixed sample evaluation and first end-to-end acceptance

**Owner:** all three. **Depends:** 6, 17–18 and their prerequisites; approved live provider for real model claims.
- **Deliver:** small fixed deterministic cases and permitted labeled sample, Go/Python/UI integration results, human-checked claims/evidence and actual token/latency/candidate-budget measurements under [controlled gates](../docs/technical-design.md#9-verification-and-release-gates).
- **Verify:** availability/relationships, extraction/negation/schedule/must-prefer, Profile-before-truncation, missing handbooks, boundaries/unknowns, IDs/evidence/partial/empty, two users/edits and failures. State sample coverage explicitly; no fabricated quality or startup results. This milestone does not wait for complete target ingestion, public mail/restore or v0.5 tests.

### Task 20: Expand to the complete target collection

**Owner:** Python, with Go metadata checks and shared source review. **Depends:** 19 and permitted target inventory.
- **Deliver:** expand the validated pipeline to the agreed complete Informatics-target AY2026 second-semester offering collection, with scope/IDs/sections reconciled against an authorized source inventory. Preconfigured other affiliations do not broaden collection.
- **Verify:** inventory counts/gaps, repeat/partial/update behavior, encoding/field/evidence checks and scoped retrieval regression. Unknown inventory/completeness remains disclosed, not “complete.” Acquire applicable handbooks only through permitted paths and reviewed applicability; their absence does not stop basic discovery.

### Task 21: Public-release security, quality and packaging

**Owner:** all three. **Depends:** 19–20; mail/origin/provider/source/hosting decisions for the parts that require them.
- **Deliver:** verification/recovery mail actions, public HTTPS/secrets, minimal Go/Python-online/frontend images plus explicit offline job execution; startup/migration/import/index/rollback and backup/deletion-replay runbooks. Apply [public gates](../docs/technical-design.md#9-verification-and-release-gates), not the former conversation/rule-engine MVP gates.
- **Verify:** frozen public benchmark and measured load/cost, packaged flows, denied role access, real token/recovery behavior and disposable encrypted restore/deletion replay. No normal-volume reset. Report unmet gates; do not deploy or incur hosting costs without separate authorization.

## Open evidence and bounded blockers

| Evidence/decision | Blocks | Does not block |
|---|---|---|
| Source-specific permitted acquisition/storage/display/embedding and real sample format | Real academic ingestion/quality claims | Foundation, accounts/Profile and synthetic contract checks |
| Complete target inventory | Claim of full target collection at 20 | Small-sample first end-to-end milestone |
| Official confirmation of configured affiliation names/relationships | Claim of official verification | User-specified configuration implementation with its unverified label |
| Applicable Informatics period definitions | Automatic mapping of ambiguous timetable phrases | Clear numeric-slot queries or clarification |
| Reviewed handbook editions/mappings | Curriculum alignment claims | Basic discovery with missing-coverage warning |
| Provider/model/data terms and measured budget | Live LLM/embedding use and real quality claims | Fake adapters and deterministic boundary tests |
| Public mail/HTTPS/backup/hosting | Public release/deployment | Controlled local flow with essential security |

Next smallest implementation work is Task 1: run the existing health test, implement its route plus the minimal Go entry point, and observe local health/shutdown results through the learning workflow if requested. This documentation task executes none of those commands.
