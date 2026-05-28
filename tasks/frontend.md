# Session C — Frontend Engineer

## Your job
Build the remaining 3 page handlers (planner, pantry, chat) and fix the API client URLs. Then wire the consume meal action.

## Already done (do not redo)
- `index.html` with all 5 page templates (auth, profile, planner, pantry, chat) + navbar + error/loading
- `router.js` — SPA routing, state management, auth guard, navigation, localStorage persistence
- `api-client.js` — HTTP methods, Bearer token, error handling (but endpoint URLs are WRONG — fix in Step 1)
- `pages/profile.js` — full profile form with metrics/preferences/budget submit handlers
- `app.js` — initialization, error/loading display
- `css/styles/main.css` — complete responsive styles

## Step-by-step execution

### Step 1 — Fix API client URLs (30 mins, unblocks everything)
**File: `frontend/src/js/api/api-client.js`**

Current wrong routes → correct routes:
- `/pantry` → `/pantry/items` (GET, POST)
- `/pantry/${id}` → `/pantry/items/${id}` (PUT → PATCH, DELETE)
- `/meal-plans` → `/plans/proposal` (POST)
- `/meal-plans/${id}` → `/plans/${id}/accept` (PUT → POST), `/plans/${id}/decline` (PUT → POST)
- `/meal-plans?date=` → `/plans/week?start=` (GET)
- Add missing: `searchIngredients(query)` → `GET /ingredients/search?q=`
- Add missing: `getRecipe(id)` → `GET /recipes/{id}`
- Fix: `updatePantryItem` should PATCH with `{ quantityDelta }` not PUT with full body
- Fix: `addMealPlan` → `createProposal` with `PlanProposalRequest` shape

### Step 2 — Planner page handler (most complex, 1-2 hours)
**File: `frontend/src/js/pages/planner.js`**

- [ ] Create `PlannerPageHandler` class matching `ProfilePageHandler` pattern
- [ ] On navigation to planner route, call `GET /plans/week?start={monday}` for the current week
- [ ] Render a 7-column grid, one column per day
- [ ] Each day column shows 4 sections: breakfast, lunch, dinner, snacks
- [ ] Each section shows: recipe name, ingredient names, macros (protein/carbs/fat/calories)
- [ ] At the bottom of each day column: daily macro total
- [ ] Below the grid: weekly macro total row
- [ ] Add a week navigation control (prev/next week buttons)
- [ ] Handle empty state (no plan yet) and error state
- [ ] Register the page in `initializeRoutes()` in `router.js`

### Step 3 — Pantry page handler (1 hour)
**File: `frontend/src/js/pages/pantry.js`**

- [ ] Create `PantryPageHandler` class
- [ ] Search input: on typing, call `GET /ingredients/search?q=` and show results
- [ ] Selecting a result opens an "add to pantry" form with quantity + unit fields
- [ ] Below search: list of current pantry items from `GET /pantry/items`
- [ ] Each item shows: food name, quantity, unit
- [ ] Each item has + / - buttons to increment/decrement via `PATCH /pantry/items/{id}`
- [ ] Each item has a delete button via `DELETE /pantry/items/{id}`
- [ ] Register in router

### Step 4 — Chat page handler with proposal actions (1-2 hours)
**File: `frontend/src/js/pages/chat.js`**

- [ ] Create `ChatPageHandler` class
- [ ] Load chat history from `GET /chat` on mount
- [ ] Action buttons row above the chat input:
  - "Create Meal" → sends action=meal to backend
  - "Create Day Plan" → action=day
  - "Create Week Plan" → action=week
  - "Create Month Plan" → action=month
- [ ] When a proposal is returned, render it below the chat as a preview card:
  - Plan type, date range, total cost
  - Days with 4 meal sections, macros per meal
- [ ] Below the preview: Accept button (POST /plans/{id}/accept) and Decline button (POST /plans/{id}/decline)
- [ ] On accept, navigate to planner page to show the accepted plan
- [ ] On decline, show a confirmation and allow regeneration
- [ ] Register in router

### Step 5 — Consume meal action (30 mins)
- [ ] In the planner page, add a "Consume" button on each meal section
- [ ] On click, call `POST /plan-meals/{id}/consume`
- [ ] Show success feedback and refresh the planner data
- [ ] Show a banner: "Pantry updated — X items deducted"
- [ ] Handle error: "Not enough pantry stock for Y ingredient"

### Step 6 — Polish
- [ ] Add loading/empty states for planner, pantry, chat
- [ ] Ensure all pages handle 401 response by redirecting to auth
- [ ] Test full flow manually: register → profile → chat → generate week → accept → planner → consume

## Dependencies

| You need from      | What                                    |
|--------------------|-----------------------------------------|
| Session B (backend)| Consumption endpoint shape (Step 4/5)   |
| Session B (backend)| Chat endpoint shape (Step 4)            |
| Session B (backend)| Generate endpoint shape (Step 4)       |
| Session D (AI)     | Example proposal JSON for rendering mock|

## Done criteria
- Planner shows 7 days x 4 sections with macros from real API data
- Pantry search/add/remove/decrement all work against backend
- Chat has action buttons, proposal preview, accept/decline
- Consume meal button deducts pantry and shows feedback
- No hardcoded/incorrect API URLs remain
