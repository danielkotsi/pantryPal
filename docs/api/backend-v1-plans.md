# Backend API v1 - Plans

Base URL (local): `http://localhost:8080`

All endpoints below require:

```http
Authorization: Bearer <token>
```

## `POST /plans/proposal`

Stores an AI-normalized proposal as a `proposal` plan.

Request:

```json
{
  "periodType": "week",
  "startDate": "2026-06-01",
  "endDate": "2026-06-07",
  "source": "ai",
  "proposalVersion": 1,
  "aiCostCentsTotal": 1234,
  "notes": "weekly demo",
  "meals": [
    {
      "scheduledDate": "2026-06-01",
      "mealSection": "breakfast",
      "recipeId": "rcp_breakfast_oats",
      "recipeName": "Banana Oat Bowl",
      "servings": 1,
      "estimatedCostCents": 180,
      "macros": {
        "calories": 430,
        "proteinG": 17,
        "carbsG": 69,
        "fatG": 10
      }
    }
  ]
}
```

Response `201`:

```json
{
  "plan": {
    "id": "pln_...",
    "periodType": "week",
    "startDate": "2026-06-01",
    "endDate": "2026-06-07",
    "status": "proposal",
    "source": "ai",
    "proposalVersion": 1,
    "aiCostCentsTotal": 1234,
    "notes": "weekly demo"
  },
  "days": [],
  "weekTotals": {
    "calories": 430,
    "proteinG": 17,
    "carbsG": 69,
    "fatG": 10
  }
}
```

## `POST /plans/:id/accept`

Accepts a proposal and promotes it to `accepted`.

Response `200`: same response shape as proposal, with `plan.status = "accepted"`.

## `POST /plans/:id/decline`

Declines a proposal and keeps decline context in plan notes.

Request:

```json
{
  "reason": "too expensive"
}
```

Response `200`: same response shape as proposal, with `plan.status = "declined"`.

## `GET /plans/week?start=`

Returns accepted plans that overlap the 7-day window beginning at `start`.

Example:

```http
GET /plans/week?start=2026-06-01
```

Response `200`:

```json
{
  "plans": [
    {
      "id": "pln_...",
      "periodType": "week",
      "startDate": "2026-06-01",
      "endDate": "2026-06-07",
      "status": "accepted",
      "source": "ai",
      "proposalVersion": 1,
      "aiCostCentsTotal": 1234,
      "notes": "weekly demo"
    }
  ],
  "days": [
    {
      "date": "2026-06-01",
      "sections": {
        "breakfast": {
          "id": "pm_...",
          "recipeId": "rcp_breakfast_oats",
          "recipeName": "Banana Oat Bowl",
          "scheduledDate": "2026-06-01",
          "mealSection": "breakfast",
          "servings": 1,
          "estimatedCostCents": 180,
          "macros": {
            "calories": 430,
            "proteinG": 17,
            "carbsG": 69,
            "fatG": 10
          },
          "isConsumed": false
        }
      },
      "totals": {
        "calories": 430,
        "proteinG": 17,
        "carbsG": 69,
        "fatG": 10
      }
    }
  ],
  "weekTotals": {
    "calories": 1630,
    "proteinG": 93,
    "carbsG": 189,
    "fatG": 52
  }
}
```
