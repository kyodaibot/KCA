# Kyoto University Course Selection Agent — Requirements (Draft)

## 1. Product vision
Build a conversational course discovery assistant for Kyoto University students. The authenticated user maintains a student profile and describes course-year/term constraints and interests in Japanese, Chinese, or English. The assistant recommends courses, explains why they match, and links to authoritative syllabus and curriculum sources.

This is a course discovery and decision-support product, **not** an official registration system and **not** an authoritative graduation-eligibility checker.

## 2. Example user request
> 私は京都大学工学部地球工学科の3回生です。水曜日の午後の最初の授業で、試験がない科目を探しています。

Equivalent Chinese example: 我是工学部地球工学科大三学生，我想找到周三下午第一节、不需要考试的课。

Expected behavior:
1. Load the authenticated student profile; interpret institution, faculty, department and grade from that profile, and merge explicit conversation constraints without silently changing saved data.
2. Resolve the meaning of 「下午第一节」 against the university's actual timetable, or ask a focused clarification if ambiguous.
3. Filter offerings by academic year, term, weekday and period.
4. Distinguish explicitly confirmed no-exam grading from unclear or mixed grading descriptions.
5. Check department/year curriculum constraints only against verified rules and appropriate effective dates.
6. Return ranked, source-backed recommendations with explicit uncertainty.

## 3. Users and use cases
- Undergraduate students browsing courses by schedule, assessment, credits and interests.
- Students checking whether a course could count toward a specific curriculum category (with appropriate qualifications).
- Students asking follow-up questions and comparing courses conversationally.

## 4. Authoritative data sources and provenance
- Public Kyoto University syllabus: https://www.k.kyoto-u.ac.jp/external/open_syllabus/top
- Official faculty/department curriculum guides, handbooks, course tables and graduation requirements, scoped by **admission cohort** and **academic year** when applicable.
- Official academic calendar and term/period information if needed.

Every imported record must retain source URL, source document/page ID if available, extraction time, academic year, update time if available, and raw source excerpt. A local copy does not replace the source of truth.

### Source feasibility verification (before implementation)
- Inspect pages and interaction patterns to determine whether public HTML, a documented API, or JavaScript-rendered navigation is involved.
- Check robots.txt, website terms, access limits, and licensing/redistribution constraints. Do not bypass access controls, CAPTCHAs or authentication.
- Establish a respectful crawl rate and incremental refresh strategy; prefer official export/API permission where available.
- Run a small proof of concept with a few pages, preserving Japanese encoding and layout before choosing a full crawler.

## 5. Functional requirements
### Course ingestion
- Retrieve permitted public syllabus details; robustly capture structured data and free text.
- Normalize to a canonical course/section/offering schema; store original HTML when lawful/useful and Markdown for retrieval.
- Stable IDs, upsert, deduplication, academic-year versioning, incremental sync, deletion handling and audit logs.
- Parse multiple class meetings, irregular/intensive courses and different campus locations if provided.

### Example course schema (verify availability from source)
- course_id, offering_id, academic_year, academic_term
- title_ja/title_en, faculty, department, instructor
- weekday, period, location, course_type, eligible_years, credits
- syllabus text, objectives, prerequisites, assessment original text
- assessment_type and `exam_required` (`yes` / `no` / `unknown`), extraction confidence and citation
- source_url, source_updated_at, fetched_at

Never assume a field is published just because it appears in this proposed schema. Keep unknown values rather than hallucinating.

### Curriculum and graduation rules
- Ingest and version official curriculum documents by faculty/department, admission cohort and effective year.
- Model clear rules as structured predicates and mappings (required/elective, category, credit minimum, level constraints, prerequisites, exclusions).
- Use RAG for retrieving/explaining relevant source passages, **not** as the sole deterministic rule engine.
- Distinguish `eligible`, `ineligible`, and `unverified` with rule/source explanations.
- Do not claim a student has completed requirements without their actual course history and relevant official rule verification.

### Search and recommendations
- LLM converts natural-language requests into a typed search intent; show or internally validate interpretation.
- Hard filters: term, day/period, level/affiliation restrictions, confirmed exam absence if explicitly required.
- Semantic retrieval: topic, interests, teaching content, related skills; start with vector retrieval and compare with hybrid keyword retrieval when justified.
- Ranking: hard constraint satisfaction first, then curriculum relevance, semantic relevance and preferences.
- Every suggestion includes course name, time, grading method, curriculum status, reasons, syllabus URL, evidence and warnings.
- Unknown fields must not be treated as definite matches for strict requests.

### Conversational UX
- Follow-up questions preserve stated constraints unless updated.
- Clarify only essential ambiguities (e.g. course period, year, missing cohort for graduation rule).
- Users can adjust filters, compare candidates and inspect source excerpts.
- Japanese-first UI, with optional Chinese/English responses.

