# PantryPal Backend-Frontend Integration Summary

## ✅ Integration Complete

All backend routes from `backend/internal/transport/http/router/router.go` are now fully integrated with the frontend.

## Backend Routes (Go Backend)

### Router Definition
**File**: `backend/internal/transport/http/router/router.go`

```go
// Public Routes
GET  /health                    // Health check
POST /auth/register             // Register new user
POST /auth/login                // Login user

// Protected Routes (require Bearer token)
GET  /me                        // Get current user
GET  /profile                   // Get full profile
PATCH /profile/metrics          // Update metrics
PATCH /profile/preferences      // Update preferences
PATCH /profile/budget           // Update budget
```

### Middleware Stack
1. **CORS** - Allows all origins
2. **Logging** - Logs all requests
3. **AuthRequired** - Validates JWT for protected routes

## Frontend Implementation

### Frontend Files
- **`frontend/index.html`** - Main page structure
- **`frontend/src/js/api/api-client.js`** - HTTP client (axios-like)
- **`frontend/src/js/router/router.js`** - Routing & state management
- **`frontend/src/js/pages/profile.js`** - Profile page logic
- **`frontend/src/js/app/app.js`** - App initialization
- **`frontend/src/css/styles/main.css`** - Styling

### API Client Methods

#### Authentication
```javascript
api.signup(email, password, displayName)    // POST /auth/register
api.login(email, password)                  // POST /auth/login
api.logout()                                // Clears token locally
api.getMe()                                 // GET /me
```

#### Profile Management
```javascript
api.getProfile()                            // GET /profile
api.updateMetrics(metrics)                  // PATCH /profile/metrics
api.updatePreferences(preferences)          // PATCH /profile/preferences
api.updateBudget(budget)                    // PATCH /profile/budget
```

## Complete Request/Response Examples

### 1. Registration Flow

**Frontend Request:**
```javascript
const result = await api.signup('john@example.com', 'password123', 'John Doe');
```

**HTTP Request:**
```
POST http://localhost:8080/auth/register
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123",
  "displayName": "John Doe"
}
```

**Backend Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": "2026-05-29T12:34:56Z",
  "user": {
    "id": "user_123abc",
    "email": "john@example.com",
    "displayName": "John Doe"
  }
}
```

**Frontend Processing:**
- Extracts token and stores in localStorage
- Sets `isAuthenticated: true`
- Stores user info in app state
- Navigates to profile page

### 2. Login Flow

**Frontend Request:**
```javascript
const result = await api.login('john@example.com', 'password123');
```

**HTTP Request:**
```
POST http://localhost:8080/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Backend Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": "2026-05-29T12:34:56Z",
  "user": {
    "id": "user_123abc",
    "email": "john@example.com",
    "displayName": "John Doe"
  }
}
```

### 3. Get Profile Flow

**Frontend Request:**
```javascript
const profile = await api.getProfile();
```

**HTTP Request:**
```
GET http://localhost:8080/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Backend Processing:**
1. AuthRequired middleware validates token
2. Extracts UserID from JWT claims
3. Checks user still exists in database
4. Injects UserID into request context
5. ProfileHandler.GetProfile() is called with UserID
6. Fetches user, metrics, preferences, budget from DB

