# Kyoto University Course Selection Agent — Requirements

## 1. Product and release scope

KCA is a course-discovery assistant built by three students. Prioritize a small, complete, verifiable flow. Recommendations support decisions; they are not official academic guidance, university registration or graduation certification. [Technical design](technical-design.md) defines implementation boundaries; [the roadmap](../plans/implementation-plan.md) defines incremental delivery.

| Release | Scope |
|---|---|
| **v0** | Register/sign in → complete and save Profile → submit one independent natural-language course query → receive up to 3–5 source-linked course cards. Initially serve only master's students in Kyoto University's Graduate School of Informatics, targeting AY2026 second semester. |
| **v0.5** | Conversational follow-ups, course explanations and course comparisons. |
| **Later** | Show more results, expand search, academic progress tracking, standalone manual search, personal course plans and timetables. Planned credits must not count as earned credits. |

v0 collects no total/category-level earned credits or completed-course histories. Do not implement automated university-system login/registration, graduation certification or a complete curriculum rule engine. Preconfigured affiliations do not broaden target users or acquisition scope.

First validate acquisition and extraction using a small permitted real sample, then expand to the complete target course collection. Complete ingestion is not required for the first working end-to-end flow. Show current coverage; sample results do not imply complete collection. Real course and handbook artifacts have not yet been acquired.

## 2. Accounts, affiliation options and Profile

