# V1 Stabilization — Final Design

**Date:** 2026-06-06
**Status:** Approved
**Scope:** Fix all known defects, unify migrations, productionize Docker, add tests, add API docs, add CI

---

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Reminder system | Keep virtual/computed model | Works for V1; no background task needed |
| Migration strategy | Remove SQL files, use GORM AutoMigrate only | Simpler, no extra dependency |
| Testing scope | Backend core unit tests only | V1-appropriate coverage |
| CI/CD | GitHub Actions (go build + go test + npm build) | Automated quality gate |
| API documentation | swaggo/swag with Swagger UI | Industry standard, auto-generated |

---

## Work Package 1: Defect Fixes

### 1.1 Frontend Store data access bug

**File:** `frontend/src/utils/request.ts`
**Problem:** Axios response interceptor returns `response.data` (unwraps `{code, message, data}` envelope). But stores access `res.data` again, yielding `undefined`.

**Fix:** Update all store files (`stores/asset.ts`, `stores/subscription.ts`, `stores/category.ts`, `stores/dashboard.ts`) to use `res` directly instead of `res.data`. The interceptor already returns the inner `data` field.

**Verification:** Load dashboard and list pages in browser; confirm data renders correctly.

### 1.2 i18n duplicate key

**Files:** `frontend/src/locales/zh-CN.ts`, `frontend/src/locales/en-US.ts`
**Problem:** Key `subscription.cycle` is defined twice — once as "Billing Cycle" / "计费周期" and once as "Custom" / "自定义". The second overwrites the first.

**Fix:** Rename to `subscription.billingCycle` (label) and `subscription.custom` (custom cycle option). Update all template references.

### 1.3 Missing 404 route

**File:** `frontend/src/router/index.ts`
**Problem:** Undefined paths show a blank page.

**Fix:** Add a catch-all route `{ path: '/:pathMatch(.*)*', redirect: '/' }` at the end of the children array.

### 1.4 Hard redirect on 401

**File:** `frontend/src/utils/request.ts`
**Problem:** `window.location.href = '/login'` causes full page reload.

**Fix:** Import router and use `router.push('/login')`.

### 1.5 Remove unused Reminder model

**Files:** `backend/internal/model/reminder.go`, `backend/pkg/database/database.go`
**Problem:** `Reminder` model and `reminders` table exist but are never used. The service computes reminders virtually from subscriptions/assets.

**Fix:** Delete `model/reminder.go`. Remove `&model.Reminder{}` from `AutoMigrate()`. The `reminders` table will remain in existing databases but won't be created for new ones.

---

## Work Package 2: Migration Unification

### 2.1 Remove SQL migration files

**Action:** Delete `backend/migrations/` directory entirely.

### 2.2 Verify AutoMigrate coverage

Confirm all models are registered in `database.AutoMigrate()`:
- `User`, `Category`, `Subscription`, `Asset`, `ReminderRead`, `ReminderDismissed`
- After removing `Reminder`, the list has 6 models.

### 2.3 Verify soft delete

Confirm `model.Model` embeds `gorm.DeletedAt` and all business models (User, Category, Subscription, Asset) embed `Model`.

---

## Work Package 3: Docker Compose Production Readiness

### 3.1 docker-compose.yml improvements

- Add `restart: unless-stopped` to all three services
- Add health check to backend: `curl -f http://localhost:8080/api/v1/auth/me || exit 1` (or a dedicated `/health` endpoint)
- Add health check to frontend: `curl -f http://localhost:80/ || exit 1`
- Make `frontend` depend on `backend` with `condition: service_healthy`
- Replace hardcoded credentials with `${VAR:-default}` pattern
- Add logging config with size limits

### 3.2 Add .dockerignore files

Create `backend/.dockerignore` and `frontend/.dockerignore` to exclude:
- `.git`, `node_modules`, `dist`, `*.md`, `.env`

### 3.3 nginx improvements

