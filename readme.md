# 2FA Code Generator API

A simple, robust, and secure REST API for generating and verifying Time-based One-Time Passwords (TOTP), commonly used for Two-Factor Authentication (2FA). This project is built in Go and follows a clean, layered architecture based on SOLID principles.

## ✨ Features

- **Time-Windowed Codes**: Generates 6-digit codes that are valid for fixed 30-second windows based on the UTC clock (e.g., `xx:xx:00` to `xx:xx:29`).
- **Deterministic Generation**: Codes are generated deterministically based on the User ID and the current time window.
- **Interactive API Documentation**: Comes with a built-in Swagger UI to easily test and explore the API.

## 🚀 Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites

- **Go**: Version 1.22 or later.
- **Swag CLI**: The command-line tool for generating Swagger documentation.

You can install the Swag CLI with the following command:
```bash
go install [github.com/swaggo/swag/cmd/swag@latest](https://github.com/swaggo/swag/cmd/swag@latest)
```
*Note: Ensure your Go binary path (`$GOPATH/bin`) is in your system's `PATH`.*

### Installation & Running

1.  **Build the Docker Image**:
    From your project's root directory, run the `docker build` command. This will compile your Go application and package it into an image named `2fa-api`.
    ```bash
    docker build -t 2fa-api .
    ```

2.  **Run the Docker Container**:
    Once the image is built, run it using the `docker run` command. This command also creates a persistent volume for the database.
    ```bash
    docker run -p 8080:8080 -v "$(pwd)/data:/app" 2fa-api
    ```
    - **`-p 8080:8080`**: Maps port 8080 on your local machine to port 8080 inside the container.
    - **`-v "$(pwd)/data:/app"`**: Creates a volume that maps a `data` folder in your project directory to the `/app` directory inside the container. This is where the `2fa.db` file will be stored, ensuring it persists.

The server is now running, and you can access the API and Swagger UI.

---

## 📖 API Documentation

This project includes interactive API documentation via Swagger UI. Once the container is running, you can access it at:

**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**


### Endpoints

#### `GET /code`

Generates a new 6-digit code for a user.

-   **Method**: `GET`
-   **Query Parameters**:
    -   `userID` (string, required)
    -   `userEmail` (string, required)
-   **Success Response (200 OK)**:
    ```json
    {
      "code": "123456",
      "generated_at": "2025-07-22T18:01:00Z",
      "expires_at": "2025-07-22T18:01:30Z"
    }
    ```

#### `POST /verify`

Verifies a code submitted by a user.

-   **Method**: `POST`
-   **Request Body**:
    ```json
    {
      "user_id": "user-123",
      "user_email": "test@example.com",
      "code": "123456"
    }
    ```
-   **Success Response (200 OK)**:
    ```json
    {
      "valid": true,
      "message": "Code verified successfully."
    }
    ```

## 🔮 Future Enhancements

Here are some ideas for future improvements to make the project even more robust and feature-rich:

-   **JWT Bearer Token Authentication**:
    -   Implement a new `/login` endpoint that, upon successful credential validation, issues a signed JWT (JSON Web Token).
    -   The JWT would include the `userID` and `userEmail` as custom claims.
    -   Protect the `/code` and `/verify` endpoints with middleware that validates this Bearer Token, removing the need to pass user details in every request.

-   **Rate Limiting**:
    -   Add middleware to limit the number of requests to the `/code` and `/verify` endpoints per user or IP address to prevent brute-force attacks and abuse.

-   **Configuration Management**:
    -   Move hardcoded values (like the database path and JWT secret key) into a configuration file (e.g., `config.yaml`) or environment variables for better security and flexibility.

-   **CI/CD Pipeline**:
    -   Set up a continuous integration/continuous deployment pipeline (e.g., using GitHub Actions) to automatically run tests, build the Docker image, and push it to a container registry on every commit.
