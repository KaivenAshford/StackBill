# StackBill V1 Stabilization Design

**Date:** 2026-06-01
**Status:** Approved

## Goal

Make the current PC browser version credible as a V1 release candidate. The work focuses on build confidence, core backend tests, security boundaries, and delivery documentation.

This stabilization pass does not add business features, redesign UI, or include mobile work.

## Scope

The stabilization scope includes:

- Establish quality gates for backend and frontend builds.
- Add focused backend tests for authentication, user data isolation, core CRUD behavior, reminders, and dashboard aggregation.
- Check high-risk paths: registration, login, JWT protection, and cross-user access.
- Fix blocking bugs found while building or testing.
- Update delivery documentation so a new developer can start, validate, and understand V1 limits.
- Keep mobile browser, PWA, and native app work out of the current plan until explicitly requested.

## Acceptance Criteria

- `go build ./...` passes in `backend`.
- `npm run build` passes in `frontend`.
- New backend tests pass for the selected core service, repository, and API/security paths.
- Tests cover user isolation for categories, subscriptions, assets, and reminders where those paths are exercised.
- README documents local validation commands, Docker startup checks, and known V1 limits.
- Project documentation states that mobile is not in current scope.

## Testing Strategy

Backend tests are the priority because the highest V1 risks are data isolation, authentication, and business rules. Frontend verification uses the existing build and TypeScript pipeline first; no new frontend testing framework is introduced in this pass.

### Service Tests

Add focused tests for:

- Auth: registration, duplicate username/email handling, password verification, and default category creation.
- Subscription: create/update behavior, next payment date calculation, user-scoped reads, and delete ownership checks.
- Asset: create/update behavior, list filtering, and user-scoped reads.
- Reminder: generated renewal and expiration reminders, mark read, mark all read, and dismiss/delete behavior.
- Dashboard: monthly and yearly expense totals, category expense aggregation, upcoming renewal counts, and expiring asset counts.

### Repository Integration Tests

Use a minimal GORM-backed test setup for critical query behavior:

- Queries include `user_id` filters.
- Pagination and filters return expected records.
- Count and aggregate methods do not leak data across users.

The goal is not exhaustive CRUD coverage. The goal is to lock down the data boundaries and calculations most likely to break V1.

### API And Auth Smoke Tests

Use Gin test utilities for a small set of route-level checks:

- Protected routes reject unauthenticated requests.
- Authenticated users can access their own records.
- Cross-user access returns the existing not-found or forbidden-style error used by the codebase.

## File Boundaries

- Add test files next to the code under test, such as `backend/internal/service/subscription_test.go`.
- If shared setup is needed, add a small `backend/internal/testutil` package.
- Keep API response shapes unchanged unless tests expose a real bug.
- Do not change frontend page structure as part of this pass.
- Frontend changes are limited to build/type fixes if the current build fails.

## Implementation Order

1. Run baseline verification:
   - `cd backend && go build ./...`
   - `cd frontend && npm run build`
2. Fix only blocking build failures.
3. Add minimal backend test infrastructure.
4. Add tests for authentication and user isolation first.
5. Add tests for subscriptions, assets, reminders, and dashboard calculations.
6. Fix bugs revealed by those tests without broad refactoring.
7. Update README and planning documentation.
8. Run final verification commands.

## Risk Handling

If pure unit tests would require large architecture changes, prefer lightweight integration tests over introducing mocks or interfaces only for testing.

If a local PostgreSQL test database is unavailable, database-backed tests may be gated behind an explicit environment variable while ordinary builds remain unblocked. The test documentation must state how to enable them.

If frontend build failure is caused by environment or dependency setup rather than code, document the required toolchain and do not mask the issue with unrelated changes.

If verification exposes a larger design issue, record it as a follow-up unless it blocks V1 stability.

## Out Of Scope

- Mobile browser experience.
- PWA.
- Native mobile app.
- New business features.
- UI redesign or polish work.
- Full end-to-end browser test framework.
- Broad architecture refactors.
- Exhaustive edge-case testing for every request parameter.

## Expected Outcome

After this pass, StackBill should have a reliable baseline for a PC browser V1 release candidate: both builds pass, core backend behavior is protected by focused tests, user isolation has explicit coverage, and README explains how to validate the application.
