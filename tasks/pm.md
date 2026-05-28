# Session A — Product Manager

## Your job
Keep the build coordinated, resolve blockers quickly, and ensure the demo flow works end-to-end.

## What is already complete
- All backend routes except consumption, chat, and AI wiring
- AI client, parser, schema, prompt builder, unit tests
- Frontend shell, router, API client, auth/profile pages, CSS
- DB schema, seed data, reset scripts

## What each agent is building (order-independent)

### Session B — Backend
| Step | Task | Est. time |
|------|------|-----------|
| 1 | Consumption endpoint | 1.5h |
| 2 | Chat endpoint | 1h |
| 3 | Wire AI into generate endpoint | 1h |
| 4 | Hardening + logging middleware | 30m |

### Session C — Frontend
| Step | Task | Est. time |
|------|------|-----------|
| 1 | Fix API client URLs | 30m |
| 2 | Planner page handler | 1.5h |
| 3 | Pantry page handler | 1h |
| 4 | Chat page with proposal actions | 1.5h |
| 5 | Consume meal button | 30m |
| 6 | Polish | 30m |

### Session D — AI
| Step | Task | Est. time |
|------|------|-----------|
| 1 | Fallback canned plan generator | 1h |
| 2 | Wire Gemini into app.go | 1h |
| 3 | Ingredient matching service | 1h (P1) |
| 4 | Prompt presets for demo | 30m |

### Session E — QA
| Step | Task | Est. time |
|------|------|-----------|
| 1 | Smoke test existing endpoints | 1h |
| 2 | Re-test as new features land | ongoing |
| 3 | Final regression + sign-off | 1h |

## Execution order (recommended)

```
Hour 0     Hour 1     Hour 2     Hour 3     Hour 4     Hour 5     Hour 6     Hour 7
├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤
B: Consumption ──────►│ B: Chat ──────────►│ B: Wire AI ──────────►│ B: Harden │
                      │                     │                        │           │
C: Fix URLs ►│ C: Planner ───────────────►│ C: Pantry ───────────►│ C: Chat+Proposals ──►│ C: Polish│
              │                            │                        │                      │          │
D: Fallback ──►│ D: Wire Gemini into app ──►│ D: Ingredient matching (P1)                  │
                │                            │                        │                      │
E: Smoke tests ────────────────►│ E: Ongoing regression ──────────────────────────────────►│ E: Signoff│
```

## Key checkpoints for you
- **End of Hour 1**: Frontend API URLs fixed; fallback generator exists
- **End of Hour 2**: Consumption endpoint testable; planner page renders static data
- **End of Hour 3**: Chat endpoint works; pantry page works; AI wiring started
- **End of Hour 4**: Accept/decline flow through chat works end-to-end
- **End of Hour 5**: Consume meal button deducts pantry; full demo loop works
- **End of Hour 6**: All P0 features complete; QA begins full regression
- **End of Hour 7**: QA sign-off; bug fixes; final demo prepped

## Blockers to watch
1. **Frontend API client URLs** — fix this first, else every frontend page will fail
2. **Gemini API key** — without it, AI generation will always use fallback. Make sure Session D has the key.
3. **DB path** — the Go binary runs from `backend/` with default DB path `../database/sqlite/pantrypal.db`. If frontend runs `go run` from a different directory, the path breaks. Standardize: `cd backend && go run ./cmd/api`
4. **CORS** — frontend runs on a different port/domain. Middleware is already in place but verify origin is allowed.

## Fallback plan
If Gemini is unavailable at demo time:
- The fallback generator in Session D Step 1 provides a deterministic week plan
- The demo flow works identically — the user just sees "AI fallback mode" banner
- No feature is blocked

If any frontend page is incomplete:
- The demo script can use `curl` commands to show the API working
- PM should prepare curl commands for each step as a safety net
