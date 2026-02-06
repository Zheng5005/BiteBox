# Server – 🔴 MVP Micro-Checklists

Purpose: break **high-risk server tasks** into atomic, deterministic steps that an agent can execute without ambiguity.

Scope: **S4, S5, S6 only**.

---

## S4. Interaction Tracking API 🔴 (Core)

### Checklist
- [ ] Confirm interaction table does not already exist
- [ ] Define interaction types enum (view, like, save)
- [ ] Create DB table with:
  - [ ] user_id (nullable for anon)
  - [ ] recipe_id (FK)
  - [ ] interaction_type
  - [ ] created_at
- [ ] Add unique constraint / dedup strategy
- [ ] Create POST `/api/interactions`
- [ ] Validate payload
- [ ] Reject invalid interaction types
- [ ] Handle anonymous users explicitly
- [ ] Implement dedup logic (same user + recipe + type)
- [ ] Return success without blocking client

### Failure Conditions
- Duplicate rows inflate counts
- Invalid interaction types stored

### Done When
- Repeated identical events do not create duplicates
- Valid interactions persist reliably

---

## S5. Personalized Feed (Rule-Based) 🔴

### Checklist
- [ ] Identify interaction weights (view < like < save)
- [ ] Aggregate interactions per user
- [ ] Join interactions → recipes
- [ ] Compute score per recipe
- [ ] Order by score DESC
- [ ] Add recency bias
- [ ] Exclude inactive recipes
- [ ] Limit + paginate results
- [ ] Implement cold-start fallback

### Failure Conditions
- Non-deterministic ordering
- Slow queries on large datasets

### Done When
- Same input produces same order
- Query time acceptable for MVP

---

## S6. AI Chef (Text Only) 🔴

### Checklist
- [ ] Accept ingredient list input
- [ ] Normalize ingredients
- [ ] Search existing recipes first
- [ ] Define minimum match threshold
- [ ] If insufficient results → call AI
- [ ] Build deterministic prompt
- [ ] Parse AI output into:
  - [ ] title
  - [ ] ingredients
  - [ ] steps
  - [ ] tags
- [ ] Validate parsed output
- [ ] Persist recipe with `ai_generated=true`
- [ ] Prevent duplicate recipes

### Failure Conditions
- Hallucinated ingredients stored
- Malformed recipes persisted

### Done When
- AI-generated recipes match schema
- Stored recipes usable like normal ones

---

## Agent Rules (Reminder)

- Follow checklist top → bottom
- Do not skip unchecked items
- Stop immediately when Done When is true

