# Recommendation Engine Specification

## Domain: Recommendation

### Requirement: Track Recipe Views
The system MUST track when a user or guest views a recipe's detail page.

#### Scenario: Authenticated User View
- GIVEN an authenticated user
- WHEN they request `GET /api/recipes/:id`
- THEN the system MUST record a view for that recipe with the user's ID
- AND the total view count for the recipe MUST increment by 1

#### Scenario: Guest View
- GIVEN a guest (unauthenticated) user
- WHEN they request `GET /api/recipes/:id`
- THEN the system MUST record a view for that recipe (user_id is NULL)
- AND the total view count for the recipe MUST increment by 1

### Requirement: List Popular Recipes
The system MUST provide a list of recipes sorted by popularity.
Popularity Score = (Views * 1) + (Likes * 5) + (AvgRating * 10).

#### Scenario: Fetch Popular Feed
- GIVEN a database with recipes having varying views, likes, and ratings
- WHEN a user requests `GET /api/recipes/popular`
- THEN the system RETURNS a list of recipes sorted by Popularity Score descending
- AND the list contains at most 10 items (or pagination limit)

#### Scenario: No Interaction Data
- GIVEN a fresh database with no views/likes
- WHEN a user requests `GET /api/recipes/popular`
- THEN the system RETURNS recipes sorted by ID or creation date (fallback)

### Requirement: Personalized Recommendations
The system SHOULD provide personalized recommendations based on the user's interaction history (views/likes) with meal types.

#### Scenario: User with History
- GIVEN a user who has viewed/liked mostly "Italian" and "Mexican" recipes
- WHEN they request `GET /api/recipes/recommended`
- THEN the system IDENTIFIES "Italian" and "Mexican" as top preferences
- AND RETURNS popular recipes specifically from those meal types
- AND EXCLUDES recipes the user has already viewed (optional, but good for discovery)

#### Scenario: User without History (Cold Start)
- GIVEN a new user with no interactions
- WHEN they request `GET /api/recipes/recommended`
- THEN the system RETURNS the generic "Popular" list
