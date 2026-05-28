# Session B - Backend Tasks

## Owner

- Backend Engineer (Go API + SQLite)

## Objective

- Deliver stable API and data layer for auth, profile, pantry, meal plans, macros, and consumption flows.

## Ordered Task List

### 1) Schema and setup first (0:30-1:30)

- [ ] `P0` Implement SQLite schema + migrations for:
  - users
  - user_body_metrics
  - user_preferences
  - budgets
  - purchases
  - ingredients
  - recipes
  - recipe_ingredients
  - pantry_items
  - meal_plans
  - plan_meals
  - consumption_log
  - chat_messages
  - favorites (optional table behind feature flag/time)
- [ ] `P0` Add one-command local DB setup/reset path.
- [ ] `P0` Seed ingredients and minimal recipes for deterministic demo.

### 2) Auth and profile APIs (1:30-2:30)

- [ ] `P0` Auth endpoints: `POST /auth/register`, `POST /auth/login`, `GET /me`.
- [ ] `P0` Password hashing and token/session expiry rules.
- [ ] `P0` Profile endpoints for metrics, preferences, budget:
  - `GET /profile`
  - `PATCH /profile/metrics`
  - `PATCH /profile/preferences`
  - `PATCH /profile/budget`
- [ ] `P0` Publish API v1 examples for frontend and AI sessions.

### 3) Pantry and recipes APIs (2:30-3:30)

- [ ] `P0` `GET /ingredients/search?q=` for pantry add flow.
- [ ] `P0` `POST /pantry/items` (ingredient + quantity + unit).
- [ ] `P0` `PATCH /pantry/items/:id` (increment/decrement quantity).
- [ ] `P0` `DELETE /pantry/items/:id`.
- [ ] `P0` `GET /recipes/:id` with ingredient quantities and macros.
- [ ] `P1` Purchases endpoints for budget tracking:
  - `POST /purchases`
  - `GET /budget/summary`

### 4) Plan persistence and macro aggregation (3:30-5:00)

- [ ] `P0` `POST /plans/proposal` ingestion endpoint for AI-normalized payload.
- [ ] `P0` `POST /plans/:id/accept` persists to meal plans.
- [ ] `P0` `POST /plans/:id/decline` archives/rejects and supports regeneration flow.
- [ ] `P0` `GET /plans/week?start=` returns 7-day calendar payload.
- [ ] `P0` Macro totals service in responses:
  - per meal
  - per day
  - per week

### 5) Consumption and pantry deduction (5:00-6:00)

- [ ] `P0` `POST /plan-meals/:id/consume` marks consumed with timestamp.
- [ ] `P0` Deduct pantry ingredients by recipe quantities.
- [ ] `P0` Prevent negative pantry values (clamp at zero + warning field).
- [ ] `P0` Log deductions to `consumption_log` for audit/debug.

### 6) Hardening for demo (6:00-7:00)

- [ ] `P0` Add consistent error format for frontend handling.
- [ ] `P0` Add lightweight request logging and key action tracing.
- [ ] `P0` Add smoke checks for auth -> plan accept -> consume flow.
- [ ] `P1` Add favorite proposal endpoint if time remains.

## Contracts Needed From Others

- PM: frozen priority and final scope boundaries.
- AI: normalized proposal JSON schema and ingredient matching policy.
- Frontend: required response fields per screen.
- QA: top-priority test scenarios and failure expectations.

## Risks

- AI names not matching ingredient table.
- Unit mismatches (g/ml/piece) creating deduction inaccuracies.
- Time overrun from building too many endpoints before core path is stable.

## Done Criteria

- API supports full happy path from auth to pantry auto-deduction.
- Weekly plan endpoint returns 7-day, 4-section layout with macro totals.
- DB setup/seed is reproducible in one command.
- Frontend can run without backend contract guesswork.
