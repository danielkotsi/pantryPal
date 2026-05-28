# Session B — Backend Engineer

## Your job
Implement the two remaining backend endpoints: consumption and chat. Then wire AI into the app.

## Already done (do not redo)
- Auth: `POST /auth/register`, `POST /auth/login`, `GET /me` with bcrypt + HMAC tokens
- Profile: `GET /profile`, `PATCH /profile/metrics`, `PATCH /profile/preferences`, `PATCH /profile/budget`
- Pantry: `GET /ingredients/search`, `GET /pantry/items`, `POST /pantry/items`, `PATCH /pantry/items/{id}`, `DELETE /pantry/items/{id}`
- Recipes: `GET /recipes/{id}` with ingredients + macros
- Plans: `POST /plans/proposal`, `POST /plans/{id}/accept`, `POST /plans/{id}/decline`, `GET /plans/week?start=`
- App bootstrap, config, DB setup, CORS + logging middleware, seed data

## Step-by-step execution

### Step 1 — Consumption endpoint (`POST /plan-meals/:id/consume`)
- [ ] Create `ConsumeService` in `backend/internal/services/consume_service.go`
- [ ] Load the plan meal by ID, verify it belongs to the requesting user
- [ ] Load the recipe (if linked via `recipe_id`) to get ingredient `fdc_id` + quantities
- [ ] For each ingredient, find the user's matching pantry item and deduct the quantity
- [ ] Clamp pantry quantities at zero
- [ ] Log each deduction to `consumption_log` with before/after quantities
- [ ] Mark `plan_meals.is_consumed = 1` and set `consumed_at`
- [ ] Create `ConsumeHandler` in `backend/internal/transport/http/handlers/consume_handler.go`
- [ ] Add route: `mux.Handle("POST /plan-meals/{id}/consume", authRequired(http.HandlerFunc(h.Consume.ConsumeMeal)))`
- [ ] Wire into `app.go` and `router.go`
- [ ] DTOs: `ConsumeMealResponse` with consumed items and warnings

### Step 2 — Chat endpoint (`POST /chat`, `GET /chat`)
- [ ] Create `ChatService` in `backend/internal/services/chat_service.go`
- [ ] `POST /chat`: accept `{ message, action? }`, store user message in `chat_messages`, return the stored message
- [ ] `GET /chat?limit=50`: return recent messages for the authenticated user
- [ ] Create `ChatHandler` in `backend/internal/transport/http/handlers/chat_handler.go`
- [ ] Add routes: `POST /chat`, `GET /chat` (both auth required)
- [ ] Wire into `app.go` and `router.go`
- [ ] DTOs: `ChatSendRequest`, `ChatMessageResponse`, `ChatHistoryResponse`

### Step 3 — Wire AI into the generate flow (collaborate with Session D)
- [ ] Create `GenerateService` or extend `ChatService` to:
  1. Accept a generation request (meal/day/week/month)
  2. Call `ai.BuildPrompt()` with user context (metrics, preferences, budget, pantry)
  3. Call `ai.Client.Generate()` with the prompt
  4. Call `ai.ParsePlanResponse()` on the result
  5. Call `PlanService.CreateProposal()` with the normalized payload
  6. Return the proposal to the frontend
- [ ] Add route: `POST /generate` (or integrate into `POST /chat` with an `action` field)
- [ ] Wire into `app.go` (instantiate Gemini client via `ConfigFromApp`)
- [ ] If Gemini is unavailable — call the fallback generator (Session D builds it in Step 3)

### Step 4 — Hardening
- [ ] Add request logging middleware with method, path, duration, status
- [ ] Run `go build ./...` and fix any compilation errors
- [ ] Test full flow: register → profile → create proposal → accept → week plan → consume → verify pantry deduction

## Dependencies

| You need from      | What                        |
|--------------------|-----------------------------|
| Session D (AI)     | Fallback generator function |
| Session D (AI)     | Gemini client is in `ai/` ready to use |
| Session C (FE)     | Feedback on consume response shape |

## Done criteria
- `POST /plan-meals/:id/consume` deducts pantry and logs consumption
- `POST /chat` + `GET /chat` store and return messages
- Gemini generation wired and reachable via an endpoint
- `go build ./...` passes
