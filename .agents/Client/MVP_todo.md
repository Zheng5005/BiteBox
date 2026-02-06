# Client MVP TODO (Agent-Optimized)

Goal: implement a **stable MVP client** that supports auth, recipe interaction, and personalization signals.

Legend: 🟢 Low | 🟡 Medium | 🔴 High

---

## C1. Auth Session Handling 🟡

**Steps**

1. Create `AuthContext` (user, token, status)
2. Store JWT in memory (optional: localStorage restore)
3. Attach token via fetch/axios interceptor
4. Handle `401` globally → logout + redirect

**Edge Cases**

* Expired token on reload
* Backend unavailable

**Done When**

* Protected routes redirect
* No silent auth failures

---

## C2. Recipe Feed (MVP) 🟡

**Steps**

1. Create feed API client with pagination params
2. Render list with loading + empty states
3. Implement load-more or infinite scroll

**Edge Cases**

* Empty feed
* Partial page fetch

**Done When**

* Pagination stable
* UI survives fetch errors

---

## C3. Recipe Detail 🟡

**Steps**

1. Add route `/recipes/:id`
2. Fetch recipe by ID
3. Render ingredients, steps, meal type
4. Add like/save actions

**Edge Cases**

* Invalid recipe ID
* Inactive recipe

**Done When**

* Correct error states
* Interaction state synced

---

## C4. Recipe Create / Edit 🔴

**Steps**

1. Build recipe form component
2. Add client-side validation
3. Meal type selector
4. Submit with optimistic UI

**Edge Cases**

* Missing required fields
* Duplicate ingredients

**Done When**

* Invalid forms blocked
* Server errors displayed

---

## C5. Interaction Tracking 🟡 (Critical)

**Steps**

1. Define interaction event types
2. Fire events on view / like / save
3. Send asynchronously (non-blocking)

**Edge Cases**

* Rapid repeated events
* Tracking API down

**Done When**

* Events sent reliably
* UI unaffected by failures

---

## C6. For You Feed UI 🟢

**Steps**

1. Add `/for-you` route
2. Gate by auth
3. Fallback to default feed

**Done When**

* Route accessible only to auth users

