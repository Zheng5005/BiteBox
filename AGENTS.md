# AGENTS.md

This file provides a comprehensive overview of the BiteBox project for agents interactions.

## Project Overview

BiteBox is a full-stack web application. The project is structured as a monorepo with a React frontend and a Go backend.

*   **Frontend:** The `Client` directory contains a React application built with Vite. It uses TypeScript, React Router for navigation, and Tailwind CSS for styling. The client application can be containerized using `Client/Dockerfile`.
*   **Backend:** The `Server` directory contains a Go application that serves a RESTful API. It uses the standard `net/http` library for routing, `pq` for connecting to a PostgreSQL database, and JWTs for authentication. The server application can be containerized using `Server/Dockerfile`.
*   **CI/CD:** The project includes GitHub Actions workflows defined in `.github/workflows/main.yml` for continuous integration and deployment.

---
## Client (Frontend)

The client-side application is located in the `Client/` directory.

### Running the Client

1.  **Navigate to the client directory:**
    ```bash
    cd Client
    ```
2.  **Install dependencies:**
    ```bash
    npm install
    ```
3.  **Run the development server:**
    ```bash
    npm run dev
    ```
    The application will be available at `http://localhost:5173` (or another port if 5173 is in use, check the Vite output).

### Key Scripts
*   `npm run dev`: Starts the development server.
*   `npm run build`: Builds the application for production.
*   `npm run lint`: Lints the codebase using `eslint.config.js`.

### Client Structure
The `Client/src` directory contains the core application logic and is organized as follows:
*   `api/`: Contains API service modules for interacting with the backend.
*   `components/`: Reusable React components.
*   `context/`: React Context providers for global state management.
*   `hooks/`: Custom React hooks.
*   `pages/`: Top-level page components for routing.
*   `test/`: Test setup and utilities.

---
## Server (Backend)

The server-side application is located in the `Server/` directory.

### Running the Server

1.  **Navigate to the server directory:**
    ```bash
    cd Server
    ```
2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Run the server:**
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

### Database
The server uses a PostgreSQL database. A `bitebox.sql` file is present in the `db/` directory, which likely contains the database schema. You will need a running PostgreSQL instance and a `.env` file with the database connection details. A `db/Dockerfile` is also available for containerizing the database.

### Server Structure
The `Server/` directory is organized as follows:
*   `db/`: Database schema and related utilities.
*   `handlers/`: Contains API route handlers, further organized by domain:
    *   `auth/`: Authentication related handlers.
    *   `comments/`: Comment related handlers.
    *   `meals/`: Meal related handlers.
    *   `recipes/`: Recipe related handlers.
    *   `users/`: User related handlers.
*   `lib/`: External library integrations, e.g., `cloudinary.go` for Cloudinary.
*   `middlewares/`: HTTP middleware, e.g., `cors.go` for CORS and `jwt.go` for JWT authentication.
*   `utils/`: Utility functions, e.g., `jwt.go` for JWT token handling.

---
## Development Conventions
*   **Monorepo:** The project is a monorepo with client and server code in separate directories.

### Frontend Development Conventions
*   **Styling:** The frontend uses Tailwind CSS for styling.
*   **Linting:** ESLint configuration is located at `Client/eslint.config.js`.

### Backend Development Conventions
*   **API:** The Go server provides a RESTful API consumed by the React client. API routes are prefixed with `/api`.
*   **Authentication:** Authentication is handled via JWTs. The `JWTMiddleware` in the Go backend protects certain routes.

---
## Agent-Specific Documentation

The `.agents/` directory contains additional markdown files specifically for agent interactions, providing more detailed guidance and context for various aspects of the project:
*   `AGENTS_WORKFLOW.md`: General workflow guidelines for agents.
*   `API_CONTRACTS.md`: Details about the API contracts between client and server.
*   `Client/`: Client-specific agent documentation:
    *   `edge_case_tests.md`: Information about edge case testing for the client.
    *   `FRONTEND.md`: Detailed frontend development guidelines.
    *   `MVP_todo.md`: To-do list for Minimum Viable Product features on the client.
    *   `test_plan.md`: Client testing plan.
*   `Server/`: Server-specific agent documentation:
    *   `BACKEND.md`: Detailed backend development guidelines.
    *   `MICROCHECKLIST.md`: Micro-checklist for server-side tasks.
    *   `MVP_todo.md`: To-do list for Minimum Viable Product features on the server.