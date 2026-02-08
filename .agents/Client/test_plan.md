# Test Plan for BiteBox Frontend

This document provides a detailed testing blueprint for Jules. It is based on the edge case analysis and is intended to guide the creation of unit and integration tests.

## App.tsx: Routing

---

### **Test Case ID:** `APP-001`

-   **Component/File:** `App.tsx`
-   **Scenario:** A non-authenticated user attempts to access the `/profile` route, which requires authentication.
-   **Test Steps:**
    1.  Render the `App` component within a test environment that supports routing (e.g., using `MemoryRouter` from `react-router-dom`).
    2.  Set the initial route to `/profile`.
    3.  Ensure that the `AuthContext` provides a `null` user, simulating a non-authenticated session.
-   **Expected Result:** The application should not render the `Profile` page. Instead, it should redirect the user to the `/login` page, and the `Login` component should be visible.
-   **Mocking/Setup:**
    -   Wrap the `App` component in `MemoryRouter`.
    -   Provide a mock `AuthContext` where `user` is `null`.

---

### **Test Case ID:** `APP-002`

-   **Component/File:** `App.tsx`
-   **Scenario:** A user navigates to a URL that does not match any defined routes.
-   **Test Steps:**
    1.  Render the `App` component within a `MemoryRouter`.
    2.  Set the initial route to a non-existent path, such as `/this/route/does/not/exist`.
-   **Expected Result:** A "Not Found" component should be rendered. If no "Not Found" page exists, the test should verify that none of the existing page components (like `MainPage`, `Login`, `Profile`, etc.) are rendered.
-   **Mocking/Setup:**
    -   Wrap the `App` component in `MemoryRouter`.

## API Layer

---

### **Test Case ID:** `API-001`

-   **Component/File:** `src/api/axiosInstance.ts` (and its usage in a component like `Profile.tsx`)
-   **Scenario:** A logged-in user's JWT has expired. They try to perform an action that makes an API call.
-   **Test Steps:**
    1.  Mock the `axios` instance to simulate a `401 Unauthorized` response for an API call (e.g., `getUserRecipes`).
    2.  Render a component that makes this API call on load (e.g., `Profile.tsx`).
    3.  Simulate a logged-in user by providing a (now invalid) token in the `AuthContext`.
-   **Expected Result:** The application should handle the `401` error gracefully. The `AuthContext`'s `logout` function should be called, clearing the user's session and redirecting them to the `/login` page.
-   **Mocking/Setup:**
    -   Mock `axios` or the specific API function (e.g., `getUserRecipes`) to return a `401` error.
    -   Use `MemoryRouter` to check for redirection.

---

### **Test Case ID:** `API-002`

-   **Component/File:** `src/api/users.ts`
-   **Scenario:** The `/users` API endpoint returns a recipe object with `id` and `rating` as strings instead of numbers.
-   **Test Steps:**
    1.  Directly test the `getUserRecipes` function.
    2.  Mock the `axiosInstance.get` call to return a response with an array of recipe objects where `id` and `rating` are strings (e.g., `id: "123"`, `rating: "4.5"`).
    3.  Call `getUserRecipes` and await its result.
-   **Expected Result:** The function should correctly parse the string values into numbers. The returned `Recipe[]` should have `id` and `rating` as numbers.
-   **Mocking/Setup:**
    -   Mock `axiosInstance.get`.

## Components

---

### **Test Case ID:** `COMP-001`

-   **Component/File:** `src/components/Avatar.tsx`
-   **Scenario:** The `Avatar` component is rendered for a user who has an account but has not uploaded a profile picture (`url_photo` is an empty string).
-   **Test Steps:**
    1.  Provide a mock `AuthContext` with a user object where `url_photo` is `""`.
    2.  Render the `Avatar` component with any `size`.
-   **Expected Result:** The `img` element's `src` attribute should point to the `avatar.iran.liara.run` stock image URL, not an empty string.
-   **Mocking/Setup:**
    -   Mock `useAuth` hook or `AuthContext.Provider`.

---

### **Test Case ID:** `COMP-002`

-   **Component/File:** `src/components/IFLButton.tsx`
-   **Scenario:** The "I feel lucky" button is clicked when there are no recipes available.
-   **Test Steps:**
    1.  Mock the `useNavigate` hook from `react-router`.
    2.  Render the `IFLButton` component, passing an empty array `[]` as the `recipes` prop.
    3.  Simulate a user click on the button.
-   **Expected Result:** The mocked `navigate` function should not be called.
-   **Mocking/Setup:**
    -   Mock `useNavigate` to provide a spy function (e.g., `jest.fn()`).

---

### **Test Case ID:** `COMP-003`

-   **Component/File:** `src/components/RecipeCard.tsx`
-   **Scenario:** A `RecipeCard` is displayed for a recipe that has not yet been rated (`rating` is 0).
-   **Test Steps:**
    1.  Create a mock `Recipe` object where `rating` is `0`.
    2.  Render the `RecipeCard` component with the mock recipe.
-   **Expected Result:** The component should display the text "BE THE FIRST ONE TO RATE IT!" instead of a star rating number.
-   **Mocking/Setup:**
    -   None required beyond creating the prop.

## Context

---

### **Test Case ID:** `CTX-001`

-   **Component/File:** `src/context/AuthContext.tsx`
-   **Scenario:** The application loads with a malformed JWT in `localStorage`.
-   **Test Steps:**
    1.  Before rendering, set `localStorage` with an invalid token (e.g., `localStorage.setItem("token", "this.is.not.a.valid.token")`).
    2.  Render a component that uses the `useAuth` hook, wrapped in the `AuthProvider`.
