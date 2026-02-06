## Product Vision & Scope
BiteBox aims to be a **social recipe platform with personalization and AI assistance**. Conceptually, it combines:
- A recipe discovery app
- A personalized "For You" recommendation feed
- An AI-powered personal chef
- Social interactions around food content

The long-term goal is to evolve from a simple recipe finder into a learning system that adapts to each user’s tastes and cooking habits.

---
## Core Pillars

### 1. Recipe Platform
- Recipes are first-class content
- Users can:
  - Create recipes
  - View recipe details
  - Like and save recipes
  - (Later) comment and share
### 2. Personalization Engine ("For You" Feed)
- Available only to authenticated users
- Driven by user behavior and interaction history
- Starts rule-based, evolves toward collaborative filtering

Tracked signals include:
- Recipe views
- Likes
- Saves
- Cooking attempts (future)
### 3. AI Personalized Chef
- Users provide available ingredients
- Backend follows a **hybrid strategy**:
  1. Search for existing recipes that match ingredients
  2. If insufficient results, generate a new recipe using AI
- AI-generated recipes are stored and reused
### 4. Social Layer (Phase 2+)
- Follow other users
- User profiles
- Engagement-driven feeds
---
## MVP Scope (Initial Release)

The MVP focuses on validating personalization and retention.
### Included
- User authentication
- Recipe CRUD
- Basic feed (chronological or popular)
- Like / Save interactions
- Interaction tracking for personalization
- AI Chef v1 (text-only recipe generation)

### Excluded (for MVP)
- Followers
- Notifications
- Chat
- Advanced ML models
---
## Personalization Strategy
### Phase 1 – Rule-Based Scoring
- Implemented with SQL + Go
- Recipe score based on:
  - Matching liked tags
  - Similar saved recipes
  - Cuisine preferences
  - Interaction recency
### Phase 2 – Collaborative Filtering
- "Users who liked X also liked Y"
- Similar user profiles
### Phase 3 – ML / Embeddings (Optional)
- Vector similarity for recipes
- AI-assisted ranking
---
## AI Chef Design
### Input
- Ingredients list
- Optional constraints (time, difficulty, style)
### Output
- Structured recipe:
  - Title
  - Ingredients
  - Steps
  - Tags
### Processing Rules
- Prefer existing recipes
- Generate only when needed
- Mark generated recipes as `ai_generated = true`
---
## Backend Architecture Guidelines
Logical modules (single service, modular design):

```
/internal
  /auth
  /recipes
  /feed
  /recommendations
  /ai
```

Responsibilities:
- Auth: JWT, users
- Recipes: CRUD, search
- Feed: pagination, ranking
- Recommendations: scoring logic
- AI: prompt building and response parsing
---
## Frontend Architecture Guidelines
Key screens:
- Auth (login / register)
- For You feed
- Recipe detail
- Create recipe
- AI Chef

State management:
- Auth: global
- Feed & recipes: server-driven
- Avoid heavy global state until necessary
---
## Metrics & Signals
All user interactions should be persisted:
- Views
- Likes
- Saves
- AI usage
- Return frequency

These metrics directly power personalization.
---
## Development Roadmap
### Phase 1 – Foundation
- Auth
- Recipes CRUD
- Basic feed
- Interaction tracking
### Phase 2 – Smart Feed
- Rule-based "For You" algorithm
- Personalized ranking
### Phase 3 – AI Chef
- Ingredient-based input
- Hybrid recipe resolution
- AI recipe persistence
### Phase 4 – Social Features
- Follow system
- Comments
- Profile pages
---
## Guiding Principles for Agents

- Prefer incremental improvements over large rewrites
- Always store interaction data
- Favor deterministic logic before ML
- Treat AI as a system component, not a black box
- Ship working features early and iterate

