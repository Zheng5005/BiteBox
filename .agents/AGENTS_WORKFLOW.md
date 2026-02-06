# AGENT_WORKFLOW.md

Purpose: define **how agents work in this repo** to deliver the MVP efficiently, with minimal tokens and maximal correctness.

---

## 0. Prime Directive

* Ship **small, deterministic increments**
* One task at a time
* Stop when acceptance criteria are met

---

## 1. Task Selection Rules

1. Pick **exactly one task** from:

   * `Client/TODO.md` **or** `Server/TODO.md`
2. Prefer 🔴 tasks first (core MVP risk)
3. Do **not** start a new task until the current one is complete

---

## 2. Execution Protocol (Strict)

For a chosen task:

1. Read the task section fully
2. Follow **Steps** in order
3. Handle listed **Edge Cases** only (no extras)
4. Implement until **Done When** is true

❌ Do not:

* Add features not listed
* Refactor unrelated code
* Expand scope

---

## 3. File Ownership Boundaries

### Client Agent

* `Client/` only
* No backend assumptions
* Treat API as a black box

### Server Agent

* `Server/` only
* No frontend changes
* No UI concerns

---

## 4. Testing Rules

* Test **only what the task affects**
* Prefer deterministic tests
* If tests are absent:

  * Add minimal coverage
  * Do not introduce new frameworks

---

## 5. Error Handling Policy

* Fail fast
* Return explicit errors
* Never swallow errors silently

---

## 6. Commit Rules

Each task = **one commit**

Commit message format:

```
<scope>: <task-id> short description
```

Examples:

* `client: C5 interaction tracking events`
* `server: S4 interaction tracking endpoint`

---

## 7. Stop Conditions

Stop immediately when:

* “Done When” criteria are satisfied
* Tests pass
* Build succeeds

Do **not** continue to the next task automatically.

---

## 8. When Blocked

If blocked:

1. Stop work
2. Document the blocker clearly
3. Do not guess or workaround

---

## 9. MVP Completion Check

MVP is complete when:

* Personalized feed works
* Interactions are tracked
* AI can generate + persist recipes
* App survives partial failures

---

## 10. Guiding Principles

* Deterministic > clever
* Simple > flexible
* Shipped > perfect

