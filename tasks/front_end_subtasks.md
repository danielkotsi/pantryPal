# Frontend Subtasks

## Task A: Build the Chat Page (`frontend/src/js/pages/chat.js`)

**Goal:** Turn the chat template into a working conversation interface with AI plan generation actions.

**Pages to touch:**
- New file: `frontend/src/js/pages/chat.js`
- Modify: `frontend/src/js/router/router.js` (wire chat page handler on navigation)
- Modify: `frontend/index.html` (add action button bar to the chat template)
- Modify: `frontend/src/js/api/api-client.js` (add `generatePlan(periodType, message?)` method)
- Modify: `frontend/src/css/styles/main.css` (chat action bar + proposal preview styles)

**Backend endpoints to call:**
- `GET /chat` → `ChatHistoryResponse { messages: [{ id, role, action?, content, createdAt }] }`
- `POST /chat` → `ChatMessageResponse` (request body: `{ message, action? }`)
- `POST /plans/generate` → `GeneratePlanResponse { proposal: ProposalResponse, fallbackActive: bool }` (request body: `{ periodType: "meal"|"day"|"week"|"month", message? }`)
- `POST /plans/{id}/accept` → `PlanSummaryResponse`
- `POST /plans/{id}/decline` → (request body: `{ reason? }`)

**Requirements:**

1. **Load chat history on page entry** — When user navigates to `/chat`, call `GET /chat` and render all past messages (user messages on the right, bot messages on the left, with timestamps). The router's route handler should call the chat page's `init()`.

2. **Fix the inline chat handler** — Remove the `handleChatSubmit` from `router.js` (currently at line 336). The new `chat.js` should own all chat logic. The send endpoint expects `{ message, action? }` (note: the current inline handler sends `{ message }` without the action field, which is correct for plain chat but needs to support the action field for generation). The response is a `ChatMessageResponse` object (not `{ reply: "..." }` as the inline handler currently assumes).

3. **Action button bar** — Add a row of action buttons above the chat input: **"Generate Meal"**, **"Generate Day"**, **"Generate Week"**, **"Generate Month"**. Each button calls `POST /plans/generate` with the corresponding `periodType`. The backend will either use Gemini AI (if configured) or return a fallback plan. The response includes a `fallbackActive` flag.

4. **Proposal preview on generation** — When the generate endpoint returns a proposal, render a preview card in the chat area showing:
   - Period type, status, dates, source (AI/fallback), version, total cost
   - Number of days and meals
   - Per-day macro totals and week totals
   - **Accept** and **Decline** buttons
   - Accept calls `POST /plans/{id}/accept`, which activates the plan. Decline calls `POST /plans/{id}/decline` with an optional reason form.

5. **Scrolling** — Auto-scroll to bottom on new messages.

6. **Styling** — Add CSS for the action button bar, proposal preview card, accept/decline buttons, and ensure the chat messages area grows/shrinks appropriately.

---

## Task B: Add AI Generate Button to the Planner Page

**Goal:** Add a "Generate Week Plan" button in the planner that uses the AI generation endpoint and displays the result.

**Pages to touch:**
- `frontend/src/js/pages/planner.js`
- `frontend/index.html` (add button position in planner header)
- `frontend/src/js/api/api-client.js` (add `generatePlan(periodType, message?)` method if not already done)

**Backend endpoints:**
- `POST /plans/generate` → `GeneratePlanResponse { proposal: ProposalResponse, fallbackActive: bool }` (request body: `{ periodType: "week", message? }`)
- `POST /plans/{id}/accept`
- `POST /plans/{id}/decline`

**Requirements:**

1. **Add a "Generate with AI" button** in the planner header area (next to the view toggle or period navigation). It should be visible only on the week view.

2. **On click**, call `POST /plans/generate` with `periodType: "week"` and optionally include the user's profile preferences/notes from state as the `message` field so the AI has context.

3. **While loading**, show a loading state on the button (e.g., "Generating..." with spinner).

4. **On success**, the response contains a full `ProposalResponse`:
   - Replace the current week view with the generated proposal's days/meals/macros
   - Show a banner/overlay indicating "AI Generated Plan — Accept or Decline?"
   - **Accept button** → calls `POST /plans/{id}/accept`, then reloads the week plan from `GET /plans/week`
   - **Decline button** → opens a small inline form for an optional reason, calls `POST /plans/{id}/decline`, then falls back to the previously accepted plan or generated mock data
   - Show `fallbackActive` badge if the plan was generated from seeded data rather than AI

5. **On error**, show an error message in the planner's error area and keep the current week view unchanged.

6. **Styling** — The generate button, accept/decline banner, and fallback badge should be clearly styled and feel like a temporary overlay that can be dismissed.
