# Proposal: Recommendation Engine

## Intent
Increase user engagement and discovery by providing personalized and popularity-based recipe recommendations on the main page.

## Scope

### In Scope
- Database schema update: `recipe_views` table.
- Backend: `GET /recipes/popular` (weighted score).
- Backend: `GET /recipes/recommended` (personalized by user's top categories).
- Backend: Background/Async view increment on recipe details fetch.
- Frontend: `Main.tsx` update to show "Most Popular" and "Recommended for You" sections.

### Out of Scope
- Advanced Machine Learning models.
- Real-time event streaming (Kafka/Redpanda).
- Complex collaborative filtering (User-User similarity).

## Approach
Use a weighted scoring algorithm in SQL: `Score = (Views * 1) + (Likes * 5) + (AvgRating * 10)`.
For personalization, identify the user's top 3 most viewed/liked `meal_type_id`s and fetch popular recipes from those types.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `Server/db/bitebox.sql` | Modified | Add `recipe_views` table. |
| `Server/handlers/recipes` | Modified | Add sorting logic and new endpoints. |
| `Client/src/api/recipes.ts` | Modified | Add `getPopular` and `getRecommended`. |
| `Client/src/pages/Main.tsx` | Modified | Replace flat list with categorized feeds. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Query Performance | Medium | Index `recipe_views(recipe_id)` and `recipe_views(user_id)`. |
| Privacy Concerns | Low | Only internal tracking; no external sharing. |

## Rollback Plan
1. Revert frontend to use `getRecipes` (flat list).
2. Remove new endpoints.
3. Drop `recipe_views` table (optional, data loss acceptable for this feature).

## Success Criteria
- [ ] Views are accurately incremented when a user opens a recipe.
- [ ] "Most Popular" feed displays recipes with high interaction counts.
- [ ] "Recommended" feed shows recipes matching the user's preferred meal types.
