# Poshta

Poshta is a secure, self-hosted messenger / chat backend written in Go.  
It provides HTTP and WebSocket APIs for user authentication, chat management, and real-time messaging.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [API Overview](#api-overview)
  - [Authentication](#authentication)
  - [Users](#users)
  - [Chats](#chats)
  - [Messages](#messages)
  - [WebSocket Endpoint](#websocket-endpoint)
  - [Healthcheck](#healthcheck)
  - [API Documentation (Swagger)](#api-documentation-swagger)
- [Database & Migrations](#database--migrations)
- [Development](#development)
- [Security, Privacy & Data Protection](#security-privacy--data-protection)
- [Digital Public Good Alignment](#digital-public-good-alignment)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [Authors](#authors)
- [License](#license)
- [Copyright](#copyright)

---

## Overview

Poshta is a backend service for chat / messaging applications.  
It is designed to be:

- **Open-source** and reusable by different frontends (web, mobile, desktop),
- **Secure** by default, with JWT-based authentication and configurable token lifetimes,
- **Deployable anywhere**, using environment variables for configuration and a standard SQL database.

This project can be used as the backend for a messenger for communities, NGOs, schools, or other organizations that need privacy-respecting, self-hosted communication tools.

---

## Features

Current / planned capabilities include:

- User registration and login with JWT-based authentication
- Access and refresh tokens with configurable TTLs
- Retrieval of a user’s public key for secure communication
- Creation and management of chats
- Sending and deleting messages
- Fetching chat messages and user chats
- Real-time messaging via WebSocket hub
- Swagger (OpenAPI) documentation endpoint
- Healthcheck endpoint for monitoring
- CORS configured for a frontend (e.g. `http://localhost:3000`)

> Note: Some features may be under active development. Check open issues and code for the latest status.

---

## Architecture

Poshta is written in Go and uses a layered architecture:

- **`cmd/app`** – Entry point of the application.
- **`internal/app`** – Application bootstrap (config, connections, startup, WebSocket hub).
- **`internal/handler`** – HTTP handlers for auth, chats, messages, WebSocket.
- **`internal/service` / `internal/usecase`** – Business logic and use cases.
- **`internal/repository`** – Database access using `sqlx`.
- **`pkg/logger`** – Centralized logging using `logrus`.
- **`migrations`** – Database schema migrations.
- **`docs`** – Auto-generated Swagger / OpenAPI documentation.

Configuration is loaded from a `.env` file and environment variables, using `caarlos0/env` and `godotenv`.

---

## Getting Started

### Prerequisites

- **Go** 1.23+ (as specified in `go.mod`)
- A running SQL database (e.g. PostgreSQL or MySQL)
- **Git** to clone the repository

Optional but helpful:

- [Air](https://github.com/cosmtrek/air) for live reload (supported via `.air.toml`)

### Configuration

The application is configured via environment variables. By default it loads from `.env` (path can be overridden via `-config` flag).

Create a `.env` file in the project root with values like:

```env
# HTTP server
SERVER_HOST=localhost
SERVER_PORT=8080

# Database (example for PostgreSQL)
DATABASE_DSN=postgres://user:password@localhost:5432/poshta?sslmode=disable

# JWT configuration
JWT_SECRET_KEY=change_me_to_a_strong_secret
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=72h
JWT_ISSUER=poshta-app
````

**Environment variables**

* `SERVER_HOST` – Host to bind the HTTP server (default: `localhost`)
* `SERVER_PORT` – Port for the HTTP server (default: `8080`)
* `DATABASE_DSN` – DSN for the SQL database (PostgreSQL/MySQL supported)
* `JWT_SECRET_KEY` – Secret key for signing JWTs
* `JWT_ACCESS_TOKEN_TTL` – Access token lifetime (e.g. `15m`)
* `JWT_REFRESH_TOKEN_TTL` – Refresh token lifetime (e.g. `72h`)
* `JWT_ISSUER` – Issuer field for JWT tokens (default: `poshta-app`)

### Running the Application

Clone the repository:

```bash
git clone https://github.com/zhanibek05/poshta.git
cd poshta
```

Run with Go:

```bash
go run ./cmd/app -config .env
```

The server will start at:

```text
http://localhost:8080
```

---

## API Overview

Base path for the API:

```text
http://{SERVER_HOST}:{SERVER_PORT}/api
```

Security:

* Bearer token (`Authorization: Bearer <access_token>`) using JWT.
* Some endpoints are public (e.g., registration, login); others are protected.

Below is a brief overview; check Swagger for full details.

### Authentication

* `POST /api/auth/register`
  Register a new user.

* `POST /api/auth/login`
  Log in and obtain access and refresh tokens.

* `POST /api/auth/refresh`
  Refresh an access token using a valid refresh token.

### Users

* `GET /api/{user_id}/public_key`
  Get the public key of a user (used for secure messaging).

* `GET /api/profile` (protected)
  Get the currently authenticated user’s profile.

### Chats

(Protected endpoints – require a valid JWT)

* `POST /api/chats`
  Create a new chat.

* `GET /api/chats/{user_id}/chats`
  Get the list of chats for a given user.

* `GET /api/chats/{chat_id}/messages`
  Get messages in a specific chat.

* `DELETE /api/chats/{chat_id}/chats`
  Delete a chat.

### Messages

(Protected endpoints)

* `POST /api/message`
  Send a new message.

* `DELETE /api/messages/{id}`
  Delete a message by ID.

### WebSocket Endpoint

* `GET /ws`
  WebSocket endpoint used for real-time communication. Clients connect with a valid token and then send/receive chat messages via the WebSocket protocol.

### Healthcheck

* `GET /healthcheck`

Returns a simple status message and HTTP `200 OK` if the server is healthy.

---

### API Documentation (Swagger)

Swagger (OpenAPI) docs are generated using `swaggo`.

After running the server, open:

```text
http://localhost:8080/swagger/index.html
```

There you can view and test all available endpoints.

---

## Database & Migrations

Database access is implemented via `sqlx`.
The schema and evolution scripts are stored under the `migrations/` directory.

Typical workflow:

1. Create the database in your SQL server.
2. Apply migration SQL files from `migrations/` in order.
3. Update `DATABASE_DSN` in your `.env` to point to the database.

> Note: The exact migration commands depend on your tooling (e.g. `goose`, `migrate`, or manual SQL). Adjust this section to match the tool you are using.

---

## Development

Run the full test suite:

```bash
go test ./...
```

Run with hot-reload using [Air](https://github.com/cosmtrek/air) (if installed):

```bash
air
```

This uses the configuration from `.air.toml` in the project root.

---

## Security, Privacy & Data Protection

* Authentication is implemented using **JWT** with configurable lifetimes for access and refresh tokens.
* Credentials and secrets are provided via **environment variables**, not hard-coded in the source code.
* The `DATABASE_DSN` should not contain production passwords committed to the repository.
* In production, the application should be run behind **HTTPS**, with TLS termination managed by a reverse proxy or load balancer.
* Logs should be configured to **avoid sensitive data** such as raw message content and passwords.
* Future improvements may include:

  * End-to-end encryption on the client side,
  * User data export and deletion endpoints,
  * Role-based access control (RBAC).

This aligns with a “privacy by design” approach, especially when used by communities, NGOs or schools.

---

## Digital Public Good Alignment

> This section is especially useful for university coursework and DPG self-assessment.

* **Relevant SDGs**: Poshta aims to support SDG 9 (*Industry, Innovation and Infrastructure*) and SDG 16 (*Peace, Justice and Strong Institutions*) by providing an open, low-cost communication backend that can be used by civic groups, schools, and community organizations.
* **Open Source**: The project is intended to be licensed under an OSI-approved license (see [License](#license)).
* **Platform Independent**: Poshta runs anywhere Go and a SQL database are available (bare metal, VM, Docker, any cloud provider).
* **Documentation**: This README and the `docs/` folder (Swagger) provide both user and contributor-oriented documentation.
* **Do No Harm / Non-Discrimination**: The project can be paired with a Code of Conduct and contribution guidelines (e.g. `CODE_OF_CONDUCT.md`, `CONTRIBUTING.md`) to foster an inclusive and safe community.

You can expand this section with links to your course’s DPG self-assessment document, Code of Conduct, and CONTRIBUTING file when they are added.

---

## Project Structure

High-level structure of the repository:

```text
cmd/
  app/               # Application entry point (main.go)

internal/
  app/               # Bootstrap: config, connections, startup, ws hub
  handler/           # HTTP & WebSocket handlers
  middleware/        # JWT and other middleware
  repository/        # Database repositories (users, chats, messages)
  service/           # Services (auth, etc.)
  usecase/           # Application use cases
  ...                # Other internal packages

migrations/          # Database migration scripts
docs/                # Swagger / OpenAPI generated docs
pkg/
  logger/            # Logging utilities

.air.toml            # Air (live reload) configuration
go.mod
go.sum
```

---

## Contributing

Contributions are welcome!

1. Fork the repository.
2. Create a new branch: `git checkout -b feature/my-feature`.
3. Make your changes and add tests if possible.
4. Run `go test ./...`.
5. Submit a pull request with a clear description of your changes.

When a `CONTRIBUTING.md` file is added, please follow the guidelines described there.

---

## Authors

* **Zhanibek Beisenov** – [GitHub](https://github.com/zhanibek05)
* **Daniyar Dautbaev** – [GitHub](https://github.com/daniyardautbaev)

---

## License

This project is intended to be licensed under the **MIT License**.
Once the `LICENSE` file is added, it will contain the full license text.

---

## Copyright

Copyright (c) 2025
**Zhanibek Beisenov**, **Daniyar Dautbaev**

See the [License](#license) section and `LICENSE` file for details.

```

If you want, I can also generate a `LICENSE` file (MIT) and a short `CONTRIBUTING.md` / `CODE_OF_CONDUCT.md` to match your DPG self-assessment.
::contentReference[oaicite:0]{index=0}
```
