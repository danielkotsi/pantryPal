# Session D - AI/RAG Tasks

## Owner

- AI Engineer (Gemini integration + output normalization)

## Objective

- Provide reliable meal generation through Gemini with deterministic output shape for backend/frontend use.

## Ordered Task List

### 1) Integration baseline (0:30-1:30)

- [ ] `P0` Implement Gemini client wrapper with timeout/retry guards.
- [ ] `P0` Configure API key loading via env variables.
- [ ] `P0` Define request templates for:
  - single meal
  - day plan
  - week plan
  - month plan

### 2) Output schema and parser (1:30-3:00)

- [ ] `P0` Define strict JSON schema for AI output.
- [ ] `P0` Include required fields per meal:
  - meal_type (breakfast/lunch/dinner/snack)
  - recipe_name
  - ingredient list with quantity and unit
  - macros (protein, carbs, fat, calories)
  - estimated cost
- [ ] `P0` Build validation/parsing layer with clear error reasons.
- [ ] `P0` Reject/repair malformed responses before backend persistence.

### 3) Context and personalization (3:00-4:30)

- [ ] `P0` Inject user context: preferences, body metrics, budget target, pantry snapshot.
- [ ] `P0` Add constraints for 4-sections/day structure.
- [ ] `P0` Add prompt guardrails for budget awareness and ingredient realism.
- [ ] `P1` Add simple conversation memory window for follow-up prompts.

### 4) Ingredient matching and macro consistency (4:30-5:30)

- [ ] `P0` Build ingredient name normalization to map AI strings to DB ingredients.
- [ ] `P0` Return unmatched ingredients list for backend/UI warnings.
- [ ] `P0` Reconcile macro totals per meal/day/week in normalized output.
- [ ] `P1` Add alias dictionary for common ingredient synonyms.

### 5) Accept/decline and fallback mode (5:30-6:30)

- [ ] `P0` Support regenerate flow with decline reason/context.
- [ ] `P0` Ensure accepted proposal payload is stable for persistence.
- [ ] `P0` Implement fallback canned generator if Gemini fails/timeouts.
- [ ] `P0` Expose fallback status flag for frontend banner.

### 6) Demo hardening (6:30-7:30)

- [ ] `P0` Add observability fields (request id, model latency, parse status).
- [ ] `P0` Prepare 2-3 tested prompt presets for live demo reliability.
- [ ] `P1` Add favorites metadata output shape for optional feature.

## Contracts Needed From Others

- Backend: ingredient canonical list, plan proposal ingest schema, error contract.
- Frontend: expected chat/proposal rendering shape.
- PM: scope lock for personalization depth and month-plan expectations.
- QA: failure-mode cases for malformed AI responses.

## Risks

- Non-deterministic model output breaking parser.
- Ingredient mismatch rates too high for pantry deduction.
- API rate limits/latency impacting live demo.

## Done Criteria

- AI output consistently validates against schema.
- Plans generate for meal/day/week/month requests with 4 sections/day.
- Ingredient mapping returns canonical IDs or explicit unmatched list.
- Fallback generator can replace Gemini without breaking the demo path.
