# AI Chef — Spec

## API Contract

### POST `/api/ai/recipes`

**Auth:** Required (JWT Bearer)

**Content-Type:** `application/json`

**Request Body:**

```json
{
  "ingredients": ["string", ...],       // Required. 1-20 items. Free-text ingredients.
  "constraints": {                       // Optional. All fields optional.
    "max_time": 30,                      // int. Max cooking time in minutes.
    "difficulty": "easy",                // string. "easy" | "medium" | "hard"
    "cuisine": "italian"                 // string. Free-text cuisine preference.
  }
}
```

**Response (200) — Server-Sent Events stream:**

The response is streamed via SSE. Two event types:

#### Event: `catalog_results`

Emitted when catalog has ≥ 3 matches. Stream closes immediately after.

```
event: catalog_results
data: {
  "source": "catalog",
  "recipes": [RecipeSummary, ...],
  "match_count": 5
}
```

#### Event: `ai_generation`

Emitted when AI generation is triggered. Multiple `partial` events stream tokens, followed by one `complete` event.

```
event: ai_generation
data: {
  "source": "ai",
  "status": "partial",
  "token": "{\"title\":\"Chicken..."
}

event: ai_generation
data: {
  "source": "ai",
  "status": "complete",
  "recipe": {
    "id": "uuid",
    "title": "string",
    "description": "string",
    "ingredients": [
      { "name": "string", "quantity": "number", "unit": "string" }
    ],
    "steps": [
      { "order": 1, "instruction": "string" }
    ],
    "tags": ["string"],
    "estimated_minutes": 25,
    "difficulty": "easy",
    "cuisine": "string",
    "ai_generated": true,
    "created_at": "2026-05-15T..."
  }
}
```

#### Event: `error`

```
event: error
data: {
  "error": "string",
  "code": "RATE_LIMITED" | "AI_TIMEOUT" | "INVALID_INPUT" | "INTERNAL_ERROR"
}
```

**Error Codes:**

| Code             | HTTP Status | Description                                  |
| ---------------- | ----------- | -------------------------------------------- |
| `INVALID_INPUT`  | 400         | Missing ingredients, empty array, > 20 items |
| `RATE_LIMITED`   | 429         | User exceeded 10 AI requests/day             |
| `AI_TIMEOUT`     | 504         | AI provider did not respond within 8s        |
| `INTERNAL_ERROR` | 500         | Unexpected server error                      |

---

### GET `/api/ai/recipes/{id}`

**Auth:** Optional

**Purpose:** Retrieve full detail of an AI-generated recipe (same as GET /api/recipes/{id} but includes `ingredients` array with quantities/units).

**Response (200):**

```json
{
  "id": "string",
  "name_recipe": "string",
  "description": "string",
  "meal_type": { "id": 1, "name": "string" },
  "img_url": "string",
  "creator_name": "string",
  "rating": "string",
  "likes": 0,
  "ai_generated": true,
  "steps": [{ "order": 1, "instruction": "string" }],
  "ingredients": [{ "name": "string", "quantity": 1.5, "unit": "cups" }],
  "estimated_minutes": 25,
  "difficulty": "easy",
  "cuisine": "italian",
  "tags": ["string"],
  "created_at": "2026-05-15T..."
}
```

---

### POST `/api/ai/recipes/{id}/feedback`

**Auth:** Required

**Request:**

```json
{
  "action": "save" | "discard"
}
```

**Response (200):**

```json
{ "status": "ok" }
```

**Purpose:** Record whether user saved or discarded an AI-generated recipe. Feeds into interaction tracking.

---

## Data Models

### Go Types

```go
// Request
type AiChefRequest struct {
    Ingredients []string        `json:"ingredients"`
    Constraints *AiConstraints  `json:"constraints,omitempty"`
}

type AiConstraints struct {
    MaxTime    *int    `json:"max_time,omitempty"`
    Difficulty *string `json:"difficulty,omitempty"`
    Cuisine    *string `json:"cuisine,omitempty"`
}

// AI Output Schema (validated against before persistence)
type AiRecipe struct {
    Title            string                 `json:"title"`
    Description      string                 `json:"description"`
    Ingredients      []AiIngredient         `json:"ingredients"`
    Steps            []AiStep               `json:"steps"`
    Tags             []string               `json:"tags"`
    EstimatedMinutes int                    `json:"estimated_minutes"`
    Difficulty       string                 `json:"difficulty"`
    Cuisine          string                 `json:"cuisine"`
}

type AiIngredient struct {
    Name     string  `json:"name"`
    Quantity float64 `json:"quantity"`
    Unit     string  `json:"unit"`
}

type AiStep struct {
    Order       int    `json:"order"`
    Instruction string `json:"instruction"`
}

// SSE Event payloads
type CatalogResultEvent struct {
    Source     string          `json:"source"`
    Recipes    []RecipeSummary `json:"recipes"`
    MatchCount int             `json:"match_count"`
}

type AiGenerationEvent struct {
    Source string     `json:"source"`
    Status string     `json:"status"` // "partial" | "complete"
    Token  string     `json:"token,omitempty"`
    Recipe *RecipeDetailFull `json:"recipe,omitempty"`
}

type ErrorEvent struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}
```

