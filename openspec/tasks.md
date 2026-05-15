# AI Chef — Tasks

## Phase 1: Database & Foundation

### Task 1.1: Database Migration

**Estimate:** 30 lines
**Files:** `Server/db/migrations/001_ai_chef.sql`

**Acceptance Criteria:**

- [ ] Migration SQL creates all new columns, tables, and indexes from spec
- [ ] Migration runs cleanly on existing database (IF NOT EXISTS guards)
- [ ] All foreign keys reference correct existing tables

---

### Task 1.2: Internal Types

**Estimate:** 80 lines
**Files:** `Server/internal/ai/types.go`, `Server/internal/recipes/types.go`, `Server/handlers/aichef/types.go`

**Acceptance Criteria:**

- [ ] `AiRecipe`, `AiIngredient`, `AiStep` types with JSON tags
- [ ] `AiChefRequest`, `AiConstraints` request types
- [ ] `CatalogResultEvent`, `AiGenerationEvent`, `ErrorEvent` SSE types
- [ ] `RecipeSummary` type for catalog results
- [ ] All types compile cleanly

---

## Phase 2: Core Backend Modules

### Task 2.1: Tag Normalization

**Estimate:** 50 lines
**Files:** `Server/internal/recipes/tags.go`, `Server/internal/recipes/tags_test.go`

**Acceptance Criteria:**

- [ ] `NormalizeIngredients([]string) []string` function
- [ ] Lowercase, trim, stop word removal, basic singularization
- [ ] Deduplication of normalized tags
- [ ] Unit tests for normalization edge cases

---

### Task 2.2: Rate Limiter

**Estimate:** 70 lines
**Files:** `Server/internal/ratelimit/limiter.go`, `Server/internal/ratelimit/limiter_test.go`

**Acceptance Criteria:**

- [ ] `RateLimiter` struct with map + mutex
- [ ] `Allow(userID string) bool` method
- [ ] TTL-based cleanup goroutine
- [ ] Thread-safe (concurrent access test)
- [ ] 10 requests/user/day default

---

### Task 2.3: Catalog Search

**Estimate:** 60 lines
**Files:** `Server/internal/recipes/catalog.go`, `Server/internal/recipes/catalog_test.go`

**Acceptance Criteria:**

- [ ] `CatalogRepo` struct with DBExecutor
- [ ] `SearchByIngredients` method using array intersection SQL
- [ ] Minimum 2 tag match threshold
- [ ] Constraint filters (max_time, difficulty, cuisine) as optional WHERE clauses
- [ ] Returns sorted by match_score DESC
- [ ] Tests with mock DB

---

### Task 2.4: Deduplication

**Estimate:** 80 lines
**Files:** `Server/internal/recipes/dedup.go`, `Server/internal/recipes/dedup_test.go`

**Acceptance Criteria:**

- [ ] Levenshtein distance function for title similarity
- [ ] Jaccard similarity function for ingredient overlap
- [ ] Combined score: 0.6 _ title + 0.4 _ ingredients
- [ ] `FindNearDuplicate` returns existing recipe if score > 0.8
- [ ] Tests for exact match, near match, and no match cases

---

### Task 2.5: AI Recipe Persistence

**Estimate:** 80 lines
**Files:** `Server/internal/recipes/persist.go`, `Server/internal/recipes/persist_test.go`

**Acceptance Criteria:**

- [ ] `SaveAiRecipe` inserts into `recipes` with ai_generated=true
- [ ] Ingredients inserted into `recipe_ingredients` table
- [ ] Steps stored as JSONB in `steps_json`
- [ ] Tags populated from normalized ingredients
- [ ] Returns recipe ID
- [ ] Tests with mock DB

---

### Task 2.6: AI Client (Gemini)

**Estimate:** 120 lines
**Files:** `Server/internal/ai/client.go`, `Server/internal/ai/client_test.go`

**Acceptance Criteria:**

- [ ] `AIClient` struct with API key and model
- [ ] `Generate(ctx, prompt) (*AiRecipe, error)` for non-streaming
- [ ] `GenerateStream(ctx, prompt) (<-chan string, <-chan error)` for streaming
- [ ] Gemini REST API call with function calling schema
- [ ] Proper error handling and timeout (8s context)
- [ ] Tests mock HTTP responses

---

### Task 2.7: Prompt Builder & Parser

**Estimate:** 100 lines
**Files:** `Server/internal/ai/prompt.go`, `Server/internal/ai/parser.go`

**Acceptance Criteria:**

- [ ] `BuildPrompt(ingredients, constraints, catalogHints) string`
- [ ] System prompt enforces strict JSON output
- [ ] Catalog hints included when partial matches exist
- [ ] `ValidateAndParse(rawJSON) (*AiRecipe, error)` validates against schema
- [ ] Tests for prompt construction and validation

---

### Task 2.8: Interaction Tracker

