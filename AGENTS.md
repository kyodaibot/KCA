# AGENTS.md — Kyoto Course Agent

## Project and documentation

KCA is a three-person student project. Build a small, complete, verifiable course-discovery flow for Kyoto University Graduate School of Informatics master's students, initially AY2026 second semester. v0 is register/sign in → save Profile → one independent natural-language query → up to 3–5 course cards. Do not restore the former oversized MVP.

Read the relevant documents before changing scope or architecture:

- [Requirements](docs/requirements.md): product behavior, v0/v0.5/later scope and acceptance.
- [Technical design](docs/technical-design.md): single authority for architecture, data permissions, interfaces and state; also the sole affiliation/コース configuration definition.
- [Implementation plan](plans/implementation-plan.md): small tasks, dependencies, deliverables and observed progress.
- [Planning instructions](docs/planning-prompt.md): constraints for future roadmap changes.
- [Source access](docs/discovery/source-access.md): dated observations, restrictions and unverified acquisition tasks.
- [Learning workflow](plans/learning-workflow.md): teaching checkpoints when the learner invokes that workflow; these do not block separately authorized documentation work.

Inspect the current branch, status, relevant diffs and actual files before work. Preserve unrelated changes and genuinely completed work; do not infer completion from a plan or assume this branch matches main. Prefer simple modules and short, testable increments. Explain new dependencies and significant trade-offs. Planning/documentation-only requests do not authorize application code, installs, database operations or deployment.

## Architecture boundaries

- `frontend/`: React + TypeScript + Vite, presentation and same-origin API client.
- `backend/`: Go + Gin with pgx + sqlc. Business modules are **Auth, Profile, Course, Curriculum, Chat**. Handlers, middleware, repositories and the RAG client support these modules. Defer StudyPlan.
- `rag/`: Python internal online HTTP API for intent extraction, scoped retrieval and evidence-based ranking, plus separate offline ingestion/chunking/embedding jobs. Go orchestrates requests; Python owns online query embeddings/vector retrieval. Keep provider integration replaceable.
- PostgreSQL + pgvector is shared with different runtime privileges. Go owns the only migration history under `backend/db/migrations/`; run migrations explicitly. Neither runtime migrates on startup. No second Alembic/ORM schema authority; sqlc generates query code, not migrations.
- Go alone accesses accounts, sessions and Profiles. Python online reads academic retrieval data only. Offline ingestion has separate restricted academic write credentials. Neither Python runtime accesses personal tables. Follow the ownership table in the technical design.
- Preserve the existing roots and Compose entry point. No top-level `ingestion/`, Redis, Kubernetes, complex agent orchestration or additional microservices.

## Scope and recommendation rules

- Registration and Profile setup are separate. Required v0 academic fields and available affiliations/コース are defined in the linked specifications. Enforce availability and parent relationships on the backend; visible disabled options do not enable service or imply data coverage.
- Load the latest server Profile for every personalized query. Year of study, admission year and offering year/semester are independent. v0 collects no earned-credit summaries, category reports or completed-course histories.
- Profile participates in intent extraction and candidate selection before any budget cutoff. Go validates model fields and applies supported hard conditions through parameterized SQL; never execute model-generated SQL. Academic relevance is normally a soft signal, not an offering-unit equality filter.
- Pass candidate offering IDs and compact metadata to Python, not complete syllabi. Python reads linked documents/indexes within that scope and does not repeat Go's SQL filtering. Both branches of any hybrid retrieval respect the same IDs.
- Resolve handbook scope only through reviewed manual mappings. Curriculum is minimal applicability mapping, not a complete predicate engine. Missing coverage permits discovery with an explicit curriculum-alignment warning.
- Facts on cards come from the database. Reasons and evidence must reference current candidate offerings/documents; never invent facts or curriculum guarantees. Vector similarity and schema validation do not prove assessment claims. Label insufficient evidence “Needs confirmation”; exclude explicit hard-condition failures even from that group.
- Candidate budgets are configurable, Profile-aware and disclosed as partial inspection. Distinguish no SQL candidates from no matches in the inspected batch. v0 has no expansion button, automatic expansion or silent relaxation.
- Two generative LLM calls are the normal path: parse intent and recommend. Track embeddings, retries, tokens and latency separately. Clarification ends that independent request; the user resubmits a complete query.
- v0.5 owns follow-ups, explanations and comparisons. Later work owns expanded search, manual search UI, progress, course plans and timetables; planned credits never count as earned credits. No automated university login/registration or graduation certification.

## Source handling

Investigate permitted public access before collection; honor robots, terms and access controls. Do not bypass authentication/CAPTCHA, guess IDs or equate reachable pages with permission. Preserve dated historical restrictions in the source-access record. Validate a small permitted real sample and extraction formats before expanding to the target collection. No real course/handbook artifacts have yet been acquired; infrastructure and account work can proceed independently.

Keep raw sources, normalized metadata, documents, chunks and embeddings distinct, linked by explicit IDs and source/version/section/year provenance. Use local files in development or object storage. Preserve original URLs, retrieval times and unknown/conflicting fields. Do not require an exam classification table, assessment-review state machine or complete category mapping for v0.

## Implementation, security and verification

- Keep handlers thin and services small. Pass `context.Context`, bound DB/internal HTTP/provider work, validate input and sanitize errors. Use parameterized pgx/sqlc queries and explicit transactions; separate generated rows from public DTOs. Store course credits exactly, preserve nulls, and avoid float64 for credit amounts.
- Python uses type annotations, focused acquisition/parsing/normalization/retrieval modules and Pydantic/Instructor structure validation. Structure checks do not establish semantic truth. Use restricted parameterized writes offline and read-only connections online.
- TypeScript uses explicit API types and separates UI from API communication. Define concise contracts in the technical design; formalize only implemented endpoints when needed.
- Resolve personal ownership from the authenticated session; never trust client owner IDs. Preserve password hashing, revocable cookies/sessions, CSRF/Origin checks, bounded attempts, user isolation, Profile edit conflicts and deletion. Public-release mail/HTTPS/backup requirements are distinct from controlled local gates; see the design.
- Send only server-projected necessary academic fields and minimized request text to Python/providers. Exclude identifiers, email, passwords, tokens, full Profiles and earned-credit information; keep personal content out of logs, shared caches and embeddings. If safe minimization fails, ask for reformulation. Model output cannot authorize Profile edits.
- Persist auth sessions and Profiles in PostgreSQL. Keep the current question/results only in page memory. No v0 chat sessions/history, accumulated filters, previous-result lists or temporary academic-override workflow. Clear personal UI state on logout/account switch. v0.5 client context remains untrusted and server Profile remains authoritative.
- Add meaningful business/security tests: Go unit/httptest and disposable PostgreSQL integration, Python pytest, frontend component/E2E. Auth/Profile fixtures precede personalized HTTP checks; same-origin/CSRF/session browser fixtures precede authenticated E2E. Test candidate boundaries, evidence, unknowns and failures independently of model quality.
- Use current documentation for library-specific implementation questions according to applicable Context7 instructions; documentation synchronization alone does not require installing tools or fetching course data. Use writing-plans, test-driven-development and systematic-debugging when applicable.
- Run relevant checks and report only actual results. Do not claim startup, quality or performance success without evidence. Keep images minimal/non-root, preserve the existing database volume and use disposable resources for tests. Production access, destructive database actions and deployment require explicit authorization; commits/pushes follow separate instructions.

When delivering work, briefly state changes/files, executed checks, unresolved limitations and the next smallest task. Explain in plain language and control effort and token usage.
