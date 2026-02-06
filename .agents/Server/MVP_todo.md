# Server MVP TODO (Agent-Optimized)

Goal: implement **deterministic, testable backend logic** for personalization and AI foundations.

Legend: 🟢 Low | 🟡 Medium | 🔴 High

---

## S1. Auth Middleware 🟡

**Steps**

1. Validate JWT signature
2. Enforce expiration
3. Apply middleware to all protected routes

**Done When**

* Invalid tokens rejected
* No unprotected endpoints

---

## S2. Recipe API Hardening 🟡

**Steps**

1. Validate create/update inputs
2. Enforce ownership checks
3. Respect active/inactive state

**Edge Cases**

* Editing inactive recipes

**Done When**

* Invalid operations rejected

---

## S3. Meal Types 🟢

**Steps**

1. Load meal types at startup
2. Cache in memory
3. Validate FK usage on recipes

**Done When**

* Stable responses

---

## S4. Interaction Tracking API 🔴 (Core)

**Steps**
See @MICROCHECKLIST.md

**Edge Cases**

* Anonymous users
* High-frequency events

---

## S5. Personalized Feed (Rule-Based) 🔴

**Steps**
See @MICROCHECKLIST.md

**Done When**

* Deterministic ranking
* Acceptable query time

---

## S6. AI Chef (Text Only) 🔴

**Steps**
See @MICROCHECKLIST.md

**Edge Cases**

* Hallucinated ingredients
* Duplicate recipes

---

## S7. Database Integrity 🟡

**Steps**

1. Add FK constraints
2. Index feed-related queries
3. Verify soft-delete consistency

**Done When**

* No orphaned data
* Feed queries performant

