# Edge Case and Test Analysis for BiteBox Frontend

This document outlines potential edge cases and suggests beneficial tests for the BiteBox frontend application. The analysis is based on the files in the `src` directory.

## App.tsx

### Routing

-   **Edge Case**: A user tries to access an authenticated route (e.g., `/profile`) without being logged in.
    -   **Test**: A test should ensure that if a non-authenticated user navigates to `/profile`, they are redirected to the login page or shown a "not authorized" message.
-   **Edge Case**: A user navigates to a route that doesn't exist (e.g., `/some/random/path`).
    -   **Test**: A test should verify that a "Not Found" or similar fallback page/component is rendered for undefined routes.

## API Layer (`src/api/`)

### `axiosInstance.ts`

-   **Edge Case**: The JWT token stored in `localStorage` is expired or invalid. When a request is made, the API will likely return a 401 or 403 error.
    -   **Test**: Mock an API call that returns a 401 error. The test should verify that the application handles this gracefully, for example, by logging the user out and redirecting to the login page.

### `users.ts` (`getUserRecipes`)

-   **Edge Case**: The API returns user recipes with unexpected data types (e.g., `id` as a string instead of a number, or a `rating` that is not a valid number).
    -   **Test**: Create a test that mocks the API response with malformed data and asserts that the `getUserRecipes` function either handles the transformation correctly (e.g., attempts to parse the number) or throws an appropriate error, preventing malformed data from propagating through the app.

## Components (`src/components/`)

### `Avatar.tsx`

-   **Edge Case**: The `user` object is not null, but the `url_photo` is an empty string.
    -   **Test**: Render the `Avatar` component with a user object where `url_photo` is `""`. The test should assert that the `src` of the rendered `img` tag falls back to the `stockImage` URL.

### `IFLButton.tsx`

-   **Edge Case**: The `recipes` prop is an empty array.
    -   **Test**: Render the `IFLButton` component with an empty array for the `recipes` prop. The test should simulate a click and assert that `navigate` is not called.

### `RecipeCard.tsx`

-   **Edge Case**: The `recipe.rating` is 0 or a negative number.
    -   **Test**: Render the `RecipeCard` with a recipe where the rating is 0. The test should check that the text "BE THE FIRST ONE TO RATE IT!" is displayed instead of a number.

## Context (`src/context/`)

### `AuthContext.tsx`

-   **Edge Case**: The JWT token in `localStorage` is malformed and cannot be parsed by `atob` or `JSON.parse`.
    -   **Test**: Set a malformed token in `localStorage` before rendering a component that uses the `AuthProvider`. The test should verify that the `logout` function is called and the user is not authenticated.

## Hooks (`src/hooks/`)

### `useMealTypes.ts`

-   **Edge Case**: The API call in `getMealTypes` fails (e.g., network error, server error).
    -   **Test**: Mock the `getMealTypes` API call to reject with an error. The test should assert that the `useMealTypes` hook returns an empty array `[]` for `mealTypes`.

## Pages (`src/pages/`)

### `Login.tsx`

-   **Edge Case**: The user submits the form with incorrect credentials, and the server returns an error message.
    -   **Test**: A test should simulate filling out the form and submitting it. Mock the `login` API call to return an Axios error with a response containing a `message`. The test should assert that the error message is displayed to the user in the UI.

### `PostRecipe.tsx`

-   **Edge Case**: A guest user (not logged in) tries to submit a recipe without providing a `guest_name`.
    -   **Test**: Render the `PostRecipe` component with the `user` from `useAuth` as `null`. A test should verify that the `guest_name` input is required, and the form cannot be submitted without it.
-   **Edge Case**: The user tries to upload a file that is not an image or is larger than the 2MB limit.
    -   **Test**: Simulate a file input change event with a non-image file type or a file larger than 2MB. The test should assert that the appropriate error message is displayed and the `imageFile` state is not set.

### `Profile.tsx`

-   **Edge Case**: A user has not posted any recipes yet.
    -   **Test**: Mock the `getUserRecipes` API call to return an empty array. The test should render the `Profile` component and assert that the message "You haven't posted any recipes yet." is displayed.

### `RecipeDetails.tsx`

-   **Edge Case**: The `id` from `useParams` does not correspond to any existing recipe, and `getRecipeById` returns an error.
    -   **Test**: Mock the `getRecipeById` call to throw an error. The test should assert that an error message like "Failed to load recipe details." is rendered.
-   **Edge Case**: A user tries to submit an empty comment or a comment with only whitespace.
    -   **Test**: Simulate filling the comment form with only spaces and submitting. The test should assert that the `postComment` API call is not made.

### `SignUp.tsx`

-   **Edge Case**: A user enters a username containing numbers or special characters.
    -   **Test**: Simulate typing a username with invalid characters into the input field. The test should assert that an error message is displayed and the form's state for `user_name` is not updated with the invalid characters (or that the form submission is prevented).
