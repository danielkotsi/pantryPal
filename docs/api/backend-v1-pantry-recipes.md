# Backend API v1 - Pantry and Recipes

Base URL (local): `http://localhost:8080`

All endpoints below require:

```http
Authorization: Bearer <token>
```

## `GET /ingredients/search?q=`

Searches USDA foods for pantry add flow.

Example:

```http
GET /ingredients/search?q=banana
```

Response `200`:

```json
{
  "items": [
    {
      "fdcId": 1105314,
      "description": "Bananas, ripe and slightly ripe, raw",
      "foodClass": "FinalFood"
    }
  ]
}
```

## `GET /pantry/items`

Response `200`:

```json
{
  "items": [
    {
      "id": "pnt_...",
      "quantity": 250,
      "unit": "g",
      "food": {
        "fdcId": 1105314,
        "description": "Bananas, ripe and slightly ripe, raw",
        "foodClass": "FinalFood"
      }
    }
  ]
}
```

## `POST /pantry/items`

Request:

```json
{
  "fdcId": 1105314,
  "quantity": 250,
  "unit": "g"
}
```

Response `201`:

```json
{
  "id": "pnt_...",
  "quantity": 250,
  "unit": "g",
  "food": {
    "fdcId": 1105314,
    "description": "Bananas, ripe and slightly ripe, raw",
    "foodClass": "FinalFood"
  }
}
```

If the same user already has the same `fdcId + unit`, quantity is incremented.

## `PATCH /pantry/items/:id`

Request:

```json
{
  "quantityDelta": -50
}
```

Response `200`:

```json
{
  "id": "pnt_...",
  "quantity": 200,
  "unit": "g",
  "food": {
    "fdcId": 1105314,
    "description": "Bananas, ripe and slightly ripe, raw",
    "foodClass": "FinalFood"
  }
}
```

Quantity clamps at `0` if the delta would go negative.

## `DELETE /pantry/items/:id`

Response `204` with empty body.

## `GET /recipes/:id`

Response `200`:

```json
{
  "id": "rcp_breakfast_oats",
  "name": "Banana Oat Bowl",
  "mealType": "breakfast",
  "servings": 1,
  "instructions": "Cook oats in milk, top with sliced banana.",
  "estimatedCostCents": 180,
  "macros": {
    "calories": 430,
    "proteinG": 17,
    "carbsG": 69,
    "fatG": 10
  },
  "ingredients": [
    {
      "fdcId": 1105314,
      "description": "Bananas, ripe and slightly ripe, raw",
      "quantity": 120,
      "unit": "g"
    }
  ]
}
```
