# Session E - QA Tasks

## Owner

- QA Engineer / Bug Triage Lead

## Objective

- Validate end-to-end happy path, catch critical failures early, and keep the demo stable.

## Current Status Snapshot

- No QA scripts, bug tracker docs, or automated tests are implemented yet.
- The only testable implemented surface right now is backend health/auth/profile plus DB reset flow.

## Ordered Task List

### 1) Test plan setup (0:30-1:00)

- [ ] `P0` Create concise test checklist aligned to demo flow.
- [ ] `P0` Define severity levels:
  - `S0` demo blocker
  - `S1` major workaround needed
  - `S2` minor issue/cosmetic
- [ ] `P0` Define bug report template (steps, expected, actual, env, screenshot/log).

### 2) Early validation pass (1:00-3:00)

- [ ] `P0` Validate auth flow: register, login, logout, bad credentials.
- [ ] `P0` Validate profile updates: metrics, preferences, budget persistence.
- [ ] `P0` Validate DB reset/setup flow and seeded demo data.
- [ ] `P0` Validate pantry flow: search/add/remove/decrement behavior once implemented.
- [ ] `P0` Confirm API error messages are actionable for UI display.

### 3) AI and plan flow validation (3:00-5:30)

- [ ] `P0` Validate chat actions for meal/day/week/month generation.
- [ ] `P0` Validate proposal structure: 4 sections/day and macros present.
- [ ] `P0` Validate accept persists plan and appears in weekly calendar.
- [ ] `P0` Validate decline triggers regenerate behavior.
- [ ] `P1` Validate fallback mode when AI service fails.

### 4) Consumption and data integrity checks (5:30-6:30)

- [ ] `P0` Mark meal consumed and verify pantry deduction accuracy.
- [ ] `P0` Verify no negative pantry values are shown.
- [ ] `P0` Verify macro totals remain consistent after consumption actions.
- [ ] `P1` Verify budget/purchase summary if implemented.

### 5) Final regression and demo sign-off (6:30-8:00)

- [ ] `P0` Execute full end-to-end regression script once before demo.
- [ ] `P0` Re-test all fixed `S0`/`S1` bugs.
- [ ] `P0` Produce final known-issues list with workarounds.
- [ ] `P0` Sign-off checklist for PM with go/no-go recommendation.

## Contracts Needed From Others

- PM: canonical demo script and must-pass checkpoints.
- Backend: stable test data/reset process and endpoint error contracts.
- Frontend: expected UI states for loading/errors/empty data.
- AI: expected fallback behavior and schema validation errors.

## Risks

- Late integration creates clustered `S0` defects.
- Test data drift causes inconsistent reproducibility.
- Fallback behavior untested until late, risking live demo failures.

## Done Criteria

- All `S0` issues are resolved or have accepted workaround paths.
- End-to-end script passes in a clean environment.
- Known issues and impact are documented for demo presenters.
- QA sign-off delivered to PM before final demo run.

## Immediate QA Focus

- Start with backend API smoke coverage for `/health`, auth, and profile.
- Create reusable test user/data notes based on the seeded local DB.
- Add pantry/plan/AI checks only after those features land.
