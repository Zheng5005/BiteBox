# BiteBox – MVP To‑Do List

This document defines **concrete, actionable tasks** required to complete the MVP of BiteBox.

The focus is on **correctness, testability, edge‑case handling, and incremental delivery**.

---

## CLIENT (React + TypeScript)

### 1. Authentication & Session Handling

**Tasks**
- Persist JWT securely (memory + refresh on reload if applicable)
- Centralized auth context/provider
- Automatic token attachment to API requests
- Graceful logout on 401 / token expiration

**Edge Cases**
- Expired or malformed token
- Backend unavailable during auth
- Multiple tabs open (stale auth state)

**Tests**
- Login success/failure flows
- Protected route access without token
- Token invalidation behavior

**Difficulty**: 🟡 Medium

---

### 2. Recipe Feed (MVP Feed)

**Tasks**
- Fetch paginated recipes (chronological or popularity-based)
- Infinite scroll or load-more button
- Loading & empty states
- Error fallback UI

**Edge Cases**
- Empty recipe list
- Partial page loads
- Network interruptions

**Tests**
- Pagination correctness
- Rendering performance for large lists
- Error state rendering

**Difficulty**: 🟡 Medium

---

### 3. Recipe Detail Page

**Tasks**
- Fetch recipe by ID
- Render ingredients, steps, meal type
- Show active/inactive status correctly
- Like & save actions

**Edge Cases**
- Recipe not found
- Inactive recipe access
- Unauthorized interaction attempts

**Tests**
- Route handling for invalid IDs
- Interaction state sync (liked/saved)

**Difficulty**: 🟡 Medium

---

### 4. Recipe Creation & Editing

**Tasks**
- Create recipe form
- Ingredient & step validation
- Meal type selection
- Submit + optimistic UI

**Edge Cases**
- Missing required fields
- Duplicate ingredients
- Very long input text

**Tests**
- Form validation rules
- Successful creation flow
- Server-side validation errors

**Difficulty**: 🔴 High

---

### 5. Interaction Tracking (Critical for MVP)

**Tasks**
- Fire events for:
  - View
  - Like
  - Save
- Ensure events are non-blocking

**Edge Cases**
- Duplicate events
- Rapid repeated interactions

**Tests**
- Event firing consistency
- No UI degradation if tracking fails

**Difficulty**: 🟡 Medium

---

### 6. Basic Personalization UI (For You)

**Tasks**
- Separate "For You" feed route
- Fallback to default feed if unauthenticated
- Visual distinction from general feed

**Edge Cases**
- New users with no history

**Tests**
- Auth vs unauth feed behavior

**Difficulty**: 🟢 Low

---

### 7. Testing & Tooling

**Tasks**
- Component tests for core screens
- API mocking
- Linting enforcement

**Difficulty**: 🟡 Medium

---

## SERVER (Go + PostgreSQL)

### 1. Auth Hardening

**Tasks**
- JWT expiration handling
- Middleware coverage audit
- Token refresh strategy (optional MVP+)

**Edge Cases**
- Token reuse
- Invalid signature

**Tests**
- Auth middleware tests

**Difficulty**: 🟡 Medium

---

### 2. Recipe API Stabilization

**Tasks**
- Enforce active/inactive logic consistently
- Input validation for create/update
- Ownership checks

**Edge Cases**
- Editing inactive recipes
- Unauthorized deletes

**Tests**
- Permission tests
- Validation failures

**Difficulty**: 🟡 Medium

---

### 3. Meal Types Endpoint

**Tasks**
- Cache meal types in memory
- Validate foreign key usage

**Edge Cases**
- Missing meal types

**Tests**
- Response consistency

**Difficulty**: 🟢 Low

---

### 4. Interaction Tracking API (MVP Critical)

**Tasks**
- Create interaction table (if not present)
- POST interaction endpoint
- Deduplicate repeated interactions

**Edge Cases**
- Anonymous users
- High-frequency events

**Tests**
- Idempotency tests
- Load tests (basic)

**Difficulty**: 🔴 High

---

### 5. Personalized Feed Logic (Rule-Based)

**Tasks**
- Aggregate user interactions
- Implement scoring function
- Rank recipes per user

**Edge Cases**
- Cold-start users
- Sparse interaction data

**Tests**
- Deterministic ranking tests
- SQL performance checks

**Difficulty**: 🔴 High

---

### 6. AI Chef – Generation Endpoint (Text Only)

**Tasks**
- Ingredient-based recipe search
- Fallback to AI generation
- Persist AI-generated recipes

**Edge Cases**
- Hallucinated ingredients
- Duplicate recipes

**Tests**
- AI response parsing
- Data integrity checks

**Difficulty**: 🔴 High

---

### 7. Database Integrity

**Tasks**
- Foreign key enforcement
- Indexing critical queries
- Soft-delete consistency

**Tests**
- Migration safety
- Constraint violation handling

**Difficulty**: 🟡 Medium

---

## MVP Exit Criteria

The MVP is considered complete when:
- Auth users receive a personalized feed
- Interactions are tracked reliably
- AI Chef can generate and persist recipes
- The app remains stable under partial failures

---

## Guiding Rules

- Prefer correctness over features
- Every interaction must be tracked
- No blocking UI for analytics
- Ship small, verify often

