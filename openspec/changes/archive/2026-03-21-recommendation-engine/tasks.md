# Tasks: Recommendation Engine

## Phase 1: Database & Core Logic
- [x] 1.1 Create `Server/db/dev/recipe_views.sql` with table definition and indexes.
- [x] 1.2 Run SQL command to create `recipe_views` table in local database.
- [x] 1.3 Add `incrementViews(recipeID int, userID *int)` helper to `Server/handlers/recipes/RecipesHandlers.go`.
- [x] 1.4 Update `RecipeONEHandler` in `Server/handlers/recipes/RecipesHandlers.go` to call `incrementViews` in a goroutine.

## Phase 2: Recommendation Endpoints
- [x] 2.1 Implement `GetPopularRecipes` handler in `Server/handlers/recipes/RecipesHandlers.go` using weighted score query.
- [x] 2.2 Implement `GetRecommendedRecipes` handler in `Server/handlers/recipes/RecipesHandlers.go` using personalization logic.
- [x] 2.3 Register new routes (`/api/recipes/popular`, `/api/recipes/recommended`) in `Server/main.go`.

## Phase 3: Frontend Integration
- [x] 3.1 Update `Client/src/api/recipes.ts` to include `getPopularRecipes` and `getRecommendedRecipes`.
- [x] 3.2 Create `Client/src/components/RecipeList.tsx` to display a horizontal/grid list of recipes (refactor from `Main.tsx` if needed).
- [x] 3.3 Update `Client/src/pages/Main.tsx` to fetch `popular` and `recommended` feeds in parallel and display them using `RecipeList`.

## Phase 4: Verification
- [x] 4.1 Verify `recipe_views` table exists in PostgreSQL.
- [x] 4.2 Test: Access a recipe detail page and verify `recipe_views` count increases.
- [x] 4.3 Test: Verify `/api/recipes/popular` returns recipes sorted by score.
- [x] 4.4 Test: Verify `Main.tsx` renders "Popular" and "Recommended" sections correctly.
