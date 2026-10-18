# KCA

Kyoto Course Agent is a three-person student project for course discovery. v0 targets Kyoto University Graduate School of Informatics master's students and AY2026 second semester: register/sign in, save a Profile, then submit one independent query for up to 3–5 source-linked course cards.

## Current status

The repository contains a Go/Gin module, an empty router and a health-endpoint test. The health route and runnable Go entry point are not yet implemented. Existing Docker/Compose scaffolding still references the previous Python backend; PostgreSQL/pgvector configuration and its named volume are present. Frontend and `rag/` are placeholders. These are existing artifacts, not verified startup success or completed account/recommendation features. No real course or handbook artifacts have been acquired.

## Target architecture

- `frontend/`: React + TypeScript + Vite.
- `backend/`: Go + Gin, pgx + sqlc; Auth, Profile, Course, minimal Curriculum and Chat orchestration.
- `rag/`: Python internal online HTTP API for parsing/retrieval/ranking plus separate offline ingestion jobs.
- PostgreSQL + pgvector with distinct personal-data, retrieval-read and offline-write privileges; explicit Go-owned migrations and local Docker Compose.

Start with a permitted real-data sample; expand target coverage after the first working flow. v0.5 adds follow-ups/explanations/comparisons. Progress, expanded search, plans and timetables come later.

## Documentation

- [Agent rules](AGENTS.md) and [requirements](docs/requirements.md).
- [Technical design](docs/technical-design.md), including authoritative affiliation configuration, data permissions and internal interfaces.
- [Implementation roadmap and current progress](plans/implementation-plan.md).
- [Planning constraints](docs/planning-prompt.md) and [interactive learning workflow](plans/learning-workflow.md).
- [Historical source observations and pending Informatics checks](docs/discovery/source-access.md).

Startup commands will be documented as their implementation tasks are verified. Do not reset the existing database volume to adapt the scaffold.