### RecipeSummary (for catalog results)

```go
type RecipeSummary struct {
    ID          string `json:"id"`
    Name        string `json:"name_recipe"`
    Description string `json:"description"`
    ImgURL      string `json:"img_url"`
    Rating      string `json:"rating"`
    Likes       int    `json:"likes"`
    MatchScore  float64 `json:"match_score"`
}
```

---

## Database Schema Migration

### Migration SQL

```sql
-- 1. Add AI Chef columns to recipes table
ALTER TABLE recipes
    ADD COLUMN IF NOT EXISTS ai_generated BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS source_prompt TEXT,
    ADD COLUMN IF NOT EXISTS generated_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS estimated_minutes INT,
    ADD COLUMN IF NOT EXISTS difficulty VARCHAR(20),
    ADD COLUMN IF NOT EXISTS cuisine VARCHAR(50),
    ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}';

-- 2. Create recipe_ingredients table
CREATE TABLE IF NOT EXISTS recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_name TEXT NOT NULL,
    quantity NUMERIC(10,2),
    unit VARCHAR(20)
);

-- 3. Create ai_chef_interactions table
CREATE TABLE IF NOT EXISTS ai_chef_interactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    ingredients_submitted TEXT[],
    results_count INT DEFAULT 0,
    recipe_selected_id INT REFERENCES recipes(id) ON DELETE SET NULL,
    was_ai_generated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 4. Add steps_json column to recipes for structured steps
-- (existing steps column is plain text; new recipes use JSON)
ALTER TABLE recipes
    ADD COLUMN IF NOT EXISTS steps_json JSONB;

-- 5. GIN index on tags for fast array queries
CREATE INDEX IF NOT EXISTS idx_recipes_tags ON recipes USING GIN (tags);

-- 6. Index on ai_generated for filtering
CREATE INDEX IF NOT EXISTS idx_recipes_ai_generated ON recipes (ai_generated);
```

---

## Catalog Search Query

```sql
-- Find recipes matching at least 2 of N provided ingredient tags
SELECT
    r.id,
    r.name_recipe,
    r.description,
    r.img_url,
    COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS rating,
    COUNT(DISTINCT rl.id) AS likes,
    -- Match score: fraction of user ingredients covered by recipe tags
    ROUND(
        CAST(array_length(ARRAY(
            SELECT unnest(r.tags)
            INTERSECT
            SELECT unnest($1::TEXT[])
        ), 1) AS numeric)
        / CAST(array_length($1::TEXT[], 1) AS numeric),
        2
    ) AS match_score
FROM recipes r
LEFT JOIN comments c ON r.id = c.recipe_id
LEFT JOIN recipe_likes rl ON r.id = rl.recipe_id
WHERE r.is_active = true
    AND r.tags && $1::TEXT[]  -- At least one tag overlap
    AND array_length(ARRAY(
        SELECT unnest(r.tags)
        INTERSECT
        SELECT unnest($1::TEXT[])
    ), 1) >= 2  -- At least 2 ingredient matches
GROUP BY r.id
ORDER BY match_score DESC, r.id DESC
LIMIT 10;
```

---

## Deduplication Algorithm

```
Input: AiRecipe candidate, existing recipes

1. Title similarity: Levenshtein ratio(candidate.title, existing.title)
   - Normalized: 1.0 - (distance / max(len_a, len_b))

2. Ingredient overlap: Jaccard similarity
   - Intersection(candidate.ingredients, existing.ingredients) / Union(...)
   - Compare on lowercase ingredient_name only

3. Combined score: 0.6 * title_similarity + 0.4 * ingredient_overlap

4. If score > 0.8 for any existing recipe → return that existing recipe
   Else → persist new AI recipe
```

---

## Rate Limiting

```go
// In-memory rate limiter
type RateLimiter struct {
    mu       sync.RWMutex
    requests map[string][]time.Time // key: "user:{id}:{date}"
    limit    int
}

// Check returns true if request is allowed
func (rl *RateLimiter) Check(userID string) bool

// Cleanup removes expired entries (called by goroutine every hour)
func (rl *RateLimiter) Cleanup()
```

---

## SSE Event Format

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

event: catalog_results
data: {"source":"catalog","recipes":[...],"match_count":3}

event: ai_generation
data: {"source":"ai","status":"partial","token":"{\"title\":..."}

event: ai_generation
data: {"source":"ai","status":"complete","recipe":{...}}
```

For non-streaming clients (SSE not supported), a fallback `Accept: application/json` returns the same data as a single JSON response (blocking).

---

## Acceptance Criteria

1. ✅ Catalog search returns existing recipes when ≥ 3 matches found
2. ✅ AI generation triggered only when < 3 catalog matches
3. ✅ Generated recipes persisted with `ai_generated = true`
4. ✅ Deduplication prevents near-duplicate AI recipes
5. ✅ Rate limiting enforces 10 requests/user/day
6. ✅ Interaction events recorded in `ai_chef_interactions`
7. ✅ SSE streaming for AI generation responses
8. ✅ JSON fallback for non-streaming clients
9. ✅ P95 catalog hit < 800ms, P95 AI generation < 6s
10. ✅ AI badge visible on generated recipes in frontend
