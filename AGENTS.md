# AGENTS.md — Kyoto Course Agent

## 1. Project Overview

Kyoto Course Agent is an AI-powered course recommendation chatbot for Kyoto University students.

The system helps students discover suitable courses based on:

- Faculty, department, and academic year
- Semester, weekday, and class period
- Course content and personal interests
- Assessment methods, including examination requirements
- Curriculum requirements and credit eligibility

### Project Documents

- `docs/requirements.md` — Product requirements and constraints
- `docs/planning-prompt.md` — Architecture planning instructions
- `plans/implementation-plan.md` — Implementation roadmap, once created
- `docs/technical-design.md` — Authoritative shared implementation rules, contracts and release gates
- `docs/discovery/source-access.md` — Dated source observations, restrictions and acquisition checks

Read relevant documents before making architectural decisions.

## 2. General Development Principles

- Prefer simple, maintainable solutions over unnecessary complexity.
- Follow the existing project architecture and conventions.
- Do not introduce dependencies without explaining their purpose.
- Avoid overengineering and premature optimization.
- Keep components modular and independently testable.
- Do not assume unspecified requirements.
- Explain significant architectural trade-offs.
- Do not silently change previously agreed requirements.

## 3. Planning First

For complex features or architectural changes:

1. Understand the requirements and existing implementation.
2. Identify dependencies, constraints, and potential risks.
3. Propose an implementation plan.
4. Break work into small, testable tasks.
5. Define acceptance criteria before implementation.

Use the `writing-plans` skill when appropriate.

For planning-only requests, do not write application code.

For small, well-scoped changes, proceed directly unless the user requests planning.

## 4. Technology Stack

The selected architecture is:

- Backend: Go + Gin for the API and business logic; pgx + sqlc for explicit, typed PostgreSQL queries (rationale in `docs/technical-design.md`).
- Offline RAG/data processing: Python under the top-level `rag/` directory, including syllabus ingestion, HTML/PDF parsing, Markdown conversion, chunking, embeddings and retrieval evaluation.
- Frontend: React + TypeScript + Vite under `frontend/`.
- Database: PostgreSQL
- Vector Search: pgvector
- LLM: Provider-agnostic integration
- Containerization: Docker
- Testing: Go tests/httptest for backend logic and API, pytest for `rag/`, and frontend component/E2E tests.

Keep `backend/` for Go, `rag/` for Python jobs, and the existing `docs/`, `plans/` and Compose entry point. Do not create a top-level `ingestion/` directory or a permanent Python API service. Go owns authentication, profiles, filtering, deterministic curriculum evaluation, online retrieval and recommendation orchestration. Python jobs own offline source artifacts and academic-data publication only.

The backend owns the single versioned SQL migration sequence under `backend/db/migrations/`; both runtimes consume that schema, and neither runtime auto-migrates on startup. sqlc generates Go query code; it does not execute migrations. Python has no separate Alembic/ORM schema or access to identity/profile tables. Version shared artifact contracts under `docs/contracts/` and test Python-produced records/vectors with Go readers. Avoid infrastructure unnecessary for MVP. The Python/FastAPI candidate in requirements §8 is historical and superseded by this explicit decision; product requirements are unchanged.

## 5. Data Collection and Processing

The main data source is Kyoto University's public syllabus system:

https://www.k.kyoto-u.ac.jp/external/open_syllabus/top

### Data Collection Rules

- Investigate available public data-access methods before implementing a scraper.
- Respect access restrictions, robots policies, and applicable terms.
- Do not bypass authentication, CAPTCHA, or access controls.
- Prefer structured data when available.
- Preserve original source URLs and retrieval timestamps.
- Separate raw data extraction from normalization.
- Support incremental updates when practical.
- Do not invent missing syllabus information.

### Data Storage

Distinguish between:

1. Original source data
2. Normalized structured course records
3. Markdown documents
4. Embeddings and vector indexes

Every normalized record should retain traceability to its original source.

## 6. Retrieval and Recommendation Architecture

Do not rely exclusively on vector similarity for course selection.

Use the appropriate mechanism for each task:

### Structured Filtering

Use structured queries for:

- Academic year and semester
- Faculty and department
- Weekday and class period
- Course identifiers
- Explicit assessment constraints

### Curriculum Rule Evaluation

Use verifiable rules for:

- Mandatory and elective requirements
- Department-specific eligibility
- Credit categories
- Graduation requirements

Rules must be associated with their applicable academic year and curriculum version.

### Semantic Retrieval

Use embeddings for:

- Course content similarity
- Natural-language interests
- Topic-based recommendations
- Matching informal descriptions to course content

LLMs may interpret user intent and explain results, but must not invent academic rules or override verified constraints.

## 7. Recommendation Accuracy

- Never fabricate courses, instructors, credits, schedules, or assessment methods.
- Distinguish confirmed facts from inferred information.
- Include source references when presenting course recommendations.
- Clearly report unknown or ambiguous attributes.
- Do not claim a course satisfies graduation requirements without verified evidence.
- Treat conflicting source data as an unresolved issue.
- Prefer asking for missing critical information over giving an unreliable answer.
- Distinguish recommendations from official academic guidance.