Registration is separate from academic Profile setup. Provide independent accounts, sign in/out and persistent server Profiles, without university credentials/student IDs or an assumption of university SSO. Email identifies a login, not university membership. Preserve the security and deletion requirements in [design §7](technical-design.md#7-security-privacy-and-state).

The sole affiliation/コース list, stable codes, types, parent relationships and availability flags are in [design §2](technical-design.md#2-affiliation-configuration-and-profile). These are **user-specified configuration requirements, not officially verified university facts**. Group options into undergraduate faculties and graduate schools/divisions. Only 情報学研究科 is selectable; all other affiliations remain visible, disabled and labeled **“Not yet available.”** Its seven required subdivisions are labeled **コース**, not 学科. Do not model other affiliations' internal subdivisions yet. Backend validation enforces availability and relationships even if a client bypasses disabled controls.

The required v0 Profile contains:

- Graduate school (情報学研究科).
- One configured コース under that school.
- Degree level (fixed to master's).
- Current year of study.
- Admission year.

Persist, retrieve and edit the Profile on the server so sign-in on another device restores it. Incomplete required fields block personalized discovery and direct the user to Profile setup. Load the latest Profile on each query; an edit affects the next query. Reject concurrent stale edits explicitly. Current year of study, admission year and searched academic year/semester are independent; do not derive or automatically advance them.

Service availability, course collection coverage and reviewed handbook coverage are independent. A selectable コース does not imply complete or verified academic coverage. Missing handbooks must not block ordinary discovery, but the application must state that curriculum alignment cannot yet be established.

## 3. One-query discovery behavior

Use a Japanese-first UI and support Japanese, Chinese and English query interpretation/reasons. The default target is explicitly shown as AY2026 second semester; never infer a different offering year from the Profile or current date.

Example: an Informatics master's student asks 「水曜日の午後の最初の授業で、試験がない科目を探しています。」 The server loads their Profile, resolves the timetable expression only with verified applicable period definitions, applies supported SQL conditions, then evaluates assessment evidence within the candidates. Historical undergraduate timetable evidence cannot resolve this request.

- Broad requests such as “recommend a course related to my field” can return Profile-based recommendations directly when the target semester is known. Do not demand a topic or arbitrary filters.
- Parse supported filters, semantic needs, **must** versus **prefer**, negation and unresolved requirements. A critical ambiguity returns a focused clarification request; the user resubmits a complete independent query. No v0 conversational memory is implied.
- Map validated supported conditions to parameterized SQL, never model-generated SQL. Day and period must refer to the same meeting; display all meetings. A strict all-meetings-fit request requires all meetings to be known and fit. Intensive/unknown schedules cannot establish a confirmed slot match.
- Treat academic relevance as a soft ranking signal unless a supported explicit hard condition applies. Do not automatically restrict the offering unit to the student's affiliation.
- Use Profile and reviewed applicable handbook scope before limiting candidates. Configurable inspection budgets favor speed; disclose partial inspection and do not claim exhaustive search or arbitrarily truncate before personalization.
- Detailed assessment evidence stays in source documents. Vector relevance is not proof of no exams. Separate or label insufficient evidence as **“Needs confirmation”**; never count it as satisfying a hard condition. An explicitly exam-requiring course is excluded entirely for a strict no-exam query.
- Missing/conflicting facts remain unknown. Missing/ambiguous handbook mappings permit discovery with a warning, without curriculum-category or graduation assurances. Do not guess applicable handbook editions.
- Distinguish an empty SQL candidate set from no matches found in the inspected batch. v0 has no show-more/expansion action, automatic expanded search or silent relaxation of must conditions.

Return at most five cards total, normally aiming for three to five when evidence supports them; fewer or zero is valid. Each card shows course name, all known schedule information, credits, syllabus link and one short recommendation reason. Basic facts come from database fields, with unknown values labeled. Reasons carry evidence references and uncertainty; “Needs confirmation” candidates are visibly separate from confirmed matches. Source excerpts are displayed only when permitted.

## 4. Data acquisition and curriculum limits

Use permitted syllabus data, official applicable handbooks and timetable/calendar evidence. [Source access](discovery/source-access.md) records historical observations and current verification tasks. Inspect rights, robots, format, scope and encoding before collection; do not bypass controls, enumerate guessed IDs or assume public reachability authorizes copying/embedding/redistribution.

Preserve raw sources separately from normalized metadata, documents, chunks and indexes, with original URLs, retrieval timestamps, source/section/version/year identifiers and missing/conflicting values. Capture only source-present fields, including multiple meetings. Do not invent data to fill a proposed schema. Validate a few real records and their extraction before full target ingestion; synthetic fixtures prove engineering only.

Curriculum in v0 means manually configured, human-reviewed handbook-applicability mappings from school/コース/degree/admission year to documents or sections. It does not mean implementing all category predicates, credit accounting, prerequisites or graduation rules. Preserve detailed academic provisions in documents and cite relevant evidence without certifying category fulfillment. Account/infrastructure work is independent of missing academic sources.

## 5. State, privacy and failure behavior

Persist accounts, authentication sessions and Profiles in PostgreSQL. Keep only the current query and results in frontend page memory; refreshing may discard them. No Redis, stored chat sessions, raw chat history, accumulated filters or server-side previously displayed result lists. v0.5 context sent by the browser must still be validated; identity and Profile remain server-authoritative.

Minimize transmission of academic context from the first parsing step. Exclude identifiers, email, passwords, tokens, full Profiles and earned-credit information from Python/provider payloads. Minimize free text or request reformulation when needed. Do not put personal content into logs, shared caches, embeddings or ingestion artifacts. Maintain owner isolation, parameterized queries, CSRF/Origin protection, password hashing and revocable sessions even for controlled prototypes.

Show clear onboarding, validation, clarification, coverage, empty, partial and provider-failure states. Provider failure must not fabricate cards, bypass authentication or relax hard conditions. Track the normal two generative LLM calls separately from embeddings, retries and document reading; measure token usage, latency and cost without recording raw queries or Profiles.

## 6. v0 acceptance and later gates

The first real-sample end-to-end milestone must demonstrate:

1. Register → sign in → save all required Profile fields → submit an independent query → source-backed cards or an honest empty/clarification response on the permitted sample. No full collection prerequisite.
2. Disabled affiliations are rejected server-side; invalid school/コース combinations fail. All configured options and seven コース display correctly, independently of coverage.
3. Two users cannot read/change one another's data. Profiles survive sign-out/sign-in, service restart and another-device sign-in; edits affect the next query. Missing fields, stale edits, logout/expiry, CSRF and provider failures behave as specified.
4. A fixed deterministic set exercises extraction/negation, schedule SQL, must/prefer, Profile-aware selection before truncation, missing handbooks, candidate boundaries, unknown assessment, evidence/IDs, partial/empty responses and provider failures. Hard failures never appear as uncertain matches.
5. Card facts equal database facts; references resolve to the correct offering/source revision. Structure/ID validation is not treated as proof of every recommendation claim. Human-check the sample's reasons and evidence.
6. Record actual tokens, retries, embeddings, latency, inspected/total candidate counts and sample coverage. Unmeasured quality/performance remains unverified.

Full v0 target coverage follows sample validation and requires a source-backed inventory/reconciliation for the agreed AY2026 second-semester collection. Local/controlled-prototype checks and public-release quality/security gates are separated in [design §9](technical-design.md#9-verification-and-release-gates). Follow-up/comparison tests belong to v0.5 and cannot block v0. Academic progress and complete rule-engine tests are outside v0.

Still unverified: permitted acquisition paths and formats, target inventory, official confirmation of configured names/relationships, applicable Informatics timetable and handbook versions, available assessment evidence, provider/model/data terms and measured costs/quality. These are blockers only for the work that depends on them.
