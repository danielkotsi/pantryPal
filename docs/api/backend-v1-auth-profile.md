# Backend API v1 - Auth and Profile

Base URL (local): `http://localhost:8080`

## Error format

All failures return:

```json
{
  "error": {
    "code": "validation_error",
    "message": "human readable message"
  }
}
```

## Authentication

- Token type: Bearer token from `Authorization: Bearer <token>`
- Token expiry: controlled by `TOKEN_TTL_HOURS` (default `24`)

## Endpoints

### `POST /auth/register`

Request:

```json
{
  "email": "newuser@pantrypal.local",
  "password": "Passw0rd!",
  "displayName": "New User"
}
```

Response `201`:

```json
{
  "token": "<bearer-token>",
  "expiresAt": "2026-05-29T12:00:00Z",
  "user": {
    "id": "usr_...",
    "email": "newuser@pantrypal.local",
    "displayName": "New User"
  }
}
```

### `POST /auth/login`

Request:

```json
{
  "email": "newuser@pantrypal.local",
  "password": "Passw0rd!"
}
```

Response `200`: same shape as register response.

### `GET /me`

Headers: `Authorization: Bearer <token>`

Response `200`:

```json
{
  "user": {
    "id": "usr_...",
    "email": "newuser@pantrypal.local",
    "displayName": "New User"
  }
}
```

### `GET /profile`

Headers: `Authorization: Bearer <token>`

Response `200`:

```json
{
  "user": {
    "id": "usr_...",
    "email": "newuser@pantrypal.local",
    "displayName": "New User"
  },
  "metrics": {
    "heightCm": 172,
    "weightKg": 70,
    "age": 29,
    "sex": "male",
    "activityLevel": "moderate",
    "goal": "maintain"
  },
  "preferences": {
    "dietType": "omnivore",
    "allergies": ["peanuts"],
    "dislikes": ["mushrooms"],
    "likes": ["rice"],
    "dailyCalorieTarget": 2300,
    "notes": "test"
  },
  "budget": {
    "month": "2026-06",
    "currency": "USD",
    "amountCents": 50000
  }
}
```

### `PATCH /profile/metrics`

Headers: `Authorization: Bearer <token>`

Request:

```json
{
  "heightCm": 172,
  "weightKg": 70,
  "age": 29,
  "sex": "male",
  "activityLevel": "moderate",
  "goal": "maintain"
}
```

Response `200`: full `GET /profile` shape.

### `PATCH /profile/preferences`

Headers: `Authorization: Bearer <token>`

Request:

```json
{
  "dietType": "omnivore",
  "allergies": ["peanuts"],
  "dislikes": ["mushrooms"],
  "likes": ["rice"],
  "dailyCalorieTarget": 2300,
  "notes": "test"
}
```

Response `200`: full `GET /profile` shape.

### `PATCH /profile/budget`

Headers: `Authorization: Bearer <token>`

Request:

```json
{
  "month": "2026-06",
  "currency": "USD",
  "amountCents": 50000
}
```

Response `200`: full `GET /profile` shape.

### `GET /health`

Response `200`:

```json
{
  "status": "ok"
}
```
