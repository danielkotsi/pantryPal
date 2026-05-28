# Backend-Frontend Route Mapping

## Current Backend Routes (from `router.go`)

```
GET  /health                    → HealthHandler.GetHealth
POST /auth/register             → AuthHandler.Register
POST /auth/login                → AuthHandler.Login
GET  /me                        → AuthHandler.Me (auth required)
GET  /profile                   → ProfileHandler.GetProfile (auth required)
PATCH /profile/metrics          → ProfileHandler.PatchMetrics (auth required)
PATCH /profile/preferences      → ProfileHandler.PatchPreferences (auth required)
PATCH /profile/budget           → ProfileHandler.PatchBudget (auth required)
```

## Authentication
- **Type**: Bearer Token (JWT)
- **Header Format**: `Authorization: Bearer <token>`
- **CORS**: Enabled for all origins
- **Allowed Methods**: GET, POST, PATCH, DELETE, OPTIONS
- **Token Validation**: Extracts UserID from JWT claims

## Frontend API Client Implementation

### Auth Endpoints
✅ `api.signup(email, password, displayName)` → `POST /auth/register`
✅ `api.login(email, password)` → `POST /auth/login`
✅ `api.getMe()` → `GET /me`
✅ `api.logout()` → Clears token locally

### Profile Endpoints
✅ `api.getProfile()` → `GET /profile`
✅ `api.updateMetrics(metrics)` → `PATCH /profile/metrics`
✅ `api.updatePreferences(preferences)` → `PATCH /profile/preferences`
✅ `api.updateBudget(budget)` → `PATCH /profile/budget`

## Frontend Features Implemented

### 1. Authentication Flow
- Register/Login with email, password, displayName
- JWT token stored in localStorage
- Auto-include Bearer token in all requests
- Token refresh on restore from localStorage

### 2. Profile Management
- View personal info (email, displayName)
- Edit body metrics (height, weight, age, sex, activity level, goal)
- Edit preferences (diet type, allergies, dislikes, likes, calorie target, notes)
- Edit budget (month, currency, amount in cents)
- Auto-reload after updates

### 3. State Management
- Persistent user state in localStorage
- Auto-restore on page reload
- Protected routes (redirect to auth if not authenticated)
- Global error and loading state handling

### 4. UI Components
- Auth page with login/signup toggle
- Profile page with editable sections
- Navigation with active route highlighting
- Error messages with auto-dismiss
- Loading spinner for API calls

## Future Integrations Needed

Based on backend module structure, these endpoints will need frontend implementation:
- **Pantry** - `internal/modules/pantry`
- **Recipes** - `internal/modules/recipes`
- **Plans** - `internal/modules/plans`
- **Chat** - `internal/modules/chat`
- **Budget** - `internal/modules/budget`
- **AI** - `internal/modules/ai`

## CORS Configuration

Backend allows:
- ✅ Origin: `*` (all origins)
- ✅ Methods: GET, POST, PATCH, DELETE, OPTIONS
- ✅ Headers: Authorization, Content-Type
- ✅ Credentials: Not required (no cookies)

Frontend can make requests from any origin.

## Response Format

All responses follow this pattern:

### Success Response
```json
{
  "token": "jwt_token",
  "expiresAt": "2026-05-29T...",
  "user": {
    "id": "user_id",
    "email": "user@example.com",
    "displayName": "User Name"
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "error_code",
    "message": "Human readable message"
  }
}
```

Error codes:
- `bad_request` - Invalid JSON or validation error
- `unauthorized` - Missing/invalid token or credentials
- `conflict` - Resource already exists (e.g., email)
- `validation_error` - Invalid data
- `internal_error` - Server error
- `not_found` - Resource not found