**Estimate:** 40 lines
**Files:** `Server/internal/interactions/tracker.go`

**Acceptance Criteria:**

- [ ] `RecordChefInteraction` INSERT into `ai_chef_interactions`
- [ ] `RecordFeedback` for save/discard actions
- [ ] Nullable user_id for anonymous

---

## Phase 3: Handler & Routing

### Task 3.1: AI Chef SSE Handler

**Estimate:** 150 lines
**Files:** `Server/handlers/aichef/handler.go`

**Acceptance Criteria:**

- [ ] `AiChefHandler` with all dependencies injected
- [ ] SSE headers set correctly
- [ ] Full flow: validate → rate limit → normalize → catalog search → threshold gate → AI or return
- [ ] `writeSSE` helper for event formatting
- [ ] JSON fallback when `Accept: application/json`
- [ ] Error events for all error codes
- [ ] Timeout handling (8s context)

---

### Task 3.2: Recipe Detail & Feedback Handlers

**Estimate:** 60 lines
**Files:** `Server/handlers/aichef/handler.go` (continued)

**Acceptance Criteria:**

- [ ] `AiRecipeDetailHandler` — GET with ingredients array from `recipe_ingredients`
- [ ] `AiFeedbackHandler` — POST save/discard, calls interaction tracker
- [ ] Both handlers return correct error codes

---

### Task 3.3: main.go Route Registration

**Estimate:** 20 lines
**Files:** `Server/main.go`

**Acceptance Criteria:**

- [ ] All new internal packages imported
- [ ] AI client, catalog repo, rate limiter, interaction tracker initialized
- [ ] Rate limiter cleanup goroutine started
- [ ] Three new routes registered with correct middleware
- [ ] Server compiles and starts

---

## Phase 4: Frontend

### Task 4.1: API Service Module

**Estimate:** 60 lines
**Files:** `Client/src/api/aiChef.ts`, `Client/src/hooks/useSSE.ts`

**Acceptance Criteria:**

- [ ] `useSSE` hook manages EventSource-like connection (using fetch + ReadableStream for POST)
- [ ] `submitAiChefRequest(ingredients, constraints)` returns event stream
- [ ] `submitAiFeedback(recipeId, action)` for feedback
- [ ] TypeScript types for all events
- [ ] Error handling with retry logic

---

### Task 4.2: AI Chef Page

**Estimate:** 200 lines
**Files:** `Client/src/pages/AIChef.tsx`

**Acceptance Criteria:**

- [ ] Ingredient input (comma-separated + tag chips)
- [ ] Constraint filter buttons (time, difficulty, cuisine)
- [ ] "Cook!" submit button
- [ ] Results display: catalog results or AI-generated recipe
- [ ] Streaming loading animation during AI generation
- [ ] Save/discard buttons on AI results
- [ ] Error states: rate limit, timeout, network error
- [ ] Responsive layout

---

### Task 4.3: AI Badge Component

**Estimate:** 30 lines
**Files:** `Client/src/components/AIBadge.tsx`

**Acceptance Criteria:**

- [ ] Small "✨ AI" pill badge
- [ ] Reusable, accepts className prop
- [ ] Used on AI-generated recipe cards and detail view

---

### Task 4.4: ChefInput & ChefLoading Components

**Estimate:** 120 lines
**Files:** `Client/src/components/ChefInput.tsx`, `Client/src/components/ChefLoading.tsx`, `Client/src/components/ConstraintFilters.tsx`

**Acceptance Criteria:**

- [ ] `ChefInput`: text input → tag chips, backspace to remove
- [ ] `ChefLoading`: animated "Chef is cooking..." with streaming indicator
- [ ] `ConstraintFilters`: quick-select buttons for time/difficulty/cuisine
- [ ] All components are accessible (keyboard, screen reader labels)

---

### Task 4.5: Route Integration

**Estimate:** 10 lines
**Files:** `Client/src/App.tsx` (or router file)

**Acceptance Criteria:**

- [ ] `/ai-chef` route added to React Router
- [ ] Route protected (requires auth) or available to all
- [ ] Navigation link added to main nav/menu

---

## Phase 5: Verification

### Task 5.1: End-to-End Backend Test

**Estimate:** 80 lines
**Files:** `Server/handlers/aichef/handler_test.go`

**Acceptance Criteria:**

- [ ] Integration test: catalog hit path returns SSE catalog_results
- [ ] Integration test: AI generation path returns SSE ai_generation
- [ ] Rate limiting test: 11th request returns 429
- [ ] Deduplication test: near-duplicate returns existing recipe
- [ ] Mock AI client for deterministic tests

---

### Task 5.2: API Contract Update

**Estimate:** 20 lines
**Files:** `.agents/API_CONTRACTS.md`

**Acceptance Criteria:**

- [ ] AI Chef section updated with final request/response shapes
- [ ] New error codes documented
- [ ] SSE event format documented
