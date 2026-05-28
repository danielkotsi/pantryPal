# Frontend-Backend Integration Guide

## Backend Routes (Go)

### Location: `backend/internal/transport/http/router/router.go`

```go
// Public endpoints (no auth required)
GET  /health                    // Health check
POST /auth/register             // Create account
POST /auth/login                // Login

// Protected endpoints (Bearer token required)
GET  /me                        // Get current user info
GET  /profile                   // Get full profile
PATCH /profile/metrics          // Update body metrics
PATCH /profile/preferences      // Update preferences
PATCH /profile/budget           // Update budget
```

## Authentication Flow

### 1. Registration
**Frontend:**
```javascript
const result = await api.signup(email, password, displayName);
// Returns: { token, expiresAt, user: { id, email, displayName } }
```

**Backend:**
- Endpoint: `POST /auth/register`
- Request:
  ```json
  { "email": "user@example.com", "password": "pass123", "displayName": "John" }
  ```
- Response (201):
  ```json
  {
    "token": "eyJ...",
    "expiresAt": "2026-05-29T12:00:00Z",
    "user": {
      "id": "user_123",
      "email": "user@example.com",
      "displayName": "John"
    }
  }
  ```
- Error (400/409):
  ```json
  {
    "error": {
      "code": "conflict",
      "message": "email already registered"
    }
  }
  ```

### 2. Login
**Frontend:**
```javascript
const result = await api.login(email, password);
// Returns same as signup
```

**Backend:**
- Endpoint: `POST /auth/login`
- Request:
  ```json
  { "email": "user@example.com", "password": "pass123" }
  ```
- Response (200): Same as registration

### 3. Token Usage
All authenticated requests must include:
```
Authorization: Bearer <token>
```

Frontend API client automatically adds this header via `buildHeaders()`.

## Profile Management

### Get Profile
**Frontend:**
```javascript
const profile = await api.getProfile();
```

**Backend:**
- Endpoint: `GET /profile`
- Headers: `Authorization: Bearer <token>`
- Response (200):
  ```json
  {
    "user": {
      "id": "user_123",
      "email": "user@example.com",
      "displayName": "John"
    },
    "metrics": {
      "heightCm": 180,
      "weightKg": 75,
      "age": 30,
      "sex": "M",
      "activityLevel": "moderate",
      "goal": "maintain"
    },
    "preferences": {
      "dietType": "omnivore",
      "allergies": ["nuts"],
      "dislikes": ["cilantro"],
      "likes": ["pasta"],
      "dailyCalorieTarget": 2000,
      "notes": "vegetarian on weekends"
    },
    "budget": {
      "month": "2026-05",
      "currency": "USD",
      "amountCents": 50000
    }
  }
  ```

### Update Metrics
**Frontend:**
```javascript
await api.updateMetrics({
  heightCm: 180,
  weightKg: 75,
  age: 30,
  sex: "M",
  activityLevel: "moderate",
  goal: "maintain"
});
```

**Backend:**
- Endpoint: `PATCH /profile/metrics`
- Headers: `Authorization: Bearer <token>`
- Request: Same as metrics object above (all fields optional)
- Response (200): Full profile object

### Update Preferences
**Frontend:**
```javascript
await api.updatePreferences({
  dietType: "omnivore",
  allergies: ["nuts"],
  dislikes: ["cilantro"],
  likes: ["pasta"],
  dailyCalorieTarget: 2000,
  notes: "vegetarian on weekends"
});
```

**Backend:**
- Endpoint: `PATCH /profile/preferences`
- Headers: `Authorization: Bearer <token>`
- Request: Same as preferences object above (all fields optional)
- Response (200): Full profile object

### Update Budget
**Frontend:**
```javascript
await api.updateBudget({
  month: "2026-05",
  currency: "USD",
  amountCents: 50000
});
```

**Backend:**
- Endpoint: `PATCH /profile/budget`
- Headers: `Authorization: Bearer <token>`
- Request: Budget object (all fields optional)
- Response (200): Full profile object

## Error Handling

### Error Response Format
All errors follow this format:
```json
{
  "error": {
    "code": "error_code",
    "message": "Human readable message"
  }
}
```

### Common Error Codes
- `bad_request` - Invalid input (400)
- `unauthorized` - Missing/invalid token (401)
- `conflict` - Resource conflict, e.g., email exists (409)
- `validation_error` - Validation failed (400)
- `internal_error` - Server error (500)
- `not_found` - Resource not found (404)

### Frontend Error Handling
The API client converts backend errors to `APIError` instances:
```javascript
try {
  await api.login(email, password);
} catch (error) {
  console.log(error.message);        // "invalid credentials"
  console.log(error.status);         // 401
  console.log(error.data.error.code); // "unauthorized"
}
```

## CORS Configuration

Backend enables CORS for all origins:
- **Allow-Origin**: `*`
- **Allow-Methods**: GET, POST, PATCH, DELETE, OPTIONS
- **Allow-Headers**: Authorization, Content-Type

Frontend can make requests from any origin without additional configuration.

## Frontend Implementation Status

### ✅ Completed
- Auth page (login/signup)
- Profile page with forms
- API client with all endpoints
- Token management
- Error handling
- Loading states
- Protected routes

### ⏳ Ready for Implementation
- Pantry management (CRUD pantry items)
- Recipes (browse, search)
- Meal planner (create plans)
- Chat (message history)
- Consumption log (track meals)

## Running the Application

### Backend
```bash
cd backend
go run cmd/api/main.go
# Server runs on http://localhost:8080
```

### Frontend
Simply open `frontend/index.html` in a browser or use a local server:
```bash
# Python 3
python -m http.server 3000

# Node.js
npx http-server -p 3000
```

Access at `http://localhost:3000` (or your preferred port).

## Testing the Integration

1. **Register**
   - Go to http://localhost:3000
   - Click "Sign Up"
   - Enter email, display name, password
   - Should redirect to profile page

2. **Login**
   - Logout (click logout button)
   - Go back to auth page
   - Enter credentials
   - Should redirect to profile page

3. **Update Profile**
   - On profile page, update metrics/preferences/budget
   - Click save
   - Should show success and reload data

4. **Check Network Tab**
   - All requests should include `Authorization: Bearer <token>` header
   - Responses should follow the format above
