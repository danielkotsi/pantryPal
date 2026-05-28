# Session C - Frontend Tasks

## Owner

- Frontend Engineer (Vanilla JS)

## Objective

- Build a clean, fast demo UI that proves all core flows using the backend API.

## Ordered Task List

### 1) App shell and API client (0:30-1:30)

- [ ] `P0` Set up page structure and shared layout (auth, profile, planner, pantry, chat).
- [ ] `P0` Build API client helpers for auth headers, JSON handling, and error display.
- [ ] `P0` Add simple route/state switching (no framework) between major sections.

### 2) Auth and profile screens (1:30-2:30)

- [ ] `P0` Register + login forms connected to API.
- [ ] `P0` Persist auth token/session in browser storage for demo.
- [ ] `P0` Profile form for body metrics, preferences, budget.
- [ ] `P0` Save/update profile and show success/error feedback.

### 3) Weekly planner UI (2:30-4:30)

- [ ] `P0` Build 7-day calendar-like layout.
- [ ] `P0` Render 4 sections per day: breakfast, lunch, dinner, snacks.
- [ ] `P0` Show ingredients, quantities, and macros per meal.
- [ ] `P0` Show macro totals per day and per week.
- [ ] `P1` Add lightweight monthly toggle view if backend payload is ready.

### 4) Pantry and consumption flow (4:30-5:30)

- [ ] `P0` Pantry ingredient search input with result list.
- [ ] `P0` Add item with quantity/unit.
- [ ] `P0` Increment/decrement/remove pantry items.
- [ ] `P0` Add "consume meal" action in planner and refresh pantry state.

### 5) Chat and plan actions (5:30-6:30)

- [ ] `P0` Chat panel with action buttons:
  - create meal
  - create day plan
  - create week plan
  - create month plan
- [ ] `P0` Render proposal preview returned by API.
- [ ] `P0` Accept and decline buttons wired to proposal endpoints.
- [ ] `P1` Favorite action button if endpoint is available.

### 6) Demo polish and resilience (6:30-7:30)

- [ ] `P0` Add loading and empty states for each screen.
- [ ] `P0` Add fail-safe banner if AI service fallback mode is active.
- [ ] `P0` Ensure mobile-usable layout and legible data density.
- [ ] `P0` Prepare one-click path for live demo navigation.

## Contracts Needed From Others

- Backend: stable endpoint contracts and sample responses.
- AI: proposal payload format and error/fallback response shape.
- PM: final demo path and required screen order.
- QA: bug priority list for final hour fixes.

## Risks

- Contract changes late in build causing UI rewiring.
- Overbuilding visuals before core actions are connected.
- Too much on one page reducing clarity in live demo.

## Done Criteria

- User can authenticate, edit profile, generate/accept a weekly plan, and consume a meal.
- Planner clearly shows 7 days x 4 sections with macro totals.
- Pantry updates reflect manual edits and automatic meal consumption.
- Chat actions are usable via buttons without manual API calls.
