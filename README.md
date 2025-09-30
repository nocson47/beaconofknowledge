# Beacon of Knowledge — Mini Forum (MVP)

This repository is a small Vite + React frontend and Go (Fiber) backend implementing a mini forum with hexagonal architecture.

The purpose of this README is to explain how to run the project locally, list implemented features, outline security notes and limitations, and provide a candid scoring against the assignment rubric shown in the attachment.

## Quick overview
- Backend: Go + Fiber, PostgreSQL, Redis, optional Mongo for reports/logs
- Frontend: Vite + React + TypeScript + Tailwind
- Architecture: Hexagonal — internal usecases define ports; adapters implement persistence, HTTP, JWT

## Implemented features (high level)
- Auth: register, login (bcrypt hashed passwords), JWT-based authentication
- Profiles: avatar URL, bio, social link, update profile (owner or admin only)
- Threads & replies: create/read/update/delete (soft-delete conventions in usecases)
- Admin: role-based middleware (admin-only endpoints), simple admin pages in frontend
- Reports: user reporting flow implemented; Mongo adapter persists reports and audit logs when available
- Caching: Redis wired in and used for some handlers
- Rate limiting: per-IP and per-user rate limiters in middleware

## What to run locally
Prerequisites: Go 1.20+, Node 18+, Docker (for running Postgres/Redis/Mongo if you prefer using docker-compose)

1) Copy env example and adjust values

```bash
cp backend/.env.example backend/.env
```

2) Start dependent services (quick dev using docker-compose)

```bash
docker-compose up -d postgres redis mongo
```

3) Start backend

```bash
cd backend
go build ./cmd && ./cmd
```

4) Start frontend

```bash
cd frontend
npm install
npm run dev
```

API runs on :3000 by default (check `cmd/main.go`). Frontend dev server runs on :5173.

## Security notes (current status)
- Authentication uses JWTs stored in `localStorage` and sent as `Authorization: Bearer <token>`.
  - Pros: This pattern avoids classic cookie-based CSRF if cookies are not used.
  - Cons: Storing tokens in localStorage is vulnerable to XSS; mitigate with CSP and careful input sanitization.
- CORS: Development CORS is permissive to `http://localhost:5173`. In production, you should configure allowed origins via environment variables and avoid `Access-Control-Allow-Origin: *` with credentials.
- CSRF: There is no explicit CSRF token implemented. Because JWTs are sent in headers by the frontend, classic CSRF risk is low — provided you do not switch to cookie-based auth without adding CSRF protection.
- Recommended immediate hardening:
  - Read allowed origins from environment and restrict in production.
  - Add middleware to validate `Origin`/`Referer` for state-changing requests.
  - Add CSP and secure headers (X-Content-Type-Options, X-Frame-Options) and sanitize inputs.

## Limitations / TODO (short)
- Tests: limited automated tests; add unit/integration tests for usecases and handlers
- Mongo audit logging currently performs non-blocking writes but should log failures and create appropriate indexes
- UpdateReportStatus behavior should be tested for parity across Postgres and Mongo adapters
- Additional production hardening (migrations, secrets management, rate limiting tuning)

## Rubric assessment (honest self-score)

I evaluated the project against the rubric in the attachment. Below is a candid score out of 100 and how the points were assigned.

- Database design / ERD, indices and relationships — 18/20
  - Reason: Schema and `init.sql` present; Postgres + Mongo used properly. Missing explicit secondary DB justification docs and some index creation for Mongo audit collection.
- Features (Auth, profiles, posts/comments, admin) — 18/20
  - Reason: Core features implemented (register/login/profile/threads/replies/report). Admin flows and role checks are implemented.
- Admin UI & features (manage users/categories/reports) — 16/20
  - Reason: Admin UI pages exist but are rudimentary and need more polish and CRUD completeness.
- Secondary datastore usage & justification — 4/5
  - Reason: Redis + Mongo used. Justification present but could be documented more clearly.
- Security basics (hash, RBAC, validation, CSRF, XSS, rate limit, file policy) — 10/15
  - Reason: Password hashing, RBAC, validation, rate limiting are present. CSRF token and CSP/XSS protections are incomplete.
- Code quality, README, documentation, and demos — 16/20
  - Reason: Code structured using hexagonal architecture and has comments; README is now added. More docs and automated tests would improve score.

Total: 82 / 100

Notes on scoring: This is a pragmatic grading for a mini-project: implementation is solid for an MVP and demonstrates architectural understanding (hexagonal, adapters, usecases). The main deductions come from missing tests, missing CSRF/CSP hardening, and admin UI polish.

## Next steps I can do for you (pick any)
- Add middleware to validate `Origin`/`Referer` for POST/PUT/DELETE and wire allowed origins via `.env` (recommended)
- Change `Cors()` to read allowed origins and `allowCredentials` from env
- Add basic CSP and security headers middleware
- Add a small unit test for `UserHandler` or `Report` usecase

If you tell me which of the above to implement first I will update the todo list and make the change.

---
Generated by developer assistance — updates and small fixes can be made on request.
