## Exploration: Recommendation Engine

### Current State
- **Database**: Recipes, Users, Comments (ratings), RecipeLikes (likes). No view tracking.
- **Backend**: `/recipes` fetches all active recipes, joined with likes/comments, unsorted.
- **Frontend**: `Main.tsx` displays a flat list of recipes.

### Affected Areas
- `Server/db/bitebox.sql`: New `recipe_views` table needed.
- `Server/handlers/recipes/RecipesHandlers.go`: Update `GetRecipes` (add sorting), create `GetPopularRecipes`, `GetPersonalizedRecipes`. Update `GetRecipe` (increment views).
- `Client/src/api/recipes.ts`: Add `getPopularRecipes`, `getPersonalizedRecipes`.
- `Client/src/pages/Main.tsx`: Update UI to show multiple sections/feeds.

### Approaches
1.  **Direct Query with Weighted Score (Recommended)**
    - Calculate score on the fly: `(views * 1) + (likes * 5) + (avg_rating * 10)`.
    - Personalization: Fetch user's top 3 meal types (based on likes/views) and recommend trending recipes in those types.
    - Pros: Real-time, simple implementation.
    - Cons: Expensive query as data grows.
    - Effort: Medium.

2.  **Cached Recommendations**
    - Pre-calculate scores/recommendations daily/hourly.
    - Pros: Fast reads.
    - Cons: Stale data, requires background worker/cron.
    - Effort: High.

### Recommendation
Proceed with **Approach 1**. It fits the current scale and provides immediate value without complex infrastructure. We can optimize later with caching or denormalization (counter columns).

### Risks
- **Performance**: `COUNT` on large `recipe_views` table will slow down over time. Mitigation: Add indexes on `recipe_id` and `user_id`.
- **Privacy**: Tracking views per user requires care (though standard practice).

### Ready for Proposal
Yes. The scope is clear: new table, new endpoints, frontend updates.
