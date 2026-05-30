# SDD Proposal — AI Chef Integration

## Problem Statement

Users want to input ingredients they have on hand and receive recipe suggestions. The system should prefer returning existing catalog recipes and only invoke AI generation when the catalog lacks sufficient matches. Generated recipes must be persisted to grow the catalog organically and reduce future AI API costs.

## Scope

**In scope:**

- Backend: new AI Chef module (`Server/internal/ai`), new handler (`Server/handlers/aichef`), DB schema migration, catalog search by ingredients, AI generation pipeline, deduplication, rate limiting, interaction tracking.
- Frontend: new AI Chef page, API service module, AI badge on recipe cards, streaming response UI.
- API contract: `POST /api/ai/recipes` endpoint with streaming support.

**Out of scope (v1 per PRD):**

- Image generation for AI recipes
- Nutritional estimation
- Shopping list / missing ingredient detection
- Multi-turn follow-up beyond 1 round
- Fine-tuned or self-hosted model

## Solution Options

### Option A: Monolithic endpoint (Recommended)

Single `POST /api/ai/recipes` endpoint that:

1. Normalizes ingredients
2. Queries catalog with ingredient matching
3. If ≥ 3 results, returns them
4. If < 3, calls AI API, deduplicates, persists, and returns

**Pros:** Simple, single round-trip, easy to rate-limit, clear flow.
**Cons:** Longer response time when AI is triggered; needs streaming to stay under P95 6s target.

### Option B: Two-phase endpoint

`POST /api/ai/recipes/search` → returns catalog results + `needs_ai` flag.
`POST /api/ai/recipes/generate` → generates if needed.

**Pros:** Faster initial response, client controls flow.
**Cons:** More complex client logic, two API calls in worst case, harder to track interaction signals atomically.

### Recommendation: Option A

Monolithic endpoint with Server-Sent Events (SSE) streaming. Catalog hits return immediately. AI generation streams structured tokens as they arrive. Matches the PRD's requirement for streaming output.

## Key Design Decisions

1. **Ingredient matching**: Use PostgreSQL array containment operators (`@>`, `&&`) on a new `tags` column in recipes. Simpler and faster than FTS for this use case.

2. **AI provider**: Gemini 2.0 Flash via REST API with function calling for structured output. No SDK dependency needed — direct HTTP call keeps the Go codebase lean.

3. **Deduplication**: Simple Levenshtein distance on title + Jaccard similarity on ingredient names. Score > 0.8 returns existing recipe.

4. **Rate limiting**: In-memory map with TTL per user ID. 10 requests/user/day. Sufficient for v1; can move to Redis later.

5. **Interaction tracking**: New `ai_chef_interactions` table. Records ingredients submitted, results returned, recipe selected. Consumed later by recommendations module.

## Database Changes

### New columns on `recipes` table

- `ai_generated BOOLEAN DEFAULT FALSE`
- `source_prompt TEXT`
- `generated_at TIMESTAMP`
- `estimated_minutes INT`
- `difficulty VARCHAR(20)` — 'easy', 'medium', 'hard'
- `cuisine VARCHAR(50)`
- `tags TEXT[]` — array of ingredient/category tags

### New table: `recipe_ingredients`

- `id SERIAL PRIMARY KEY`
- `recipe_id INT REFERENCES recipes(id)`
- `ingredient_name TEXT`
- `quantity NUMERIC`
- `unit VARCHAR(20)`

### New table: `ai_chef_interactions`

- `id SERIAL PRIMARY KEY`
- `user_id INT` (nullable for anonymous)
- `ingredients_submitted TEXT[]`
- `results_count INT`
- `recipe_selected_id INT` (nullable)
- `was_ai_generated BOOLEAN`
- `created_at TIMESTAMP DEFAULT NOW()`

## Risk Assessment

| Risk                                     | Impact | Mitigation                                                   |
| ---------------------------------------- | ------ | ------------------------------------------------------------ |
| AI API latency exceeds 6s P95            | High   | Streaming SSE, timeout at 8s, fallback to "try again later"  |
| Recipe table grows large without indexes | Medium | Add GIN index on `tags` column                               |
| Rate limiting leaks memory               | Medium | TTL-based cleanup goroutine, max map size                    |
| Deduplication false positives            | Low    | Threshold of 0.8 is conservative; manual review option later |

## Artifacts to Produce

1. `spec.md` — Detailed API contract, data models, error codes
2. `design.md` — Module architecture, Go package structure, data flow
3. `tasks.md` — Ordered implementation tasks with acceptance criteria

## Next Phase: Spec

Recommend proceeding to spec to define the exact API request/response shapes, error codes, and database migration SQL.