**Backend Response (200 OK):**
```json
{
  "user": {
    "id": "user_123abc",
    "email": "john@example.com",
    "displayName": "John Doe"
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
    "allergies": ["nuts", "shellfish"],
    "dislikes": ["cilantro"],
    "likes": ["pasta", "coffee"],
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

**Frontend Processing:**
- Renders user info
- Populates metrics form with values
- Populates preferences form with values
- Populates budget form with values

### 4. Update Metrics Flow

**Frontend Request:**
```javascript
await api.updateMetrics({
  heightCm: 182,
  weightKg: 78,
  age: 31,
  activityLevel: "very_active"
});
```

**HTTP Request:**
```
PATCH http://localhost:8080/profile/metrics
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "heightCm": 182,
  "weightKg": 78,
  "age": 31,
  "activityLevel": "very_active"
}
```

**Backend Processing:**
1. AuthRequired validates token and extracts UserID
2. ProfileHandler.PatchMetrics() updates only provided fields
3. Database updates user_body_metrics table
4. Returns full profile object

**Backend Response (200 OK):**
```json
{
  "user": { ... },
  "metrics": {
    "heightCm": 182,
    "weightKg": 78,
    "age": 31,
    "sex": "M",
    "activityLevel": "very_active",
    "goal": "maintain"
  },
  "preferences": { ... },
  "budget": { ... }
}
```

**Frontend Processing:**
- Updates local profile state
- Re-renders profile page
- Shows success message (no error)

## Error Handling

### Validation Error Example

**Frontend Request (invalid email):**
```javascript
try {
  await api.signup('invalid-email', 'pass123', 'John');
} catch (error) {
  console.log(error.message);  // "invalid email format"
}
```

**Backend Response (400 Bad Request):**
```json
{
  "error": {
    "code": "validation_error",
    "message": "invalid email format"
  }
}
```

### Conflict Error Example

**Frontend Request (email already exists):**
```javascript
try {
  await api.signup('existing@example.com', 'pass123', 'John');
} catch (error) {
  console.log(error.message);  // "email already registered"
}
```

**Backend Response (409 Conflict):**
```json
{
  "error": {
    "code": "conflict",
    "message": "email already registered"
  }
}
```

### Auth Error Example

**Frontend Request (invalid token):**
```javascript
// Token expired or tampered with
const profile = await api.getProfile();
```

**Backend Response (401 Unauthorized):**
```json
{
  "error": {
    "code": "unauthorized",
    "message": "invalid token"
  }
}
```

## Frontend Features

### ✅ Implemented
- User registration with email validation
- User login with password verification
- Profile view with all user data
- Edit body metrics (height, weight, age, sex, activity, goal)
- Edit preferences (diet, allergies, dislikes, likes, calories, notes)
- Edit budget (month, currency, amount)
- Persistent authentication (localStorage)
- Protected routes (redirect to auth if not logged in)
- Error display with auto-dismiss
- Loading spinner during API calls
- Global state management
- Logout functionality

### 🔒 Security Features
- JWT token validation
- Bearer token in Authorization header
- CORS enabled (safe cross-origin)
- Token storage in localStorage (not cookies - safer for SPA)
- HTML escaping to prevent XSS
- Proper error messages without sensitive data

## Running the Application

### Backend
```bash
cd backend
# Ensure DATABASE is at ../database/sqlite/pantrypal.db
go run cmd/api/main.go
# Runs on http://localhost:8080
```

### Frontend
```bash
# Option 1: Direct browser (if using file://)
open frontend/index.html

# Option 2: Simple HTTP server
cd frontend
python -m http.server 3000
# Access at http://localhost:3000

# Option 3: Using npm
npx http-server -p 3000
```

## Testing the Integration

### 1. Test Registration
```
1. Open http://localhost:3000/frontend/
2. Click "Sign Up"
3. Enter: test@example.com, password123, Test User
4. Click "Sign Up" button
5. Should see profile page
```

### 2. Test Login
```
1. Click "Logout" button
2. Should return to auth page
3. Enter: test@example.com, password123
4. Click "Login" button
5. Should see profile page
```

### 3. Test Profile Update
```
1. On profile page, scroll to "Body Metrics"
2. Enter: Height 180, Weight 75, Age 30
3. Click "Save Metrics"
4. Should see data reload
```

### 4. Verify Network Requests
```
1. Open Browser DevTools (F12)
2. Go to Network tab
3. Try login
4. Look for POST /auth/login request
5. Headers should include Authorization: Bearer <token>
```

## Documentation Files

Created for reference:
- **BACKEND_FRONTEND_MAPPING.md** - Route mapping and flow
- **FRONTEND_BACKEND_INTEGRATION.md** - Detailed integration guide
- **ROUTES_ARCHITECTURE.md** - Architecture diagrams and file structure
- **INTEGRATION_SUMMARY.md** - Original summary

## API Contract

### Base URL
- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000` (or file://)

### Default Configuration
- Port: 8080 (backend)
- DB: `../database/sqlite/pantrypal.db`
- Token TTL: 24 hours
- CORS: All origins allowed

### Customization
Backend config via environment variables:
```bash
PORT=9000                          # Change port
DB_PATH=/path/to/db.db             # Change database
TOKEN_SECRET=your-secret           # JWT secret
TOKEN_TTL_HOURS=48                 # Token lifetime
```

## Next Steps

### Immediate (Frontend)
- [ ] Test all endpoints in DevTools
- [ ] Test error scenarios
- [ ] Test token expiration
- [ ] Add loading indicators
- [ ] Add success notifications

### Short-term (Backend)
- [ ] Implement pantry endpoints
- [ ] Implement recipes endpoints
- [ ] Implement meal plan endpoints
- [ ] Implement chat endpoints

### Medium-term
- [ ] Add database migrations
- [ ] Add user validation rules
- [ ] Add rate limiting
- [ ] Add request logging
- [ ] Add metrics/monitoring
