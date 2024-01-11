# Chirpy

Chirpy is a Go project that provides a web API for managing chirps. It utilizes the Chi router and a PostgreSQL database.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [License](#license)

## Installation

1. Clone the repository: `git clone https://github.com/andrejscek/Chirpy.git`
2. Navigate to the project directory: `cd Chirpy`
3. Install the dependencies: `go mod download`
4. Build the project: `go build`
5. Create a `.env` file in the project root directory and set the following environment variables:
    - `JWT_SECRET`: Secret key for JWT token generation, can be whatever
    - `POLKA_API_KEY`: API key for fake Polka integration, can be whatever
6. Start the server by running `./Chirpy`

## Usage

Once the server is running, you can access the API endpoints using a tool like cURL or Postman. The base URL for the API is `http://localhost:8080/api`.

## API Endpoints

- `GET /api/healthz`: Check the health status of the server.
- `GET /api/chirps`: Get a list of all chirps. Supports optinal query parameters 'author_id' and 'sort'.
- `GET /api/chirps/{id}`: Get a specific chirp by ID.
- `POST /api/chirps`: Create a new chirp. You need to be provide a valid JWT token.
- `DELETE /api/chirps/{id}`: Delete a specific chirp by ID. You need to provide a valid JWT token.
- `POST /api/users`: Create a new user. ou need to provide a username and password.
- `POST /api/login`: Log in as a user. You need to provide a valid username and password. Returns a JWT token and user ID.
- `PUT /api/users`: Update user information. You need to provide a valid JWT token.
- `POST /api/refresh`: Refresh the access token. You need to provide a valid JWT refresh token.
- `POST /api/revoke`: Revoke the refresh token. You need to provide a valid JWT refresh token.
- `POST /api/polka/webhooks`: Handle Polka webhooks. You need to provide a valid Polka API key. Used to change the user status to is_chirpy_red true.
- `GET /api/reset`: Reset the visited count for app endpoints.
- `GET /admin/metrics`: Get the visited count for app endpoints.

## License

This project is licensed under the MIT License.
