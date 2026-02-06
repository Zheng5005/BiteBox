# AGENTS.md

This file provides a comprehensive overview of the BiteBox project for agents interactions.

## Project Overview

BiteBox is a full-stack web application. The project is structured as a monorepo with a React frontend and a Go backend.

*   **Frontend:** The `Client` directory contains a React application built with Vite. It uses TypeScript, React Router for navigation, and Tailwind CSS for styling.
*   **Backend:** The `Server` directory contains a Go application that serves a RESTful API. It uses the standard `net/http` library for routing, `pq` for connecting to a PostgreSQL database, and JWTs for authentication.
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
*   `npm run lint`: Lints the codebase.
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
The server uses a PostgreSQL database. A `bitebox.sql` file is present in the `db/` directory, which likely contains the database schema. You will need a running PostgreSQL instance and a `.env` file with the database connection details.

---
## Development Conventions
*   **Monorepo:** The project is a monorepo with client and server code in separate directories.

### Frontend Development Conventions
*   **Styling:** The frontend uses Tailwind CSS for styling.

### Backend Development Conventions
*   **API:** The Go server provides a RESTful API consumed by the React client. API routes are prefixed with `/api`.
*   **Authentication:** Authentication is handled via JWTs. The `JWTMiddleware` in the Go backend protects certain routes.
