# PantryPal Hackathon Plan (8 Hours)

This is the master coordination board for a short demo build.

## Current Status Snapshot

- Backend foundation is implemented for app bootstrap, auth, profile, token auth, DB migrations, demo seeds, and local DB reset scripts.
- Backend docs exist for auth/profile API and DB setup.
- Frontend is not implemented yet beyond folder scaffolding.
- AI/RAG integration is not implemented yet beyond folder scaffolding.
- QA automation/checklists are not implemented yet.

## Demo Goal

Ship a working end-to-end flow where a user can:
1. Sign up/login.
2. Set body metrics, food preferences, and monthly budget.
3. Generate a daily/weekly meal plan with AI.
4. View meals in a calendar-style UI with macros and ingredient quantities.
5. Accept a generated plan and store it.
6. Manage pantry items and mark meals as consumed to auto-deduct pantry stock.

## Scope Guardrails (Hackathon)

- Prioritize **working happy paths** over edge-case completeness.
- Build daily + weekly plan first; monthly can reuse same pipeline.
- Keep auth simple (email/password + session/JWT).
- Use SQLite locally; no production deployment requirements.
- Cost values can come from AI output for demo purposes.

## Priorities

- `P0` = must have for demo.
- `P1` = strong add if time permits.
- `P2` = optional stretch.

## Master Task Order

### Phase 0 - Project setup (0:00-0:30)

- [x] `P0` Confirm repo structure and ownership for sessions:
  - `/tasks/TASKS.md` (this file)
  - `/tasks/pm.md`
  - `/tasks/backend.md`
  - `/tasks/frontend.md`
  - `/tasks/ai.md`
  - `/tasks/bugs.md`
- [ ] `P0` Define API contract skeleton (auth, profile, pantry, plans, chat).
- [x] `P0` Define auth/security minimums (password hashing, token expiry).
- [ ] `P0` Decide timebox checkpoints at 2h / 4h / 6h / 8h.

### Phase 1 - Data model and auth foundation (0:30-2:00)

- [x] `P0` Create SQLite schema + migrations for:
  - users
  - user_preferences
  - user_body_metrics
  - budgets / purchases
  - USDA foods dataset tables
  - recipes + recipe_ingredients
  - pantry_items
  - meal_plans (day/week/month)
  - plan_meals
  - consumption_log
  - chat_messages
- [x] `P0` Provide one-command DB setup/reset for all sessions.
- [x] `P0` Seed ingredient data and demo app data.
- [x] `P0` Implement auth endpoints:
  - register
  - login
  - current user (`/me`)
- [x] `P0` Implement profile endpoints to read/update:
  - body metrics
  - preferences
  - monthly budget
- [ ] `P0` Freeze API v1 contract at 2h mark so frontend/AI can continue in parallel.

### Phase 2 - Pantry + recipes + macros pipeline (2:00-3:30)

- [ ] `P0` Pantry endpoints:
  - search ingredients
  - add item with quantity
  - remove item / decrement quantity
- [ ] `P0` Recipe read endpoints with ingredients + macros.
- [ ] `P0` Macro aggregation service:
  - per meal
  - per day
  - per week
- [ ] `P1` Budget tracking endpoint for purchases and current monthly spend.

### Phase 3 - AI meal plan generation and acceptance (3:30-5:30)

- [ ] `P0` Integrate Gemini API service wrapper.
- [ ] `P0` Add fallback generator (seeded/canned plan) if Gemini is unavailable.
- [ ] `P0` Chat endpoint with actions:
  - request meal
  - request day plan
  - request week plan
  - request month plan
- [ ] `P0` Normalize AI output into internal plan format:
  - day split into 4 sections (breakfast, lunch, dinner, snacks)
  - ingredient names mapped to DB ingredients
  - quantities, macros, and AI-provided cost captured
- [ ] `P0` Accept/decline flow:
  - decline => regenerate proposal
  - accept => persist to meal_plans + plan_meals
- [ ] `P1` Add favorite action (store proposal without scheduling).
- [ ] `P1` Decide chat retention for demo (store full messages vs summary metadata).

### Phase 4 - Frontend demo flows (5:30-7:00)

- [ ] `P0` Vanilla JS auth screens (register/login/logout).
- [ ] `P0` Profile screen for metrics, preferences, budget editing.
- [ ] `P0` Calendar-like weekly view with:
  - 7 days
  - 4 meal sections/day
  - total macros/day and week
- [ ] `P0` Pantry management UI:
  - ingredient search
  - add/remove quantity
- [ ] `P0` Chat interface with action buttons:
  - create meal
  - create day plan
  - create week plan
  - create month plan
  - accept / decline
- [ ] `P0` Mark meal consumed action triggers pantry deduction.

### Phase 5 - Stabilize, test, and demo prep (7:00-8:00)

- [ ] `P0` End-to-end happy-path test script:
  - sign up
  - set profile
  - generate week plan
  - accept plan
  - view macros
  - consume one meal
  - verify pantry deduction
- [ ] `P0` Fix top-priority bugs and broken UI/API links.
- [ ] `P0` Prepare demo data for one realistic user journey.
- [ ] `P1` Add monthly view polish.
- [ ] `P2` Optional: change one meal on a selected day.
- [ ] `P2` Optional: supermarket list from pantry gaps.

## Definition of Done for Demo

- User can authenticate and update profile (metrics/preferences/budget).
- User can generate and accept a weekly plan from chat.
- Weekly view shows meals and macro totals.
- Pantry can be edited manually and auto-deducts on meal consumption.
- System runs locally with SQLite and seeded data.

## Implemented Now

- Go backend entrypoint and app wiring.
- SQLite schema for core auth/profile/pantry/recipes/plans/chat data plus USDA food tables.
- Demo seed SQL for one user, pantry items, recipes, and recipe ingredients.
- Local DB bootstrap/reset scripts.
- `GET /health`, `POST /auth/register`, `POST /auth/login`, `GET /me`.
- `GET /profile`, `PATCH /profile/metrics`, `PATCH /profile/preferences`, `PATCH /profile/budget`.
- Auth/profile API documentation.

## Still Missing

- Pantry search and CRUD endpoints/UI.
- Recipe read endpoints and macro rollups.
- Meal plan creation, proposal, accept/decline, and weekly view payloads.
- Consumption logging endpoint and pantry auto-deduction logic.
- Gemini integration, fallback planner, and chat interface.
- Any frontend implementation.
- QA checklist, regression script, and automated tests.

## Session Ownership Map

- Session A (Product Manager): roadmap, priorities, specs, integration coordination.
- Session B (Backend): DB schema, auth, pantry, plans, macros, consumption logic.
- Session C (Frontend): auth/profile UI, calendar views, pantry UI, chat UI.
- Session D (AI/RAG): Gemini integration, prompt format, output normalization, ingredient matching.
- Session E (QA): test checklist, bug triage, regression pass before demo.

## Handoff Notes for Next Step

Next files to generate from this master:
- `/tasks/pm.md`
- `/tasks/backend.md`
- `/tasks/frontend.md`
- `/tasks/ai.md`
- `/tasks/bugs.md`

Each file should include:
- owner
- objective
- ordered task list
- API/UI contracts needed from others
- risk list
- done criteria
