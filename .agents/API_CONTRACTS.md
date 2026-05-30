# API_CONTRACTS.md

Purpose: define **stable, minimal API contracts** between Client and Server for the BiteBox MVP.

Rules:

- Backward-compatible changes only during MVP
- No undocumented fields
- Errors are explicit and consistent

Base URL: `/api`
Auth: `Authorization: Bearer <JWT>` (required unless stated)

---

## Auth

### POST `/auth/login`

**Auth:** No

**Request**

```json
{
  "email": "string",
  "password": "string"
}
```

**Response 200**

```json
{
  "token": "string",
  "user": { "id": "uuid", "email": "string" }
}
```

**Errors**

- `401` invalid credentials

---

### POST `/auth/register`

**Auth:** No

**Request**

```json
{
  "email": "string",
  "password": "string"
}
```

**Response 201**

```json
{
  "token": "string",
  "user": { "id": "uuid", "email": "string" }
}
```

---

## Recipes

### GET `/recipes`

**Auth:** Optional

**Query Params**

- `page` (int)
- `limit` (int)

**Response 200**

```json
{
  "items": [Recipe],
  "page": 1,
  "has_more": true
}
```

---

### GET `/recipes/{id}`

**Auth:** Optional

**Response 200**

```json
Recipe
```

**Errors**

- `404` not found or inactive

---

### POST `/recipes`

**Auth:** Yes

**Request**

```json
{
  "title": "string",
  "ingredients": ["string"],
  "steps": ["string"],
  "meal_type_id": "int"
}
```

**Response 201**

```json
Recipe
```

---

### PUT `/recipes/{id}`

**Auth:** Yes

**Rules**

- Only owner
- Cannot edit inactive recipe

---

### DELETE `/recipes/{id}`

**Auth:** Yes

**Rules**

- Soft delete (deactivate)

---

## Interactions (S4)

### POST `/interactions`

**Auth:** Optional

**Request**

```json
{
  "recipe_id": "uuid",
  "type": "view | like | save"
}
```

**Response 200**

```json
{ "status": "ok" }
```

**Rules**

- Idempotent per user + recipe + type
- Anonymous users allowed (no user_id)

---

## Personalized Feed (S5)

### GET `/feed/for-you`

**Auth:** Yes

**Query Params**

- `page`
- `limit`

**Response 200**

```json
{
  "items": [Recipe],
  "page": 1,
  "has_more": true
}
```

**Rules**

- Excludes inactive recipes
- Cold-start fallback applied

---

## AI Chef (S6)

### POST `/api/ai/recipes`

**Auth:** Yes

**Content-Type:** `application/json`

**Request**

```json
{
  "ingredients": ["string"],
  "constraints": {
    "max_time": 30,
    "difficulty": "easy",
    "cuisine": "italian"
  }
}
```

**Response (SSE streaming):**

```text
event: catalog_results
data: {"source":"catalog","recipes":[...],"match_count":3}

event: ai_generation
data: {"source":"ai","status":"starting"}

event: ai_generation
data: {"source":"ai","status":"partial","token":"..."}

event: ai_generation
data: {"source":"ai","status":"complete","recipe":{...}}

event: error
data: {"error":"...","code":"RATE_LIMITED"|"AI_TIMEOUT"|"INVALID_INPUT"|"INTERNAL_ERROR"}
```

**Rules**

- Prefer existing recipes (≥ 3 matches → return catalog)
- AI generation only when < 3 catalog matches
- Generated recipes persisted with `ai_generated=true`
- Rate limited: 10 requests/user/day
- SSE streaming by default; `Accept: application/json` for blocking response

### GET `/api/ai/recipes/{id}`

**Auth:** Optional

**Response 200**

```json
{
  "id": "string",
  "name_recipe": "string",
  "description": "string",
  "img_url": "string",
  "ingredients": [{ "name": "string", "quantity": 1.5, "unit": "cups" }],
  "steps": [{ "order": 1, "instruction": "string" }],
  "estimated_minutes": 25,
  "difficulty": "easy",
  "cuisine": "italian",
  "ai_generated": true,
  "created_at": "iso-date"
}
```

### POST `/api/ai/recipes/{id}/feedback`

**Auth:** Yes

**Request**

```json
{"action": "save" | "discard"}
```

**Response 200**

```json
{ "status": "ok" }
```

---

## Meal Types

### GET `/meal-types`

**Auth:** No

**Response 200**

```json
[{ "id": 1, "name": "Breakfast" }]
```

---

## Shared Types

### Recipe

```json
{
  "id": "uuid",
  "title": "string",
  "ingredients": ["string"],
  "steps": ["string"],
  "meal_type": { "id": 1, "name": "string" },
  "ai_generated": false,
  "active": true,
  "created_at": "iso-date"
}
```

---

## Error Format (Global)

```json
{
  "error": "string"
}
```

---

## Contract Rules for Agents

- Do not change response shapes without updating this file
- Client must not infer missing fields
- Server must not return extra fields
