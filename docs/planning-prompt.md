# KCA Planning Instructions

Read `AGENTS.md`, `docs/requirements.md`, `plans/implementation-plan.md`, `docs/technical-design.md` and `docs/discovery/source-access.md` before changing the architecture or roadmap. Product requirements remain binding; technology candidates in requirements §8 do not override the selected architecture.

## Document responsibilities

- `plans/implementation-plan.md`: concise architecture overview, incremental development order, milestones, tasks/dependencies, deliverables, critical edge cases, acceptance criteria and current blockers.
- `docs/technical-design.md`: the single shared implementation definition for architecture boundaries, migration/schema ownership, models, assessment acceptance, search/rules, security/privacy, credits, overrides, API contracts and release thresholds. Reference it from tasks instead of repeating it.
- `docs/discovery/source-access.md`: dated official-source observations, access/rights decisions, unresolved evidence and the permitted sample workflow. Distinguish verified observations from assumptions and untested feasibility.
- `docs/requirements.md`: product scope and constraints. Do not silently weaken or rewrite requirements to simplify a plan.

## Planning checklist

1. Keep Go/Gin + pgx/sqlc in `backend/`, finite Python data/RAG jobs in top-level `rag/`, React/TypeScript/Vite in `frontend/`, and shared PostgreSQL/pgvector. Preserve existing roots/Compose entry point; avoid microservices, a permanent Python API and top-level `ingestion/`.
2. State migration authority and runtime data privileges once in technical design. Introduce versioned SQL before dependent schemas, sqlc only after representative queries exist, and Python-write/Go-read/vector compatibility at the owning stages. No second ORM/Alembic schema authority.
3. Start with five separate increments: local health endpoint → minimal Go image → Compose PostgreSQL/pgvector → pgx readiness → first versioned migration. Keep academic schema, sqlc, vectors and cross-runtime tests later; final packaging/restore/release gates do not block the first local Docker exercise.
4. Allow source-permission research independently. Check official pages/IDs/encoding/navigation, robots/terms/rate/reuse constraints and cohort/year-specific curriculum/period documents before real collection. Honor recorded access restrictions; never bypass controls, crawl at scale or claim a crawler proven from an index page. Synthetic fixtures demonstrate engineering only.
5. Preserve source provenance/versioning, candidate-versus-accepted assessment states, unknowns, deterministic hard filters and reviewed curriculum predicates. Semantic/RAG retrieval cannot decide hard constraints or graduation completion. Keep multilingual quality, citations and every existing acceptance threshold.
6. Include core MVP registration/login/logout, persistent Student Profiles, onboarding, profile editing/deletion and credit management. Required fields: faculty, department, grade, enrollment year; optional track and all three self-reported credit summaries remain skippable and saveable without curriculum coverage. Category reports and verified progress remain distinct.
7. Keep authentication/user ownership, session/CSRF, revision conflicts, temporary client-confirmed overrides, no model-authorized persistence, provider allowlist/redaction, credit minimization and deletion/backup policy covered by the technical design and owning tests.
8. Require genuine two-user auth/onboarding fixtures before personalized HTTP checks and session bootstrap/CSRF/proxy/isolated browser fixtures before feature E2E. Include Go business-rule/httptest/integration checks, Python artifact/job/retrieval tests, UI checks and late load/restore/release evaluation.
9. Preserve observable sync/coverage/assessment/rule/retrieval/provider metrics and aggregate feedback without personal logs. Keep source permission, scope/reviewer, provider/model/data terms/budget, mail/origin and deployment choices visible as blockers for only the stages that depend on them.

## Task format and shared workflow

Each task needs an objective, explicit dependencies, a small useful deliverable, critical edge cases and observable completion criteria. Choose precise future files/signatures/dependency pins and task-specific commands at task start, rather than prewriting the application. State test/implementation/review workflow once. Remove repetitive failing-test/implement/rerun/stage/commit instructions, redundant assertions and delivery/self-review narration. No mandatory unavailable execution skills.

Use the writing-plans workflow where useful, adapting its default template to this concise format and the owner's latest instructions. Keep first steps small enough for a personal learning project; split oversized combined infrastructure/features instead of adding ceremony.

Before finishing, compare the previous plan and requirements with the revised roadmap/design for coverage, security/unknown-state semantics, edge cases and exact thresholds. Check authoritative definitions, dependency validity, links and section references. Planning is documentation-only: no application implementation, dependency installation, database operations, bulk crawling, deployment or commits. Begin future implementation with a fresh Git status/file review instead of preserving transient working-tree notes in long-lived instructions.