### User authentication (core MVP)
- Provide registration, login and logout with independent accounts and persistent profiles; do not assume access to Kyoto University internal authentication or require university credentials/student IDs.
- MVP architecture choice: email/password authentication in the existing backend, Argon2id password hashes through a maintained library, and opaque, revocable PostgreSQL-backed sessions in an HttpOnly cookie. Do not implement cryptographic primitives or store plaintext passwords/browser-readable bearer credentials.
- Use HTTPS and Secure/SameSite cookies in hosted environments, CSRF protection for state-changing requests, bounded login attempts and generic authentication failures. Logout revokes the server session immediately.
- Email is a login identifier, not proof of university membership. Verification and password recovery must use expiring, single-use tokens and a configured mail adapter before a public launch; local development uses a fake adapter and never sends mail automatically.

### Student profile and onboarding (core MVP)
- After registration, guide the user through a short Japanese-first onboarding flow. Require faculty, department, current academic year (grade), and enrollment year before personalized search or recommendations.
- Grade is year of study, enrollment year is admission cohort, and search academic year is the year of an offering; keep these independent and never infer one from another.
- Optional fields: academic track/specialization; total earned credits; major/specialized credits; general education credits; additional curriculum-specific credit categories. Optional credits can be skipped and added later; null means not reported, distinct from zero.
- Users can save the total, major/specialized and general education credit summaries even when curriculum rules or category mappings are unavailable. These are self-reported summaries with provenance and report timestamps, separate from version-scoped curriculum category reports and officially validated progress. Missing rules must not block summary entry or imply a verified category mapping.
- Persist profiles in PostgreSQL; support viewing, editing and updating academic progress over time. Do not automatically increment grade or infer a track.
- Model faculty/department relationships, optional tracks and versioned curriculum credit categories. Storage must support multiple faculties/departments/cohorts even while the course-data pilot remains limited to one scope. A saved affiliation does not establish source coverage or valid curriculum rules.
- Validate affiliation/track combinations and nonnegative decimal credits. Categories can overlap; do not force subtotals to sum to total credits unless a verified curriculum rule explicitly defines that relationship. Label all manually entered credit values as self-reported, with an update timestamp.
- Ensure ownership on the server for every profile operation. The authenticated account identifies the owner; ignore/reject client-supplied owner IDs. Users can access only their own profile.
- Offer both personal-profile deletion (returning the account to incomplete onboarding) and full account/personal-data deletion with recent password reauthentication. Remove credit data, revoke sessions on account deletion, clear client state and explain backup retention. Public course/catalog data is unaffected.

### Authenticated chatbot integration (core MVP)
- At conversation start, load the authenticated user's current profile on the backend; refresh it for every personalized request. Stored affiliation, grade, enrollment year and track are authoritative defaults, not browser-supplied ownership claims.
- Block personalized requests until required onboarding fields are complete, with a clear onboarding action; anonymous access may expose public catalog/scope information only.
- Resolve applicable curriculum rules from faculty, department, enrollment year, optional track, and the requested offering academic year. Missing essential track/version evidence or overlapping rules yields `unverified`; do not guess a track/version.
- Use profile and self-reported progress to personalize recommendations and explain relevant categories. Do not ask again for stored information. Ask only for missing material details, including search year/term or an essential optional track.
- Conversation-specific what-if changes must be explicit and separate from saved profile edits; changes to the profile invalidate prior profile-derived context. Preserve unrelated schedule/topic constraints.
- Both personalized search and chat accept an optional `temporary_profile_override` object containing academic `values` and an explicit client `confirmed` flag, bound to the request's current `profile_revision`. Divergent academic intent without confirmation prompts clarification; the UI resubmits the exact reviewed values with `confirmed=true`. The model may propose values but cannot supply confirmation. Overrides apply only to the request and explicitly continued hypothetical conversation, are labeled hypothetical, and never write the saved profile. Only authenticated, owner-scoped profile-management endpoints can persist profile changes.
- Credit totals do not establish course completion, prerequisites, exclusions, category eligibility or graduation completion. Without actual course history and verified applicable rules, completion remains unknown; self-reported totals never certify graduation eligibility.
- Persist accounts/profiles, not raw conversations. Keep identity, secrets and personal profile content out of logs, analytics, shared caches and embeddings. LLM academic-context allowlist: faculty/department codes or names, grade, enrollment cohort, optional track, and requested offering academic year/term, each only when needed. Credits are excluded by default; for an explicit credit/progress question, send only the requested total/major/general summary or relevant category amount and label it `self_reported`/unverified. Never send account/student identifiers, email addresses, passwords, tokens, unnecessary personal information or the full profile object. Project academic fields on the server; remove identity/secrets and unrequested credit information from free-text and prior context before provider calls, and ask for a reformulation if a safe minimal request cannot be formed. Hypothetical context uses the same allowlist and carries a hypothetical label.