**File:** `frontend/nginx.conf`
- Add gzip compression for text/css/js/json/svg
- Add security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `X-XSS-Protection: 1; mode=block`
- Add cache headers for static assets (`.js`, `.css`, `.svg`, `.png` — 7 days)

### 3.4 Add /health endpoint

**File:** `backend/internal/router/router.go`
- Add unauthenticated `GET /api/v1/health` returning `{"code": 0, "message": "ok"}`

---

## Work Package 4: Swagger API Docs + Error Codes

### 4.1 Swagger integration

- Add `swaggo/swag` dependency
- Add Swagger annotations to all API handlers (auth, user, subscription, asset, category, reminder, dashboard)
- Register Swagger UI route at `/swagger/*`
- Generate docs via `swag init` (checked into repo)

### 4.2 Unified error codes

**New file:** `backend/internal/service/error_codes.go`
- Define constants for all error codes currently hardcoded as literals
- Pattern: `ErrInvalidCredentials`, `ErrForbidden`, `ErrNotFound`, etc.
- Update all service files to use constants

Existing error codes:
| Constant | Code | HTTP | Description |
|----------|------|------|-------------|
| ErrInvalidCredentials | 40101 | 401 | Invalid credentials |
| ErrInvalidReminderID | 40001 | 400 | Invalid reminder ID |
| ErrIncorrectPassword | 40002 | 400 | Incorrect old password |
| ErrForbidden | 40301 | 403 | Forbidden |
| ErrNotFound | 40400 | 404 | Resource not found |
| ErrDuplicateUsername | 40901 | 409 | Username already exists |
| ErrDuplicateEmail | 40902 | 409 | Email already exists |
| ErrDuplicateCategory | 40903 | 409 | Category name already exists |

---

## Work Package 5: Backend Tests + GitHub Actions CI

### 5.1 Test infrastructure

- Use SQLite in-memory database for tests (GORM supports it natively)
- Create `backend/internal/testutil/testutil.go` with helper functions:
  - `SetupTestDB()` — initializes in-memory SQLite + AutoMigrate
  - `CreateTestUser(db, username)` — creates a user and returns user + token
  - `NewTestRouter(services)` — sets up a test Gin engine

### 5.2 Test coverage

| Test File | Covers |
|-----------|--------|
| `service/auth_test.go` | Register, Login, GetCurrentUser, duplicate username/email |
| `service/subscription_test.go` | CRUD, filtering by status/category, next payment date calculation |
| `service/asset_test.go` | CRUD, filtering by type/status, expiring soon |
| `service/category_test.go` | CRUD, default category initialization |
| `service/reminder_test.go` | List computed reminders, mark read, dismiss |
| `service/dashboard_test.go` | Monthly/yearly spend, upcoming renewals, expiring assets |
| `api/middleware_test.go` | JWT auth middleware: valid token, expired token, no token |

### 5.3 Data isolation tests

Each service test includes cross-user isolation:
- User A cannot read/update/delete User B's resources
- All list queries return only the authenticated user's data

### 5.4 GitHub Actions workflow

**File:** `.github/workflows/ci.yml`
- Triggers: push to `master`, pull_request to `master`
- Jobs:
  1. `backend`: `go build ./...`, `go test ./... -race`
  2. `frontend`: `npm ci`, `npm run build`
- Go module caching via `actions/cache`
- Node module caching via `actions/cache`

---

## Acceptance Criteria

- [ ] `go build ./...` passes
- [ ] `go test ./... -race` passes (all new tests green)
- [ ] `npm run build` passes (frontend type-check + build)
- [ ] `docker compose up` starts all services successfully
- [ ] Swagger UI accessible at `/swagger/index.html`
- [ ] GitHub Actions CI passes on master
- [ ] No hardcoded SQL migration files remain
- [ ] No known frontend data-access bugs
- [ ] README updated with new features (Swagger URL, health endpoint)
