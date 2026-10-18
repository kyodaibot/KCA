# KCA Planning Instructions

Read [AGENTS.md](../AGENTS.md), [requirements](requirements.md), [technical design](technical-design.md), [implementation plan](../plans/implementation-plan.md) and [source access](discovery/source-access.md) before changing the roadmap. The current product decisions supersede the former oversized MVP. Keep this a three-person student project with small, verifiable tasks and controlled effort/token usage.

## Document responsibilities

- Requirements define product behavior, release scope and acceptance.
- Technical design is the single shared definition of architecture, data permissions, interfaces and state, including the authoritative affiliation/コース configuration. Reference it rather than duplicating lists/contracts.
- Implementation plan defines tasks, dependencies, contributor work, integration points, observable verification and actual progress. Preserve completed work; do not equate a scaffold or expected output with completion.
- Source access retains dated observations/restrictions and pending Informatics verification. Historical Engineering material is not evidence for the new scope.
- [Learning workflow](../plans/learning-workflow.md) controls teaching pace when invoked. Preserve its step approvals for learning tasks; they do not block separately authorized documentation updates.

## Binding planning constraints

1. **v0:** register/sign in → save required Profile → one independent natural-language query → up to 3–5 course cards, for Informatics master's students and AY2026 second semester. Validate a small permitted real sample before full target ingestion; full ingestion is not needed for the first working flow. Account and infrastructure work proceed independently of acquisition.
2. **v0.5:** follow-ups, explanations and comparisons. **Later:** show more/expanded search, manual-search UI, progress, plans and timetables. No v0 earned credits, completed history, conversation memory, temporary Profile overrides, expansion button or full rule engine. Planned credits never become earned credits. No university login/registration automation or graduation certification.
3. Retain React/TypeScript/Vite in `frontend/`, Go/Gin + pgx/sqlc in `backend/`, Python internal online HTTP API and separate offline jobs in `rag/`, PostgreSQL/pgvector and Compose. No top-level `ingestion/`, extra microservices, Redis, Kubernetes or complex agents. Go business modules are Auth, Profile, Course, Curriculum and Chat; defer StudyPlan.
4. Go owns migrations and personal data. Python online reads academic retrieval data only; offline ingestion uses separate restricted academic writes. Neither accesses personal tables or owns another schema. Use internal HTTP, candidate offering IDs and compact metadata; Python reads linked documents/indexes without repeating Go's SQL hard filtering.
5. Configure all affiliation options once in the design; only 情報学研究科 is enabled. Use stable codes, types and availability flags, seven required コース and backend relationship validation. Names are user-specified, not verified official facts. Availability, corpus coverage and handbook coverage remain separate. Do not add other affiliations' internal subdivisions.
6. Profile participates in the first parsing call and selection before truncation. Go validates fields and parameterizes supported SQL; the model never supplies executable SQL. Reviewed manual handbook mappings resolve scope; missing handbooks allow discovery with a curriculum warning. No mandatory assessment classification table/review machine or complete category mapping.
7. Python performs scoped retrieval/evidence-based ranking. Candidate budgets are configurable and disclosed; hybrid branches share the scope. Unknown evidence means “Needs confirmation,” never hard-condition satisfaction; explicit exam requirements fail strict no-exam requests. Distinguish empty SQL scope from no inspected-batch match. Do not imply exhaustive search or silently relax requirements.
8. Normal requests use two generative calls (parse and recommend). Count embeddings, retries and other processing separately; document reading/retrieval need no inherent extra LLM call. Instructor/schema/ID checks validate structure/boundaries, not all claim semantics.
9. PostgreSQL persists auth sessions and Profiles; page memory holds the current question/results. Future browser conversation context is untrusted. Retain essential security, minimal provider payloads, evidence/unknown handling and user isolation. Separate controlled-prototype checks from existing public mail/HTTPS/quality/backup gates; do not add unrelated workflows.
10. Fixed tests target availability/relationships, intent/negation/schedule/must-prefer, Profile-aware retrieval/budgets, missing handbooks, IDs/evidence/unknowns/partial/empty results, two-user isolation, edits and failures. Measure actual tokens/latency; never invent results. v0.5 tests are not v0 gates.

## Task and review format

Each task states objective, owner (frontend/Go/Python), dependencies, a small useful deliverable and observable completion criteria with critical edge cases. Define precise files/interfaces and verify current library documentation at implementation time; do not prewrite application code or install the future stack while planning. Keep test-first guidance proportional and in one shared workflow. No mandatory unavailable execution skills, automatic commits or new competing handoff/design documents.

Before finishing, check scope coverage, dependency order, diagrams, links, authoritative definitions and acceptance criteria across Markdown, including ignored local prompts when relevant. Documentation-only work does not authorize code/migration/Compose edits, dependency installation, data downloads, database operations, deployment or commit/push. Future implementation starts with fresh branch/status/diff inspection and the existing learning checkpoints if requested.
