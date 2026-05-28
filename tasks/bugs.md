# Session E — QA Engineer

## Your job
Validate the full demo flow, catch bugs early, and sign off before demo.

## Already testable now
- Backend endpoints for auth, profile, pantry, recipes, plans (see full list below)
- AI parser unit tests (4 passing)
- Frontend auth + profile pages

## Step-by-step execution

### Step 1 — Smoke test existing backend endpoints (1 hour)
Run a full manual pass against the running API:

- [ ] `GET /health` → 200
- [ ] `POST /auth/register` with valid fields → 201 + token
- [ ] `POST /auth/register` with invalid email → 400
- [ ] `POST /auth/register` with weak password → 400
- [ ] `POST /auth/register` with duplicate email → 409
- [ ] `POST /auth/login` with correct credentials → 200 + token
- [ ] `POST /auth/login` with wrong password → 401
- [ ] `GET /me` with valid token → 200 + user
- [ ] `GET /me` without token → 401
- [ ] `GET /me` with expired/invalid token → 401
- [ ] `GET /profile` → 200 with metrics/prefs/budget
- [ ] `PATCH /profile/metrics` → 200, verify changes persist on GET
- [ ] `PATCH /profile/preferences` → 200, verify changes persist
- [ ] `PATCH /profile/budget` → 200, verify changes persist
- [ ] `GET /ingredients/search?q=chicken` → 200 with results
- [ ] `GET /ingredients/search?q=` (empty) → 400
- [ ] `GET /pantry/items` → 200 with demo seed items
- [ ] `POST /pantry/items` with valid fdcId/qty/unit → 201
- [ ] `POST /pantry/items` with invalid fdcId → 400
- [ ] `POST /pantry/items` with zero qty → 400
- [ ] `PATCH /pantry/items/{id}` with positive delta → 200, quantity increases
- [ ] `PATCH /pantry/items/{id}` with negative delta → 200, quantity decreases
- [ ] `DELETE /pantry/items/{id}` → 204
- [ ] `DELETE /pantry/items/{id}` (already deleted) → 404
- [ ] `GET /recipes/rcp_breakfast_oats` → 200 with ingredients + macros
- [ ] `GET /recipes/invalid` → 404
- [ ] `POST /plans/proposal` with valid payload → 201
- [ ] `POST /plans/{id}/accept` → 200, status changes to accepted
- [ ] `POST /plans/{id}/accept` (already accepted) → 409
- [ ] `POST /plans/{id}/decline` → 200, status changes to declined
- [ ] `POST /plans/{id}/decline` (already declined) → 409
- [ ] `GET /plans/week?start=2026-06-01` → 200 with days + macro totals

### Step 2 — Run AI parser tests (already passing)
- [ ] Run `go test ./internal/modules/ai/...` and confirm all 4 tests pass

### Step 3 — Test new backend endpoints as they ship (ongoing)

**When consumption endpoint lands:**
- [ ] `POST /plan-meals/{id}/consume` marks meal as consumed
- [ ] Pantry quantities decrease by meal ingredient amounts
- [ ] Pantry quantities clamp at zero (not negative)
- [ ] `consumption_log` has entries for each deducted ingredient
- [ ] Calling consume on already-consumed meal returns appropriate error

**When chat endpoint lands:**
- [ ] `POST /chat` stores message and returns it
- [ ] `GET /chat` returns recent messages
- [ ] Auth required for both

### Step 4 — Test frontend as pages land (ongoing)

**When planner page lands:**
- [ ] Week navigation works (prev/next)
- [ ] 7 columns with 4 sections each render
- [ ] Macro totals per day and per week are correct
- [ ] Empty state when no plan exists
- [ ] Error state when API fails

**When pantry page lands:**
- [ ] Search returns results from API
- [ ] Add item creates pantry entry
- [ ] +/- buttons update quantity
- [ ] Delete removes item

**When chat page lands:**
- [ ] Action buttons trigger generation
- [ ] Proposal preview renders correctly
- [ ] Accept/decline work and update state

**When consume button lands:**
- [ ] Meal marked consumed after click
- [ ] Pantry updates reflected in pantry page

### Step 5 — Final regression + demo sign-off (last hour)
- [ ] Run full end-to-end flow without any manual DB edits:
  1. Reset DB with `./scripts/db/local_db.sh reset-all`
  2. Register a new user
  3. Set profile metrics, preferences, budget
  4. Generate a week plan via chat
  5. Accept the plan
  6. View the plan in the weekly calendar
  7. Consume a meal
  8. Verify pantry deduction
  9. Logout and login again
- [ ] Document all S0/S1 issues with workarounds
- [ ] Deliver sign-off report to PM

## Bug severity rubric
- `S0` — Demo blocker. Full flow breaks, no workaround.
- `S1` — Major. Flow works but with significant UX friction or manual steps.
- `S2` — Minor. Cosmetic issue, non-critical error message, edge case.
