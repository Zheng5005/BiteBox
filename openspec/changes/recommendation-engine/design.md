# Design: Recommendation Engine

## Technical Approach
Implement popularity and personalization logic directly in the PostgreSQL database using aggregation queries. Use a new `recipe_views` table to track engagement. The frontend will fetch "Popular" and "Recommended" feeds via new API endpoints.

## Architecture Decisions

### Decision: Score Calculation Location
**Choice**: SQL Aggregation
**Alternatives considered**: Application-level calculation
**Rationale**: SQL is faster for aggregation and sorting large datasets. Reduces data transfer.

### Decision: View Tracking
**Choice**: Asynchronous (Goroutine)
**Alternatives considered**: Synchronous (Blocking), Message Queue
**Rationale**: Incrementing views should not delay the recipe details response. A simple goroutine is sufficient for current scale; MQ is overkill.

### Decision: Personalization Strategy
**Choice**: Category-based (Meal Types)
**Alternatives considered**: User-User Collaborative Filtering
**Rationale**: Collaborative filtering requires sparse matrix operations, complex to implement in SQL. Category affinity is a strong enough signal for MVP.

## Data Flow

1.  **View Tracking**:
    User -> `GET /recipes/:id` -> `RecipeONEHandler` -> (Spawn Goroutine -> `INSERT INTO recipe_views`) -> Return Recipe

2.  **Popular Feed**:
    User -> `GET /recipes/popular` -> `GetPopularHandler` -> `SELECT ... ORDER BY score DESC` -> Return JSON

3.  **Recommended Feed**:
    User -> `GET /recipes/recommended` -> `GetRecommendedHandler` ->
    a. `SELECT meal_type_id FROM recipe_views/likes WHERE user_id = ? GROUP BY meal_type_id ORDER BY count DESC LIMIT 3`
    b. `SELECT * FROM recipes WHERE meal_type_id IN (...) ORDER BY score DESC` -> Return JSON

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `Server/db/bitebox.sql` | Modify | Add `recipe_views` table definition. |
| `Server/handlers/recipes/RecipesHandlers.go` | Modify | Add `GetPopularRecipes`, `GetRecommendedRecipes`, `incrementViews`. |
| `Server/handlers/recipes/types.go` | Modify | Add `RecipeScore` struct (if needed). |
| `Client/src/api/recipes.ts` | Modify | Add `getPopularRecipes`, `getRecommendedRecipes`. |
| `Client/src/pages/Main.tsx` | Modify | Refactor to use `RecipeList` and fetch multiple feeds. |
| `Client/src/components/RecipeList.tsx` | Create | Reusable component for displaying a list/grid of recipes. |

## Interfaces / Contracts

### API: Get Popular Recipes
`GET /api/recipes/popular`
Response: `[RecipeMainPage] (same as existing)`

### API: Get Recommended Recipes
`GET /api/recipes/recommended`
Response: `[RecipeMainPage]`

### DB: Recipe Views
```sql
CREATE TABLE recipe_views (
    id SERIAL PRIMARY KEY,
    recipe_id INT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_views_recipe ON recipe_views(recipe_id);
CREATE INDEX idx_views_user ON recipe_views(user_id);
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit (Go) | `incrementViews` | Mock DB, verify INSERT call. |
| Integration (Go) | `GetPopularRecipes` | Seed DB with views/likes, verify order of results. |
| Integration (Go) | `GetRecommendedRecipes` | Seed user history, verify returned categories match. |
| Component (React) | `RecipeList` | Render with mock data, verify clicks/layout. |
| E2E | Full Flow | User clicks recipe -> View count up -> Appears in Popular. |

## Migration / Rollout
Run SQL command to create `recipe_views` table. No data migration needed (starts empty).

## Open Questions
- None.