## 8. Coding Standards

### Go (`backend/`)

- Use small feature packages and explicit service/repository interfaces. Keep Gin handlers thin; authentication, profile ownership, curriculum predicates and ranking belong in Go services.
- Pass `context.Context` through database/provider calls, wrap errors without personal data, and enforce request timeouts and input validation.
- Use parameterized sqlc queries with pgx transactions. Keep generated database structs separate from public API DTOs; never use GORM AutoMigrate or a second schema authority.
- Represent credit decimals exactly and preserve omitted/null/zero JSON semantics; do not use float64 for credits. Define APIs in `docs/contracts/openapi.yaml` and validate Go/frontend implementations against them.
- Use gofmt, Go unit tests, real PostgreSQL integration tests and httptest for middleware/ownership/session behavior.

### Python (`rag/`)

- Use type annotations.
- Follow PEP 8.
- Prefer small, focused functions.
- Use Pydantic for offline artifact/manifest validation; Python does not implement the public API or personal-data business logic.
- Keep acquisition, parsing, normalization, Markdown conversion, chunking, embedding and evaluation separate. Use psycopg parameterized writes against the backend-owned schema; never create/mutate schemas from jobs.
- Use structured logging where appropriate.

### TypeScript

- Use TypeScript rather than untyped JavaScript.
- Define explicit types for API contracts.
- Keep UI components modular.
- Separate presentation logic from API communication.

### Database

- Use versioned migrations for schema changes.
- Use parameterized queries.
- Add indexes based on actual query requirements.
- Keep curriculum rules separate from course descriptions.

## 9. Testing and Verification

- Add tests for meaningful business logic.
- Use Go unit/integration/httptest checks for backend services and routes, pytest for Python RAG fixtures/jobs, and frontend component/E2E checks. Test cross-runtime schema/vector/decimal/null compatibility against one disposable PostgreSQL instance.
- Test course filtering and curriculum rules independently from the LLM.
- Include edge cases, missing fields, and conflicting data.
- Use deterministic fixtures for recommendation logic.
- Run relevant tests after making changes.
- Authentication/session/profile fixtures must exist before personalized HTTP checks; session bootstrap, CSRF-aware client, same-origin proxy and isolated browser fixtures must exist before authenticated E2E checks.
- Never claim tests passed unless they were executed successfully.

Use `test-driven-development` and `systematic-debugging` when appropriate.

## 10. Security and Privacy

- Never commit API keys, credentials, or personal information.
- Use environment variables or appropriate secret management.
- Validate external input.
- Apply least-privilege database access.
- Do not access production resources without explicit authorization.
- Do not perform destructive database operations without explicit approval.
- Avoid collecting unnecessary student personal data.
- Resolve account ownership from the authenticated server session; enforce owner-scoped access for every personal-data operation and never trust a client-supplied owner ID.
- Keep identity, secrets and personal profile content out of logs, analytics, shared caches and embeddings; provide explicit personal-data deletion and document backup retention. Provider payloads use a server-projected academic allowlist: faculty/department codes or names, grade, enrollment cohort, optional track, and requested offering academic year/term, each only when needed. Exclude account/student identifiers, email addresses, passwords, tokens, unnecessary personal information and full profile objects. Remove identity/secrets and unrequested credit information from free-text and prior context before provider calls; ask for reformulation when a safe minimal request cannot be formed.
- Exclude credits from LLM payloads by default. For an explicit credit/progress question, send only the requested total/major/general summary or relevant category amount with a `self_reported`/unverified label. Apply the same allowlist to hypothetical context and label it hypothetical.
- Keep profiles separate from verified curriculum rules. Grade, enrollment cohort and offering academic year are independent; self-reported credits are not proof of course completion or graduation eligibility.
- Allow self-reported total, major/specialized and general education summaries without curriculum-rule coverage; keep them separate from versioned category mappings and officially validated progress.
- Load persisted profiles on the server for personalized requests; require completed onboarding and invalidate stale profile-derived context after edits. Do not assume university SSO access.
- Temporary academic overrides require explicit client confirmation of reviewed values and the current profile revision; model output cannot supply authorization. They never write saved profiles. Restrict persistent profile updates to authenticated, owner-scoped profile-management endpoints.

## 11. Docker and Deployment

- Prefer reproducible development environments.
- Keep Docker images minimal.
- Use `.dockerignore` and non-root execution where appropriate.
- Separate development and production configurations.
- Plan a Go multi-stage backend image and an opt-in one-shot Python `rag` job image. Preserve the existing PostgreSQL volume; run migrations as an explicit operator/release step before either runtime uses a new schema.
- Provide health checks for deployed services.
- Never deploy to production without explicit approval.

Use `docker-devops` when appropriate.

## 12. Communication and Deliverables

When completing a development task:

1. Briefly summarize the changes.
2. List important files created or modified.
3. Report tests executed and their results.
4. Mention unresolved issues or limitations.
5. Suggest the next logical development step.

Prefer concise technical explanations.

For architecture decisions, explain the rationale and trade-offs before introducing unnecessary complexity.
