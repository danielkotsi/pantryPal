# Quick Reference Guide

## Backend Routes at a Glance

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/health` | ❌ | Health check |
| POST | `/auth/register` | ❌ | Create account |
| POST | `/auth/login` | ❌ | Login user |
| GET | `/me` | ✅ | Get current user |
| GET | `/profile` | ✅ | Get full profile |
| PATCH | `/profile/metrics` | ✅ | Update metrics |
| PATCH | `/profile/preferences` | ✅ | Update preferences |
| PATCH | `/profile/budget` | ✅ | Update budget |

✅ = Requires Bearer token in `Authorization` header

## Frontend API Methods

```javascript
// Auth
api.signup(email, password, displayName)
api.login(email, password)
api.logout()
api.getMe()

// Profile
api.getProfile()
api.updateMetrics(metrics)
api.updatePreferences(preferences)
api.updateBudget(budget)
```

## Frontend Pages

| Page | Route | Auth Required | Purpose |
|------|-------|---------------|-----------| 
| Auth | `/` or `#auth` | ❌ | Login/Register |
| Profile | `#profile` | ✅ | View & edit profile |
| Planner | `#planner` | ✅ | Meal planning (WIP) |
| Pantry | `#pantry` | ✅ | Manage items (WIP) |
| Chat | `#chat` | ✅ | Chat assistant (WIP) |

## Common Errors

| Status | Code | Meaning |
|--------|------|---------|
| 400 | `bad_request` | Invalid input |
| 400 | `validation_error` | Failed validation |
| 401 | `unauthorized` | Invalid/missing token |
| 409 | `conflict` | Email already exists |
| 500 | `internal_error` | Server error |

## File Structure

```
frontend/
├── index.html                   # Main page
├── src/
│   ├── js/
│   │   ├── api/api-client.js   # HTTP client
│   │   ├── router/router.js     # Routes & state
│   │   ├── app/app.js           # Initialization
│   │   └── pages/profile.js     # Profile logic
│   └── css/styles/main.css      # Styling

backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── transport/http/
│   │   └── router/router.go     # All routes
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   └── repositories/            # Database access
└── migrations/                  # DB schema
```

## Quick Start

### Start Backend
```bash
cd backend
go run cmd/api/main.go
# Runs on http://localhost:8080
```

### Start Frontend
```bash
cd frontend
python -m http.server 3000
# Open http://localhost:3000 in browser
```

### Test Registration
1. Click "Sign Up"
2. Enter email, password, display name
3. See profile page

### Test Login
1. Click "Logout"
2. Enter credentials
3. See profile page

## Token Storage

- **Frontend**: `localStorage.authToken`
- **Format**: JWT (JSON Web Token)
- **Used In**: `Authorization: Bearer <token>` header
- **Lifetime**: 24 hours (configurable)

## Environment Variables (Backend)

```bash
PORT=8080                   # API port
DB_PATH=./db.db             # Database path
TOKEN_SECRET=secret-key     # JWT secret
TOKEN_TTL_HOURS=24          # Token lifetime
```

## Response Format

All responses are JSON:

**Success:**
```json
{
  "token": "...",
  "expiresAt": "...",
  "user": {
    "id": "...",
    "email": "...",
    "displayName": "..."
  }
}
```

**Error:**
```json
{
  "error": {
    "code": "error_code",
    "message": "Error description"
  }
}
```

## Frontend State

```javascript
router.state = {
  user: null,              // Current user
  isAuthenticated: false,  // Login status
  loading: false,          // API call status
  error: null              // Error message
}
```

## Useful Debugging Tips

1. **Check Token**
   ```javascript
   // In browser console
   localStorage.getItem('authToken')
   ```

2. **Check State**
   ```javascript
   // In browser console
   router.getState()
   ```

3. **Test API**
   ```javascript
   // In browser console
   api.login('email@example.com', 'password')
   ```

4. **Check Network Tab**
   - F12 → Network tab
   - Perform action
   - Click request
   - Check Headers and Response

5. **Enable Logging**
   - Backend logs to console automatically
   - Frontend errors show in DevTools

## Database

- **Type**: SQLite
- **Location**: `../database/sqlite/pantrypal.db`
- **Schema**: `backend/migrations/001_init_schema.sql`
- **Demo Data**: `backend/seeds/001_seed_demo.sql`

## Authentication Flow

```
┌──────────────────────────┐
│  User enters credentials │
└───────────┬──────────────┘
            ↓
┌──────────────────────────┐
│  POST /auth/login        │
└───────────┬──────────────┘
            ↓
┌──────────────────────────┐
│  Backend validates       │
│  Generates JWT           │
└───────────┬──────────────┘
            ↓
┌──────────────────────────┐
│  Return token + user     │
└───────────┬──────────────┘
            ↓
┌──────────────────────────┐
│  Store token locally     │
│  Update app state        │
│  Navigate to profile     │
└──────────────────────────┘
```

## For Future Development

Placeholder routes (ready to implement):
- `GET /pantry` - Get pantry items
- `POST /pantry` - Add item
- `GET /recipes` - Get recipes
- `GET /meals` - Get meal plans
- `POST /chat` - Send message

See backend modules in `internal/modules/` for structure.
