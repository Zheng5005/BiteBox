# BiteBox Client

This document provides a comprehensive overview of the BiteBox client application, intended to be used as a context for an AI assistant.

## Project Overview

BiteBox is a web application for browsing and sharing recipes. It is built with React, TypeScript, and Vite, and uses Tailwind CSS for styling. The application allows users to view, search, and filter recipes. Authenticated users can also post new recipes.

### Technologies

*   **Framework:** React
*   **Language:** TypeScript
*   **Build Tool:** Vite
*   **Styling:** Tailwind CSS
*   **Routing:** React Router
*   **Linting:** ESLint

### Architecture

The application is structured into the following directories:

*   **`src/components`:** Contains reusable UI components like `Navbar` and `Avatar`.
*   **`src/pages`:** Contains the main pages of the application, such as `MainPage`, `Login`, and `RecipeDetails`.
*   **`src/context`:** Contains the `AuthContext` for managing user authentication.
*   **`src/hooks`:** Contains custom hooks, such as `useMealTypes`.
*   **`src/`:** The root of the source code, containing the main entry point (`main.tsx`) and the main application component (`App.tsx`).

## Building and Running

The following scripts are available in `package.json`:

*   **`npm run dev`**: Starts the development server.
*   **`npm run build`**: Builds the application for production.
*   **`npm run lint`**: Lints the codebase using ESLint.
*   **`npm run preview`**: Starts a local server to preview the production build.

### Development

To run the application in development mode, use the following command:

```bash
npm run dev
```

### Production

To build the application for production, use the following command:

```bash
npm run build
```

## Development Conventions

### Coding Style

The project uses ESLint to enforce a consistent coding style. The ESLint configuration is located in `eslint.config.js`.

### Testing

- Use Vitest as a test Framework
- Use React-Testing-Framework