-   **Expected Result:** The `useEffect` in `AuthProvider` should catch the parsing error. The `logout` function should be called, `localStorage` should be cleared of the token, and the `user` state should be `null`.
-   **Mocking/Setup:**
    -   Manipulate `localStorage` before the test run.

## Hooks

---

### **Test Case ID:** `HOOK-001`

-   **Component/File:** `src/hooks/useMealTypes.ts`
-   **Scenario:** The API call to fetch meal types fails.
-   **Test Steps:**
    1.  Mock the `getMealTypes` API function to reject with an error.
    2.  Use a test component or `renderHook` from `@testing-library/react-hooks` to call the `useMealTypes` hook.
-   **Expected Result:** The hook should catch the error and return an empty array `[]`.
-   **Mocking/Setup:**
    -   Mock the `getMealTypes` function in `src/api/meals.ts`.

## Pages

---

### **Test Case ID:** `PAGE-001`

-   **Component/File:** `src/pages/Login.tsx`
-   **Scenario:** User attempts to log in with incorrect credentials.
-   **Test Steps:**
    1.  Mock the `login` API function to simulate a failure, returning an Axios error with a `401` status and a response body like `{ message: "Invalid credentials" }`.
    2.  Render the `Login` component.
    3.  Simulate user input for email and password.
    4.  Simulate form submission.
-   **Expected Result:** The `info.error` state should be updated with "Invalid credentials", and this message should be visible in the rendered component. The `isSubmiting` state should be set to `false` after the attempt.
-   **Mocking/Setup:**
    -   Mock the `login` function in `src/api/auth.ts`.

---

### **Test Case ID:** `PAGE-002`

-   **Component/File:** `src/pages/PostRecipe.tsx`
-   **Scenario:** A guest user (not logged in) submits the form without filling in the "By:" (`guest_name`) field.
-   **Test Steps:**
    1.  Provide a mock `AuthContext` where `user` is `null`.
    2.  Render the `PostRecipe` component.
    3.  Fill in all other required fields (name, description, etc.) but leave `guest_name` empty.
    4.  Simulate form submission.
-   **Expected Result:** The browser's default HTML5 validation should prevent submission because the `guest_name` input is `required`. No API call should be made.
-   **Mocking/Setup:**
    -   Mock `useAuth` to return `{ user: null }`.

---

### **Test Case ID:** `PAGE-003`

-   **Component/File:** `src/pages/PostRecipe.tsx`
-   **Scenario:** User tries to upload a `.txt` file instead of an image.
-   **Test Steps:**
    1.  Render the `PostRecipe` component.
    2.  Find the `input` with `type="file"`.
    3.  Simulate a change event with a mock file object that has a `type` of `text/plain`.
-   **Expected Result:** The `info.error` state should be updated with the message "Only PNG or JPEG images". The `previewImage` and `imageFile` states should remain `null`.
-   **Mocking/Setup:**
    -   Create a mock `File` object for the test.

---

### **Test Case ID:** `PAGE-004`

-   **Component/File:** `src/pages/Profile.tsx`
-   **Scenario:** A user visits their profile but has not posted any recipes.
-   **Test Steps:**
    1.  Mock the `getUserRecipes` API function to return a promise that resolves to an empty array `[]`.
    2.  Render the `Profile` component.
-   **Expected Result:** The component should display the text "You haven't posted any recipes yet."
-   **Mocking/Setup:**
    -   Mock `getUserRecipes` in `src/api/users.ts`.

---

### **Test Case ID:** `PAGE-005`

-   **Component/File:** `src/pages/RecipeDetails.tsx`
-   **Scenario:** A user navigates to a recipe detail page with an ID that does not exist.
-   **Test Steps:**
    1.  Mock the `useParams` hook to return an invalid `id`.
    2.  Mock the `getRecipeById` function to throw an error or return a `404` status.
    3.  Render the `RecipeDetails` component.
-   **Expected Result:** The component should display an error message like "Failed to load recipe details." and not the loading text.
-   **Mocking/Setup:**
    -   Mock `useParams` from `react-router`.
    -   Mock `getRecipeById` in `src/api/recipes.ts`.

---

### **Test Case ID:** `PAGE-006`

-   **Component/File:** `src/pages/RecipeDetails.tsx`
-   **Scenario:** A logged-in user tries to submit a comment that contains only whitespace.
-   **Test Steps:**
    1.  Mock `useAuth` to return a valid user.
    2.  Render the `RecipeDetails` component for any valid recipe.
    3.  Navigate to the "Leave a comment" tab.
    4.  Simulate typing "   " (spaces) into the comment `textarea`.
    5.  Simulate submitting the comment form.
-   **Expected Result:** The `postComment` API function should not be called. The form might show a validation message, or simply do nothing.
-   **Mocking/Setup:**
    -   Mock `postComment` from `src/api/comments.ts` with a spy to verify it wasn't called.
    -   Mock `useAuth`.

---

### **Test Case ID:** `PAGE-007`

-   **Component/File:** `src/pages/SignUp.tsx`
-   **Scenario:** A user types a username containing a number.
-   **Test Steps:**
    1.  Render the `SignUp` component.
    2.  Simulate a user typing "User123" into the username input field.
-   **Expected Result:** As the user types the number, an error message "Numbers not allowed in username." should appear. The form submission should be blocked if the error is present.
-   **Mocking/Setup:**
    -   None required.
