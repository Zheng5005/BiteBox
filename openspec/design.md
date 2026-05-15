# AI Chef — Design

## Architecture Overview

```
Client (React)
    │
    │ POST /api/ai/recipes (SSE)
    ▼
┌─────────────────────────────┐
│  AIChefHandler               │
│  (handlers/aichef/)          │
│  ┌───────────────────────┐  │
│  │ 1. Parse & validate   │  │
│  │ 2. Rate limit check   │  │
│  │ 3. Normalize tags     │  │
│  │ 4. Catalog search     │  │
│  │ 5. Threshold gate     │  │
│  │ 6. AI call / dedup    │  │
│  │ 7. Persist & track    │  │
│  └───────────────────────┘  │
├─────────────────────────────┤
│  internal/ai/                │
│  - Prompt builder           │
│  - LLM caller (Gemini)      │
│  - Response validator       │
│  - Streaming SSE writer     │
├─────────────────────────────┤
│  internal/recipes/ (new)    │
│  - Tag normalization        │
│  - Catalog query            │
│  - Deduplication            │
│  - AI recipe persistence    │
├─────────────────────────────┤
│  internal/ratelimit/         │
│  - In-memory limiter        │
│  - TTL cleanup goroutine    │
├─────────────────────────────┤
│  internal/interactions/      │
│  - AI interaction events    │
│  - Feedback recording       │
└─────────────────────────────┘
    │
    ▼
PostgreSQL (recipes, recipe_ingredients, ai_chef_interactions)
```

## Package Structure

```
Server/
├── main.go                          # New route registration
├── internal/
│   ├── ai/
│   │   ├── client.go                # Gemini API HTTP client
│   │   ├── prompt.go                # System prompt + context builder
│   │   ├── parser.go                # JSON parsing & validation
│   │   └── types.go                 # AiRecipe, AiIngredient, AiStep
│   ├── recipes/
│   │   ├── catalog.go               # Catalog search query
│   │   ├── dedup.go                 # Levenshtein + Jaccard dedup
│   │   ├── tags.go                  # Ingredient → tag normalization
│   │   ├── persist.go               # Save AI recipe to DB
│   │   └── types.go                 # Repository types
│   ├── ratelimit/
│   │   └── limiter.go               # In-memory rate limiter
│   └── interactions/
│       └── tracker.go               # AI interaction event recording
├── handlers/
│   └── aichef/
│       ├── handler.go               # Main SSE handler
│       └── types.go                 # Request/response types
└── db/
    └── migrations/
        └── 001_ai_chef.sql          # Schema migration
```

## Data Flow

### Catalog Hit Path (fast)

```
Client → AIChefHandler → normalizeIngredients()
  → catalogSearch() → match_count >= 3
  → recordInteraction(catalog, match_count)
  → SSE: catalog_results event → close
```

### AI Generation Path (slow)

```
Client → AIChefHandler → normalizeIngredients()
  → catalogSearch() → match_count < 3
  → rateLimiter.Check(userID) → if exceeded → SSE: error(RATE_LIMITED)
  → SSE: ai_generation(status: "starting")
  → aiClient.Generate(prompt) → stream tokens
  → SSE: ai_generation(status: "partial", token: "...")
  → parser.Validate(parsed) → if invalid → retry once
  → dedup.Check(candidate) → if match > 0.8 → SSE: ai_generation(existing recipe)
  → persist.Save(candidate) → SSE: ai_generation(status: "complete", recipe: {...})
  → recordInteraction(ai, 1, recipe.ID)
  → close
```

## Module Details

### 1. `internal/ai/client.go` — Gemini HTTP Client

```go
type AIClient struct {
    apiKey     string
    model      string
    httpClient *http.Client
}

func (c *AIClient) GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error)
func (c *AIClient) Generate(ctx context.Context, prompt string) (*AiRecipe, error)
```

- Uses `net/http` directly — no external SDK
- Calls `https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent`
- Function calling schema passed in request for structured output
- Stream channel yields raw JSON tokens for SSE forwarding

### 2. `internal/ai/prompt.go` — Prompt Construction

```go
func BuildPrompt(ingredients []string, constraints *AiConstraints, catalogHints []string) string
```

System prompt enforces:

- Strict JSON output matching AiRecipe schema
- No markdown code fences
- Quantities must be numbers, units must be standard
- Steps must be sequential with order field
- Use only provided ingredients (no extras)

When catalogHints exist (partial matches), prompt instructs: "These similar recipes already exist. Create something different."

### 3. `internal/recipes/tags.go` — Tag Normalization

```go
func NormalizeIngredients(raw []string) []string
```

- Lowercase, trim whitespace
- Remove common stop words ("a", "an", "some", "fresh")
- Singularize plurals (simple suffix rules: "ies"→"y", "s"→"")
- Deduplicate
- Result used as tags for DB query

### 4. `internal/recipes/catalog.go` — Catalog Search

```go
type CatalogRepo struct {
    DB db.DBExecutor
}

func (r *CatalogRepo) SearchByIngredients(ctx context.Context, tags []string, constraints *AiConstraints) ([]RecipeSummary, error)
```

- Uses the SQL query from spec (array intersection with ≥ 2 matches)
- Applies constraint filters as additional WHERE clauses
- Returns sorted by match_score DESC

### 5. `internal/recipes/dedup.go` — Deduplication

```go
func (r *CatalogRepo) FindNearDuplicate(ctx context.Context, candidate *AiRecipe) (*RecipeSummary, error)
```

- Scans active recipes with matching tags
- Computes combined score (0.6 title + 0.4 ingredients)
- Returns existing recipe if score > 0.8

### 6. `internal/recipes/persist.go` — Persistence

```go
func (r *CatalogRepo) SaveAiRecipe(ctx context.Context, recipe *AiRecipe, prompt string) (int, error)
```

