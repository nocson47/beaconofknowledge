# Beacon of Knowledge — Mini Forum (MVP)

A compact forum application (MVP) implemented with a Vite + React frontend and a Go (Fiber) backend. The project uses a hexagonal architecture: core business logic lives in internal usecases and concrete behavior is provided by adapters for HTTP, persistence, email, and other infrastructure.

This README explains how to run the project locally, the implemented features, security considerations, and suggested next steps.

## Tech stack and architecture
- Backend: Go (Fiber)
- Database: PostgreSQL (primary), Redis (cache), optional MongoDB (reports/audit)
- Frontend: Vite + React + TypeScript + Tailwind CSS
- Architecture: Hexagonal (usecases, repositories, adapters)

## Features (high level)
- Authentication: register and login (bcrypt password hashing) with JWT-based sessions
- Profiles: avatar, bio, social links, profile update (owner or admin only)
- Threads & replies: create, read, update, delete (soft-delete semantics in usecases)
- Admin: role-based middleware and simple admin UI pages
- Reporting: user report flow with optional Mongo persistence for audit logs
- Caching: Redis used where appropriate
- Rate limiting and basic request throttling middleware

## Quick start (local development)
Prerequisites: Go 1.20+, Node 18+, Docker (optional for DB/Redis/Mongo)

1. Copy example environment file and edit values

```bash
cp backend/.env.example backend/.env
# Edit backend/.env to set DB, Redis, and other values
```

2. (Optional) Start local services with docker-compose

```bash
docker-compose up -d postgres redis mongo
```

3. Start the backend

```bash
cd backend
go build ./cmd && ./cmd
# or for a quick dev run: go run ./cmd
```

4. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

Defaults: backend listens on port 3000 (see `cmd/main.go`), frontend dev server runs on port 5173.

## Configuration highlights
- `backend/.env` contains runtime configuration for DB, Redis, Mongo, and SMTP.
- For development, the SMTP settings are intentionally left blank so the app prints reset links to logs (ConsoleEmailSender). To enable real SMTP, set `SMTP_HOST`, `SMTP_USER`, `SMTP_PASS`, `SMTP_PORT`, and `SMTP_FROM` and restart backend.

Security note: never commit real API keys or secrets to the repository. Use `.gitignore` and a secrets manager for production credentials.

## Security considerations (current status)
- Authentication: uses JWTs sent in the `Authorization` header. This avoids cookie-based CSRF but storing tokens in `localStorage` is vulnerable to XSS. Mitigate with strong Content Security Policy (CSP) and input sanitization.
- CORS: development is permissive to `http://localhost:5173`. In production, configure allowed origins via environment variables and avoid overly permissive CORS.
- CSRF: there is no CSRF token system implemented. If you switch to cookie-based authentication, add CSRF protections.
- Recommended hardening:
  - Restrict allowed origins from env and validate `Origin`/`Referer` on state-changing requests.
  - Add CSP and secure response headers (X-Content-Type-Options, X-Frame-Options, Referrer-Policy, etc.).
  - Use a secrets manager for API keys and credentials; never commit `.env` with secrets.

## Limitations & TODO
- Improve automated tests: add unit and integration tests for core usecases and HTTP handlers
- Production readiness: database migrations, robust logging, observability, and secrets management
- Mongo audit logging currently uses non-blocking writes; consider retries and indexing for production
- More admin UI polish and functionality (user / report management)

## Run & test password reset locally (development flow)
1. Ensure `backend/.env` has no SMTP credentials so ConsoleEmailSender is active
2. Trigger "forgot password" from the frontend or curl the API; the backend will log a ResetURL with a token
3. Open the logged ResetURL in your browser to reset the password

If you want the app to actually send email, configure SendGrid (or other SMTP provider):
- Configure `SMTP_HOST`, `SMTP_PORT` (587 for STARTTLS recommended), `SMTP_USER` (SendGrid uses `apikey`), and `SMTP_PASS` (API key)
- Verify the `SMTP_FROM` address in SendGrid (Single Sender or Domain Authentication) to avoid 550 errors

## Troubleshooting notes
- If emails are rejected with a 550 from SendGrid, verify the `From` address is a verified sender identity in your SendGrid account.
- If a secret was accidentally committed, rotate/revoke the key immediately and remove the file from Git tracking.

## Suggested next steps I can implement
- Enforce origin/referrer validation middleware for state-changing endpoints and read allowed origins from env
- Make CORS configuration environment-driven (`CORS_ALLOWED_ORIGINS`) and toggle `allowCredentials`
- Add basic CSP and security headers middleware
- Add a focused unit test for `UserHandler` or a core usecase

Tell me which of the next steps you'd like me to implement first and I'll add it to the todo list and apply the change.

---
Generated by developer assistance. Small edits or further expansion are available on request.
