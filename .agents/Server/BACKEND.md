# BiteBox Server

## Project Overview

This project is the back-end server for **BiteBox**, a recipe and meal-sharing application. It is a RESTful API built with Go. The server handles user authentication, recipe management, comments, and meal types. It uses a PostgreSQL database for data storage and Cloudinary for image hosting.

**Key Technologies:**

*   **Language:** Go
*   **Database:** PostgreSQL
*   **Authentication:** JWT (JSON Web Tokens)
*   **Image Hosting:** Cloudinary
*   **Web Framework:** Standard Go `net/http` library
*   **Dependencies:**
    *   `github.com/lib/pq`: PostgreSQL driver
    *   `github.com/golang-jwt/jwt/v5`: JWT implementation
    *   `github.com/joho/godotenv`: Environment variable loading
    *   `github.com/cloudinary/cloudinary-go/v2`: Cloudinary Go library

## Architecture

The application follows a standard Go web application structure:

*   **`main.go`**: The entry point of the application. It initializes the database connection, sets up the HTTP server, and registers all the API routes.
*   **`db/`**: Contains database-related files.
    *   `db.go`: Handles the database connection.
    *   `bitebox.sql`: The SQL schema for the PostgreSQL database.
*   **`handlers/`**: Contains the HTTP handlers for the different API endpoints, organized by resource (e.g., `auth`, `recipes`, `users`).
*   **`middlewares/`**: Contains HTTP middleware, such as the JWT authentication and CORS handlers.
*   **`lib/`**: Contains helper libraries, such as the Cloudinary client.
*   **`utils/`**: Contains utility functions, such as JWT generation.

## Building and Running

### Prerequisites

*   Go (version 1.24.2 or higher)
*   PostgreSQL
*   A `.env` file with the necessary environment variables (see below).

### Environment Variables

Create a `.env` file in the root of the project with the following variables:

```
SECRET_KEY=your_jwt_secret_key
CLOUDINARY_URL=your_cloudinary_url
```

### Running the Server

1.  **Set up the database:**
    *   Make sure your PostgreSQL server is running.
    *   Create a database named `bitebox`.
    *   Load the schema into the database:
        ```bash
        psql -d bitebox -f db/bitebox.sql
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

## Development Conventions

*   **Code Style:** Follow standard Go conventions (`gofmt`).
*   **Testing:** Unit tests are located in the same package as the code they are testing, using the `_test.go` suffix.
*   **Dependencies:** Dependencies are managed with Go Modules.
*   **API Design:** The API is designed to be RESTful, with endpoints organized around resources.
*   **Authentication:** JWTs are used to secure endpoints that require a logged-in user. The JWT is passed in the `Authorization` header as a Bearer token.
