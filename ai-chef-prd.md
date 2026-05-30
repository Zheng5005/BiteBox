# AI Chef Integration — PRD

**Product:** BiteBox | **Phase:** 3 | **Version:** 1.0 · Draft

---

## Overview

The AI Chef allows authenticated users to input available ingredients and receive personalized recipe suggestions. The system first queries the existing recipe catalog; only when results are insufficient does it invoke the AI generation pipeline. Generated recipes are persisted and reused to reduce API costs and build the catalog organically.

> **Core principle: prefer existing recipes, generate only when needed.** The AI is a catalog-building engine, not a real-time oracle.

---

## Processing Flow

| Step | Name | Description |
|------|------|-------------|
| 01 | User input | Ingredients + optional constraints |
| 02 | Catalog search | Match by ingredient tags |
| 03 | Threshold check | ≥ 3 results? Serve them |
| 04 | AI generation | Prompt → structured recipe |
| 05 | Persist & serve | Stored with `ai_generated = true` |

> Steps 4–5 are only executed when step 3 returns fewer than 3 results.

---

## User-Facing Features

- **Ingredient input** — Free-text or tag-based entry. Users list what they have — no required quantities.
- **Optional constraints** — Max time (minutes), difficulty level, cuisine style. All optional, surfaced as quick filters.
- **Recipe results** — Ranked list with match score. AI-generated recipes show a subtle "AI" badge.
- **Chef conversation** — Users can ask follow-up questions ("make it vegetarian", "what if I don't have X").

---

## Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-01 | **Ingredient parsing** — Backend must normalize free-text ingredient input into canonical tags before querying the catalog. | P1 |
| FR-02 | **Catalog search** — Query must support partial ingredient matching (recipe uses a subset of provided ingredients). Minimum match: 2 of N provided. | P1 |
| FR-03 | **Threshold gate** — AI generation is triggered only when catalog search returns fewer than 3 relevant results. Threshold is configurable per environment. | P1 |
| FR-04 | **AI generation** — The AI module must produce a fully structured recipe (title, ingredients with quantities, steps, tags, estimated time, difficulty). | P1 |
| FR-05 | **Deduplication** — Before persisting an AI recipe, check for semantic near-duplicates by title similarity and ingredient overlap. Reuse existing if match score > 0.8. | P1 |
| FR-06 | **Persistence** — All AI-generated recipes are saved with `ai_generated = true`, the source prompt, and generation timestamp. | P1 |
| FR-07 | **Constraint filtering** — Constraints (time, difficulty, style) are applied at catalog search stage first, then passed as prompt instructions to AI generation. | P2 |
| FR-08 | **Interaction tracking** — AI Chef usage must emit interaction events (ingredients submitted, results viewed, recipe selected) into the existing signals pipeline. | P1 |
| FR-09 | **Follow-up chat** — Chef session supports at least one round of follow-up (e.g. "substitute X with Y"). Follow-up context is appended to the original prompt. | P2 |
| FR-10 | **AI badge labeling** — Any AI-generated recipe shown in results or detail view must display a visible "AI" label. No deception about provenance. | P1 |

---

## AI Output Schema

The generation pipeline must return a strictly typed structure. Loose prose responses are not acceptable — use function calling or structured output mode.

```json
{
  "title": "string",
  "description": "string (1–2 sentences)",
  "ingredients": [{ "name": "string", "quantity": "number", "unit": "string" }],
  "steps": [{ "order": "number", "instruction": "string" }],
  "tags": ["string"],
  "estimated_minutes": "number",
  "difficulty": "easy | medium | hard",
  "cuisine": "string"
}
```

---

## Backend Module Responsibilities

| Module | Responsibility |
|--------|---------------|
| `/internal/ai` | Prompt construction, LLM call, response parsing, validation against schema |
| `/internal/recipes` | Ingredient normalization, catalog search query, deduplication logic, persistence |
| `/internal/feed` | Surfacing AI Chef results in the For You feed when a recipe was generated or used |
| `/internal/recommendations` | Consuming AI Chef interaction events as personalization signals |

---

## Non-Functional Requirements

| Metric | Target |
|--------|--------|
| P95 latency — catalog hit | < 800ms |
| P95 latency — AI generation | < 6s |
| Min catalog hit rate at launch | 60% |
| AI call rate limit | 10 / user / day |

AI generation must stream output to the client. The UI must show a loading state and begin rendering as tokens arrive — a full 6s blocking spinner is not acceptable.

---

## Out of Scope (v1)

- Image generation for AI recipes
- Nutritional estimation
- Shopping list / missing ingredient detection
- Multi-turn agentic chef (more than 1 follow-up round)
- Fine-tuned or self-hosted model

---

## Success Metrics

| ID | Metric |
|----|--------|
| M-01 | Catalog hit rate ≥ 60% within 30 days of launch (measures catalog growth from AI persistence) |
| M-02 | Chef → Recipe save rate ≥ 15% (user found a result worth keeping) |
| M-03 | 7-day return rate for Chef feature users ≥ 35% |
| M-04 | Avg. AI recipes generated per new unique ingredient combination < 1.2 (dedup is working) |
