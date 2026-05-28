# Backend Routes Architecture

## Complete Route Map

```
┌─────────────────────────────────────────────────────────────┐
│                    PantryPal Backend API                    │
│                  (http://localhost:8080)                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────── PUBLIC ROUTES ──────────────────────┐
│                                                              │
│  GET /health                                                │
│  └─> HealthHandler.GetHealth()                             │
│      Returns: { status: "ok" }                              │
│                                                              │
│  POST /auth/register                                        │
│  └─> AuthHandler.Register()                                │
│      Input:  { email, password, displayName }              │
│      Output: { token, expiresAt, user }                    │
│                                                              │
│  POST /auth/login                                           │
│  └─> AuthHandler.Login()                                   │
│      Input:  { email, password }                           │
│      Output: { token, expiresAt, user }                    │
│                                                              │
└──────────────────────────────────────────────────────────────┘

┌──────────────── PROTECTED ROUTES (Bearer Token) ────────────┐
│                                                              │
│  Middleware: AuthRequired()                                 │
│  └─> Extracts JWT token from Authorization header          │
│  └─> Validates token                                       │
│  └─> Injects UserID into request context                   │
│                                                              │
│  GET /me                                                    │
│  └─> AuthHandler.Me()                                      │
│      Returns: { user: { id, email, displayName } }         │
│                                                              │
│  GET /profile                                               │
│  └─> ProfileHandler.GetProfile()                           │
│      Returns: { user, metrics, preferences, budget }       │
│                                                              │
│  PATCH /profile/metrics                                     │
│  └─> ProfileHandler.PatchMetrics()                         │
│      Input:  { heightCm?, weightKg?, age?, ... }           │
│      Output: { user, metrics, preferences, budget }        │
│                                                              │
│  PATCH /profile/preferences                                 │
│  └─> ProfileHandler.PatchPreferences()                     │
│      Input:  { dietType?, allergies?, dislikes?, ... }     │
│      Output: { user, metrics, preferences, budget }        │
│                                                              │
│  PATCH /profile/budget                                      │
│  └─> ProfileHandler.PatchBudget()                          │
│      Input:  { month?, currency?, amountCents? }           │
│      Output: { user, metrics, preferences, budget }        │
│                                                              │
└──────────────────────────────────────────────────────────────┘

┌──────────────── MIDDLEWARE STACK ────────────────────────────┐
│                                                               │
│  Request → CORS → Logging → [AuthRequired?] → Handler       │
│                                                               │
│  CORS:                                                       │
│  • Allow-Origin: *                                           │
│  • Allow-Methods: GET, POST, PATCH, DELETE, OPTIONS         │
│  • Allow-Headers: Authorization, Content-Type               │
│                                                               │
│  Logging:                                                    │
│  • Logs: METHOD PATH (duration)                             │
│  • Example: "PATCH /profile/metrics (15ms)"                 │
│                                                               │
│  AuthRequired:                                               │
│  • Extracts "Bearer <token>" from Authorization header      │
│  • Validates JWT token                                      │
│  • Verifies user still exists                               │
│  • Injects UserID into request context                      │
│                                                               │
└──────────────────────────────────────────────────────────────┘

┌──────────────── FRONTEND API CLIENT ────────────────────────┐
│                                                               │
│  class APIClient {                                           │
│    baseURL: "http://localhost:8080"                         │
│                                                               │
│    // Auth Methods                                           │
│    signup(email, password, displayName)                     │
│    login(email, password)                                   │
│    logout()                                                  │
│    getMe()                                                   │
│                                                               │
│    // Profile Methods                                        │
│    getProfile()                                              │
│    updateMetrics(metrics)                                   │
│    updatePreferences(preferences)                           │
│    updateBudget(budget)                                     │
│                                                               │
│    // Utilities                                              │
│    setToken(token)                                          │
│    buildHeaders()      // Adds Authorization header         │
│    onError(callback)   // Error event listener              │
│    onLoadingChange()   // Loading state listener            │
│  }                                                           │
│                                                               │
└──────────────────────────────────────────────────────────────┘

┌──────────────── REQUEST/RESPONSE FLOW ──────────────────────┐
│                                                               │
│  Frontend                          Backend                   │
│  ────────────────────────────────────────────────────────   │
│                                                               │
│  1. User submits login form                                 │
│     ↓                                                        │
│  2. api.login(email, password)                              │
│     ↓                                                        │
│  3. POST /auth/login (JSON)                                 │
│     ├─ Content-Type: application/json                      │
│     └─ Body: { email, password }                           │
│        ────────────────────────────────→                   │
│                                         AuthHandler.Login()  │
│                                         ↓                    │
│                                         Verify credentials   │
│                                         ↓                    │
│                                         Generate JWT token   │
│                                         ↓                    │
│        ←────────────────────────────────                   │
│        200 OK                                               │
│        {                                                    │
│          "token": "eyJ...",                                │
│          "expiresAt": "2026-05-29T...",                   │
│          "user": { "id", "email", "displayName" }         │
│        }                                                    │
│     ↓                                                        │
│  4. Store token in localStorage                            │
│     ↓                                                        │
│  5. Navigate to profile page                               │
│     ↓                                                        │
│  6. api.getProfile()                                       │
│     ↓                                                        │
│  7. GET /profile                                           │
│     ├─ Authorization: Bearer <token>                       │
│     └─ Content-Type: application/json                      │
│        ────────────────────────────────→                   │
│                                         AuthRequired()       │
│                                         ↓                    │
│                                         Validate token       │
│                                         ↓                    │
│                                         ProfileHandler      │
│                                         .GetProfile()        │
│                                         ↓                    │
│        ←────────────────────────────────                   │
│        200 OK                                               │
│        {                                                    │
│          "user": { ... },                                  │
│          "metrics": { ... },                               │
│          "preferences": { ... },                           │
│          "budget": { ... }                                 │
│        }                                                    │
│     ↓                                                        │
│  8. Render profile page with user data                     │
│                                                               │
└──────────────────────────────────────────────────────────────┘

┌──────────────── ERROR HANDLING ─────────────────────────────┐
│                                                               │
│  Backend Response (Error)                                   │
│  ─────────────────────────                                 │
│  HTTP 401 Unauthorized                                      │
│  {                                                          │
│    "error": {                                               │
│      "code": "unauthorized",                                │
│      "message": "invalid token"                             │
│    }                                                        │
│  }                                                          │
│                                                               │
│  Frontend Handling                                          │
│  ──────────────────                                        │
│  ↓ Catch APIError                                          │
│  ↓ Extract error.message: "invalid token"                 │
│  ↓ Call router.setState({ error: message })               │
│  ↓ Display error in UI                                     │
│  ↓ Auto-dismiss after 5 seconds                            │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

## Backend File Structure

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go                 # Entry point
│   └── dbseed/
│       └── main.go                 # Database seeder
├── internal/
│   ├── app/
│   │   └── app.go                  # App initialization
│   ├── config/
│   │   └── config.go               # Configuration
│   ├── modules/
│   │   ├── auth/
│   │   ├── profile/
│   │   ├── pantry/
│   │   ├── recipes/
│   │   ├── plans/
│   │   ├── chat/
│   │   ├── budget/
│   │   └── ai/
│   ├── platform/
│   │   ├── auth/
│   │   │   └── token.go            # JWT handling
│   │   ├── db/
│   │   │   ├── bootstrap.go        # DB initialization
│   │   │   └── sqlite.go           # SQLite driver
│   │   ├── httpserver/
│   │   ├── id/
│   │   │   └── id.go               # ID generation
│   │   └── logger/
│   ├── repositories/
│   │   └── user_repository.go      # User DB access
│   ├── services/
│   │   ├── auth_service.go         # Auth logic
│   │   └── profile_service.go      # Profile logic
│   └── transport/
│       └── http/
│           ├── handlers/
│           │   ├── auth_handler.go
│           │   ├── health_handler.go
│           │   ├── profile_handler.go
│           │   └── response.go
│           ├── middleware/
│           │   ├── auth.go         # JWT validation
│           │   ├── common.go       # CORS, Logging
│           │   └── context.go      # Context helpers
│           ├── dto/
│           │   └── types.go        # Request/Response types
│           └── router/
│               └── router.go       # Route registration
├── migrations/
│   └── 001_init_schema.sql         # DB schema
├── seeds/
│   └── 001_seed_demo.sql           # Demo data
└── go.mod
```

## Frontend File Structure

```
frontend/
├── index.html                       # Main HTML
├── src/
│   ├── js/
│   │   ├── api/
│   │   │   └── api-client.js       # HTTP client
│   │   ├── router/
│   │   │   └── router.js           # Routing & state
│   │   ├── app/
│   │   │   └── app.js              # App initialization
│   │   └── pages/
│   │       └── profile.js          # Profile logic
│   └── css/
│       └── styles/
│           └── main.css            # Styles
└── [other static files]
```

## Key Integration Points

1. **Authentication**
   - Backend: JWT token generation
   - Frontend: Token storage and Bearer header injection

2. **Error Handling**
   - Backend: Standard error response format
   - Frontend: APIError class and error display

3. **CORS**
   - Backend: Allows all origins
   - Frontend: Can make requests from any origin

4. **State Management**
   - Backend: Stateless (JWT-based)
   - Frontend: Client-side state with localStorage persistence

5. **API Format**
   - Backend: RESTful with JSON
   - Frontend: Structured APIClient class
