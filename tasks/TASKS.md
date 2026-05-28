# PantryPal Hackathon Plan (8 Hours)

## Demo Goal

Ship a working end-to-end flow:
1. Sign up/login.
2. Set body metrics, food preferences, and monthly budget.
3. Generate a daily/weekly meal plan with AI.
4. View meals in a calendar-style UI with macros and ingredient quantities.
5. Accept a generated plan and store it.
6. Manage pantry items and mark meals as consumed to auto-deduct pantry stock.

## Overall Execution Order (for PM coordination)

### Step 1 — Backend: Consumption endpoint + Chat endpoint (parallel with Step 2)
**Owner: Session B**
- `POST /plan-meals/:id/consume` — deduct pantry, log to `consumption_log`
- `POST /chat` + `GET /chat` — store/retrieve messages

### Step 2 — Frontend: Fix API client + Planner page + Pantry page
**Owner: Session C**
- Fix stale endpoint URLs in `api-client.js`
- Implement `planner.js` — 7-day x 4-section weekly calendar from `GET /plans/week`
- Implement `pantry.js` — ingredient search, add/remove items from pantry API

### Step 3 — AI: Fallback generator + Wire Gemini into app
**Owner: Session D**
- Build canned fallback plan generator (static week plan)
- Wire Gemini client + parser + plan service into `app.go` via a generate endpoint or chat handler

### Step 4 — Frontend: Chat page + Consume button
**Owner: Session C**
- `chat.js` — action buttons (meal/day/week/month) + proposal preview + accept/decline
- Consume meal button in planner wired to consumption endpoint

### Step 5 — QA: Full pass + sign-off
**Owner: Session E**
- Smoke all endpoints
- Run end-to-end script
- Bug triage and known-issues list

## What Is Already Complete (do not redo)
- All backend routes for auth, profile, pantry CRUD, recipes, plan proposals/accept/decline/week read
- AI: Gemini client, JSON schema, parser with repair, prompt builder, unit tests
- Frontend: index.html, SPA router, API client shell, auth/profile pages, CSS
- DB schema, migrations, seed data, reset scripts

## Priorities
- `P0` = must have for demo
- `P1` = strong add if time permits
- `P2` = optional stretch