## 6. Nonfunctional requirements
- Accuracy: source-backed, strict filtering, no unsupported graduation assurances.
- Testability: extraction fixtures, parser unit tests, structured-query tests and rule evaluation tests.
- Performance: define realistic p95 latency target during planning and measure it.
- Security/privacy: accounts and student profiles are required in MVP; no student ID, transcript or detailed course history is required. Isolate users, minimize personal data, support deletion, retain no raw chats by default, and apply secret management and least privilege. Encrypted backups expire within 30 days; restores must reapply deletion records before serving traffic. Deletion records contain only opaque account IDs and expire after the backup window.
- Cost: modest hosting, local development, cloud-agnostic where reasonable.
- Observability: sync failures, extraction field coverage, search results, ranking latency, user feedback.

## 7. Suggested MVP scope
1. Choose one faculty/department and one academic year for a small pilot.
2. Collect a limited permitted syllabus sample and one official curriculum document.
3. Implement course schedule filtering and exam assessment classification with an `unknown` state.
4. Add registration/login/logout, persistent student profiles, required onboarding, optional self-reported credit progress, profile editing and personal-data deletion.
5. Build basic semantic retrieval, authenticated chat and source-linked personalized recommendations using saved profiles.
6. Validate using hand-labeled examples and cross-user security/onboarding tests before expanding scope.

**Out of MVP**: campus-wide ingestion, automated course registration, unofficial credit certification, personalized completed-credit audits, multi-agent orchestration, Kubernetes.

## 8. Suggested technology candidates — evaluate, don't blindly adopt
- Backend: Python + FastAPI + Pydantic.
- Importer: requests/httpx + BeautifulSoup/selectolax; Playwright only where needed and permitted.
- Storage: PostgreSQL for structured data and curriculum rules, pgvector if embedding retrieval is justified.
- Raw artifacts: Markdown plus source metadata in local files or object storage.
- Embeddings: compare multilingual models suited to Japanese syllabi; begin with one low-cost baseline.
- Frontend: Next.js/React or a minimal prototype UI based on development effort.
- Testing: pytest and a small golden dataset, frontend E2E later.
- Deployment: local Docker Compose initially; cloud hosting after pilot works.

Ask the planner to compare alternatives, costs, operations, and trade-offs, not treat this list as mandatory.

## 9. Evaluation
- Parsing: required-field coverage, assessment extraction precision and unknown rate on human-labeled samples.
- Filtering: precision/recall for weekday, period, term and no-exam constraints.
- Curriculum: rule-based eligibility accuracy, citation accuracy, unknown/ambiguous cases.
- Retrieval: Hit@k / Recall@k with page/course IDs; compare SQL filtering + vector vs SQL + hybrid + rerank.
- Answer quality: groundedness, citation correctness, constraints satisfied and conversational follow-up behavior.
- Prepare 30–50 realistic questions with gold courses or explicitly no result.

## 10. Open decisions for discovery
- Public site's legal/technical crawling feasibility; stable identifiers; how exam modes are expressed.
- Exact academic period definitions; whether 「Wednesday first period in the afternoon」 has an unambiguous mapping.
- Department/cohort-specific policy documents and machine-readable availability.
- Which syllabus fields are missing/ambiguous and therefore require `unknown`.
- Single vs multiple meeting slots, timetable changes and historical offerings.
- Hosting budget and language priorities.

## 11. Authentication/profile acceptance criteria and scope changes

- Two independently registered users retain distinct profiles after logout/login and service restart; neither can read or change the other's data via request tampering, cached responses or reused conversation state.
- Registration leads to onboarding; missing any required field blocks personalized recommendations. Omitting all optional fields still permits profile completion and ordinary discovery; missing rule-critical track may keep curriculum status unverified.
- Editing affiliation/cohort/track or credits changes the next recommendation context without re-registration or repeated questions. Concurrent edits return an explicit revision conflict rather than silently overwriting progress.
- Without any reviewed curriculum rules, all three self-reported credit summaries can be saved/reloaded/edited, including null and zero; none establishes verified progress. A confirmed temporary override changes recommendations without changing saved data; unconfirmed/model-generated confirmation cannot authorize an override or a persistent write.
- Negative/invalid credits and mismatched faculty/department/track IDs are rejected. Blank optional credits persist as null; zero persists as zero. Curriculum-specific categories remain version-scoped across profile changes.
- Logout/expiry invalidate access; state-changing requests without valid CSRF protection fail. Profile/account deletion requires recent reauthentication, removes owned data as specified, and cannot delete another user's records or public catalog data.
- A profile with apparently sufficient aggregate credits never receives graduation certification or assumed prerequisite completion. Unsupported departments/cohorts remain clearly marked as lacking curriculum coverage.
- Existing source acquisition restrictions, multilingual behavior, strict filters, citations, unknown states, benchmark and cost constraints remain intact.

Necessary changes to earlier scope: persistent accounts/profiles now supersede the previous plan's deferral of accounts and entirely client-held student context. Optional self-reported credit summaries are in MVP; personalized completed-credit audits and official graduation certification remain out of MVP. Online serving needs narrowly scoped writes for authentication/profile tables while retaining read-only access to academic catalog/rules. The frontend choice remains the existing React/TypeScript/Vite plan; authentication alone does not require Next.js or a new service.
