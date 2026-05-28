# Session D — AI / RAG Engineer

## Your job
Build the fallback plan generator, wire the Gemini client into the app, and build ingredient matching against the USDA foods table.

## Already done (do not redo)
- `ai/client.go` — Gemini HTTP client with timeout, retry, config from env
- `ai/schema.go` — `PlanResponse`, `DayPlan`, `MealDetail`, `IngredientInput`, `MacroInput`, `NormalizedPlan`
- `ai/parser.go` — `ParsePlanResponse()` with markdown fence repair, numeric string coercion, strict unknown field rejection, meal type normalization, duplicate section detection, date validation, cost conversion
- `ai/prompts.go` — `BuildPrompt()` for meal/day/week/month with user context injection (metrics, preferences, budget, pantry snapshot)
- `ai/parser_test.go` — 4 passing tests

## Step-by-step execution

### Step 1 — Fallback canned plan generator (1 hour, independent)
**File: `backend/internal/modules/ai/fallback.go`**

- [ ] Create `GenerateFallbackWeekPlan(userID string, startDate string)` function
- [ ] Return a `NormalizedPlan` with 7 days, 4 meals/day using the seeded demo recipes
- [ ] Map to the same `PlanResponse` → `ParsePlanResponse` flow so output is identical format
- [ ] Use demo recipes from `001_seed_demo.sql`: Banana Oat Bowl (breakfast), Chicken Rice Plate (lunch), Egg Fried Rice Lite (dinner), Yogurt Apple Snack (snacks)
- [ ] Rotate or repeat across 7 days with slight variations
- [ ] Set realistic macros and costs matching seed data
- [ ] Export as `func FallbackWeekPlan(startDate string) (NormalizedPlan, error)` so backend can call it

### Step 2 — Wire Gemini into the generate endpoint (1 hour, coordinate with Session B)
**File: `backend/internal/services/generate_service.go`** (new) or extend **`backend/internal/services/chat_service.go`**

- [ ] Create a service function: `GeneratePlan(ctx, userID, requestType, userMessage)`
- [ ] Logic:
  1. Build prompt via `ai.BuildPrompt()` with user context (fetch from profile/pantry APIs)
  2. Call `ai.Client.Generate()` with the prompt
  3. Call `ai.ParsePlanResponse()` on the result
  4. Call `backend PlanService.CreateProposal()` with `NormalizedPlan.Proposal`
  5. Return the proposal response to the caller
- [ ] If Gemini client returns error or times out → call `FallbackWeekPlan()` instead
- [ ] Return a `fallbackActive: true` flag in the response so frontend can show a banner
- [ ] Wire into `app.go` — initialize Gemini client via `ai.NewClient(ai.ConfigFromApp(cfg))`
- [ ] Test: set `GEMINI_API_KEY=test`, make a request, verify parser repairs work on real output

### Step 3 — Ingredient matching service (1 hour, optional P0/P1)
**File: `backend/internal/services/ingredient_service.go`** (new)

- [ ] Build `MatchIngredients(ingredients []NormalizedIngredient) ([]MatchedIngredient, []UnmatchedIngredient)`
- [ ] Query `usda_foods` by `description LIKE %name%` for each ingredient name
- [ ] Use a simple fuzzy match: if a single result matches, use `fdc_id`; if multiple, pick the closest by Levenshtein or word overlap
- [ ] Return matched (with `fdc_id`) and unmatched (with original name) separately
- [ ] This is used by the consumption flow to know what pantry items to deduct
- [ ] Optional: add a synonym map in `ai/ingredients.go` for common mismatches ("chicken breast" ↔ "Chicken, breast, meat only, raw")

### Step 4 — Prompt presets for demo (30 mins)
- [ ] Prepare 2-3 tested prompt templates that reliably produce valid JSON
- [ ] Document the exact `GEMINI_API_KEY`, `GEMINI_MODEL`, and `GEMINI_TIMEOUT_SECONDS` env vars needed
- [ ] Test that `ParsePlanResponse()` handles real Gemini output (not just test fixtures)

## Dependencies

| You need from      | What                                    |
|--------------------|-----------------------------------------|
| Session B (backend) | PlanService.CreateProposal() is ready   |
| Session B (backend) | Generate endpoint to hook into         |
| Session C (frontend)| Feedback on fallback banner flag shape |

## Done criteria
- `FallbackWeekPlan()` returns valid `NormalizedPlan` that can be persisted
- Gemini → parse → proposal flow works end-to-end when API key is set
- Ingredient matching returns `fdc_id` or explicit unmatched list
- Fallback mode activates when Gemini is unavailable