- Inserts into `recipes` with `ai_generated = true`, `source_prompt`, `generated_at`
- Inserts each ingredient into `recipe_ingredients`
- Inserts steps as JSONB into `steps_json`
- Populates `tags` from normalized ingredients
- Returns recipe ID

### 7. `internal/ratelimit/limiter.go` — Rate Limiter

```go
type RateLimiter struct {
    mu       sync.RWMutex
    requests map[string][]time.Time
    limit    int
    ttl      time.Duration
}

func NewRateLimiter(limit int, ttl time.Duration) *RateLimiter
func (rl *RateLimiter) Allow(userID string) bool
func (rl *RateLimiter) StartCleanup(interval time.Duration)
```

- Key format: `user:{id}`
- Stores timestamps, prunes entries older than TTL
- Cleanup goroutine runs every hour

### 8. `internal/interactions/tracker.go` — Interaction Tracking

```go
type InteractionTracker struct {
    DB db.DBExecutor
}

func (t *InteractionTracker) RecordChefInteraction(userID int, ingredients []string, resultsCount int, selectedRecipeID *int, wasAIGenerated bool) error
func (t *InteractionTracker) RecordFeedback(userID int, recipeID int, action string) error
```

- Lightweight INSERT-only operations
- Feedback endpoint records save/discard actions

### 9. `handlers/aichef/handler.go` — SSE Handler

```go
type AIChefHandler struct {
    CatalogRepo     *recipes.CatalogRepo
    AIClient        *ai.AIClient
    RateLimiter     *ratelimit.RateLimiter
    InteractionTrk  *interactions.InteractionTracker
    SecretKey       string
    Threshold       int
}

func (h *AIChefHandler) AiChefHandler(w http.ResponseWriter, r *http.Request)
func (h *AIChefHandler) AiRecipeDetailHandler(w http.ResponseWriter, r *http.Request)
func (h *AIChefHandler) AiFeedbackHandler(w http.ResponseWriter, r *http.Request)
```

SSE handler sets headers:

```go
w.Header().Set("Content-Type", "text/event-stream")
w.Header().Set("Cache-Control", "no-cache")
w.Header().Set("Connection", "keep-alive")
w.Header().Set("X-Accel-Buffering", "no") // Nginx passthrough
```

Helper function for writing events:

```go
func writeSSE(w http.ResponseWriter, event string, data interface{})
```

## Frontend Architecture

### New Files

```
Client/src/
├── api/
│   └── aiChef.ts              # SSE streaming API service
├── pages/
│   └── AIChef.tsx             # Main AI Chef page
├── components/
│   ├── AIRecipeCard.tsx       # Recipe card with AI badge
│   ├── ChefInput.tsx          # Ingredient input with tag chips
│   ├── ChefLoading.tsx        # Streaming loading animation
│   └── ConstraintFilters.tsx  # Time/difficulty/cuisine filters
└── hooks/
    └── useSSE.ts              # Generic SSE hook
```

### SSE Streaming Hook

```typescript
function useSSE(url: string, body: object) {
  const [events, setEvents] = useState<SSEEvent[]>([]);
  const [status, setStatus] = useState<
    "idle" | "streaming" | "complete" | "error"
  >("idle");

  // Returns event stream, handles parsing, connection lifecycle
  return { events, status, error, send };
}
```

### AIChef Page Flow

1. User enters ingredients (comma-separated or tag chips)
2. Optionally sets constraints (quick filter buttons)
3. Clicks "Cook!" → POST to `/api/ai/recipes` via SSE
4. If catalog results → immediate display with match scores
5. If AI generation → streaming animation, progressive render
6. Results shown with AI badge on generated recipes
7. User can save/discard (triggers feedback API)

### AI Badge Component

```typescript
// Small pill badge overlaid on recipe cards
<AIBadge /> // Renders "✨ AI" pill, subtle opacity
```

## Error Handling Strategy

| Layer      | Error               | Client Response                               |
| ---------- | ------------------- | --------------------------------------------- |
| Validation | Empty ingredients   | Show inline error                             |
| Rate limit | 429 RATE_LIMITED    | "Try again tomorrow" banner                   |
| AI timeout | 504 AI_TIMEOUT      | "Chef is busy, retry" + show catalog fallback |
| Network    | Connection lost     | Retry button                                  |
| Parse      | Invalid AI response | "Invalid response, retry"                     |

## main.go Changes

```go
// New imports
import (
    "github.com/Zheng5005/BiteBox/handlers/aichef"
    "github.com/Zheng5005/BiteBox/internal/ai"
    "github.com/Zheng5005/BiteBox/internal/recipes"
    "github.com/Zheng5005/BiteBox/internal/ratelimit"
    "github.com/Zheng5005/BiteBox/internal/interactions"
)

// In main():
aiClient := ai.NewAIClient(os.Getenv("GEMINI_API_KEY"), "gemini-2.0-flash")
catalogRepo := recipes.NewCatalogRepo(db.DB)
rateLimiter := ratelimit.NewRateLimiter(10, 24*time.Hour)
rateLimiter.StartCleanup(1 * time.Hour)
interactionTracker := interactions.NewInteractionTracker(db.DB)

aiChefHandler := aichef.NewAIChefHandler(catalogRepo, aiClient, rateLimiter, interactionTracker, secret, 3)

mux.HandleFunc("POST /api/ai/recipes", middleware.JWTMiddleware(aiChefHandler.AiChefHandler))
mux.HandleFunc("GET /api/ai/recipes/", aiChefHandler.AiRecipeDetailHandler)
mux.HandleFunc("POST /api/ai/recipes/feedback", middleware.JWTMiddleware(aiChefHandler.AiFeedbackHandler))
```
