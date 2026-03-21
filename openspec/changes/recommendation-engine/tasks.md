# Tasks: Recommendation Engine

## Phase 1: Database & Core Logic
- [ ] 1.1 Create `Server/db/dev/recipe_views.sql` with table definition and indexes.
- [ ] 1.2 Run SQL command to create `recipe_views` table in local database.
- [ ] 1.3 Add `incrementViews(recipeID int, userID *int)` helper to `Server/handlers/recipes/RecipesHandlers.go`.
- [ ] 1.4 Update `RecipeONEHandler` in `Server/handlers/recipes/RecipesHandlers.go` to call `incrementViews` in a goroutine.

## Phase 2: Recommendation Endpoints
- [ ] 2.1 Implement `GetPopularRecipes` handler in `Server/handlers/recipes/RecipesHandlers.go` using weighted score query.
- [ ] 2.2 Implement `GetRecommendedRecipes` handler in `Server/handlers/recipes/RecipesHandlers.go` using personalization logic.
- [ ] 2.3 Register new routes (`/api/recipes/popular`, `/api/recipes/recommended`) in `Server/main.go`.

## Phase 3: Frontend Integration
- [ ] 3.1 Update `Client/src/api/recipes.ts` to include `getPopularRecipes` and `getRecommendedRecipes`.
- [ ] 3.2 Create `Client/src/components/RecipeList.tsx` to display a horizontal/grid list of recipes (refactor from `Main.tsx` if needed).
- [ ] 3.3 Update `Client/src/pages/Main.tsx` to fetch `popular` and `recommended` feeds in parallel and display them using `RecipeList`.

## Phase 4: Verification
- [ ] 4.1 Verify `recipe_views` table exists in PostgreSQL.
- [ ] 4.2 Test: Access a recipe detail page and verify `recipe_views` count increases.
- [ ] 4.3 Test: Verify `/api/recipes/popular` returns recipes sorted by score.
- [ ] 4.4 Test: Verify `Main.tsx` renders "Popular" and "Recommended" sections correctly.
