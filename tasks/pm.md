# Session A - Product Manager Tasks

## Owner

- Product Manager / Integration Lead

## Objective

- Keep the 8-hour demo focused on a working end-to-end flow.
- Lock scope quickly, unblock cross-session dependencies, and manage risk/time.

## Ordered Task List

### 0) Kickoff and scope lock (0:00-0:30)

- [ ] `P0` Confirm demo storyline: signup -> profile setup -> generate week plan -> accept -> pantry update -> consume meal.
- [ ] `P0` Confirm out-of-scope items for this sprint (full optimization, advanced settings, heavy polish).
- [ ] `P0` Assign session owners and communication cadence (check-ins at 2h/4h/6h/8h).
- [ ] `P0` Create shared decision log with timestamps and owner.

### 1) Spec freeze for parallel work (0:30-2:00)

- [ ] `P0` Freeze API v1 contract by 2h mark.
- [ ] `P0` Freeze AI response shape v1 (meals/day split, ingredient fields, macros, cost).
- [ ] `P0` Freeze frontend MVP screens and navigation map.
- [ ] `P0` Define acceptance rules for plan proposal (accept, decline/regenerate, favorite optional).

### 2) Integration management (2:00-6:30)

- [ ] `P0` Resolve cross-team blockers in <=15 minutes each.
- [ ] `P0` Verify backend/frontend contract alignment after every major endpoint.
- [ ] `P0` Verify AI ingredient matching strategy and fallback behavior.
- [ ] `P1` Keep release notes updated with what is demo-ready vs partially complete.

### 3) Demo readiness (6:30-8:00)

- [ ] `P0` Run full walkthrough with QA observer.
- [ ] `P0` Prepare demo script with exact clicks and API actions.
- [ ] `P0` Mark risks and fallback moves for each step (especially AI outage).
- [ ] `P0` Final go/no-go checklist with owners for each section.

## Contracts Needed From Others

- Backend: stable endpoint list, payload examples, migration/setup command.
- Frontend: final routes/screens and expected loading/error states.
- AI: prompt template, JSON schema, fallback plan output format.
- QA: final bug severity rubric and test checklist coverage.

## Risks

- API contract drift after 2h causing rework.
- AI output variability breaking parser.
- Pantry deduction logic unclear at edge cases.
- Time spent on optional features reducing demo stability.

## Done Criteria

- API and AI contracts are frozen and shared.
- End-to-end demo script passes at least once without manual DB edits.
- Each session has explicit must-have vs nice-to-have status.
- Fallback demo path exists if AI service fails.
