# AionApi

_This repository is part of an ongoing study and development project. Many features are still being built and tested, so expect frequent changes as new tools and patterns are explored._

## Aion: Empowering you to take control of your time, habits, and aspirations

> _Aion is an innovative habit management system designed to help you organize, track, and analyze your daily routine for improved physical, mental, and emotional well-being. It combines cutting-edge technology with a user-centered approach to make your productivity and self-improvement journey seamless and insightful._
>
> Whether you’re focusing on fitness, learning, or personal growth, Aion is your companion in building the discipline you need to achieve sustainable success.

## **Table of Contents**

- [Overview](#overview)
- [Current and Upcoming Features](#current-and-upcoming-features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development](#development)
- [Docker Integration](#docker-integration)
- [API Endpoints](#api-endpoints)
    - [REST Endpoints](#rest-endpoints)
    - [GraphQL Operations](#graphql-operations)
- [License](#license)

---

## Overview

AionAPI is a backend service written in **Go** that exposes both REST and GraphQL endpoints. It relies on **PostgreSQL** for persistent storage and **Redis** for caching. The project follows the **Ports & Adapters (Hexagonal)** architecture, enabling clear separation between domain logic and external technologies. Observability is handled with **OpenTelemetry**, **Jaeger**, **Prometheus**, and **Grafana**, while structured logging is powered by **Zap**.

**Technology Stack**

- Go 1.24
- chi HTTP router and gqlgen for GraphQL
- GORM ORM for PostgreSQL
- Redis cache
- Docker & Docker Compose
- OpenTelemetry with Jaeger tracing and Prometheus metrics
- Grafana dashboards
- zap for structured logging

---

## Current and Upcoming Features

- **Streamlined Habit Management** — organize and track your habits effortlessly.
- **Data-Driven Insights** — visualize your progress and analyze behavior patterns.
- **Modern Integrations** — sync with tools and platforms for extended usability.
- **Developer-Friendly API** — clean, extensible endpoints for all your needs.

---

## Installation

### Prerequisites

- [Go](https://go.dev/doc/install) 1.24 or newer
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/) (for containerized development)

### Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/lechitz/AionApi.git
   cd AionApi
   ```
2. **Install development tools** (optional)
   ```bash
   make tools-install
   ```
3. **Download Go dependencies**
   ```bash
   go mod tidy
   ```

---

## Configuration

1. **Copy the example environment file**
   ```bash
   cp infrastructure/docker/example/.env.example infrastructure/docker/dev/.env.dev
   ```
2. **Edit `.env.dev`** with values that match your local setup.
3. **Run database migrations** (optional)
   ```bash
   make migrate-up
   ```
4. **Start the development environment**
   ```bash
   make dev-up
   ```

---

## Development

The project is organized as follows:

```text
cmd/            - application entry point
infrastructure/ - migrations, docker files, observability
internal/       - domain logic and adapters
makefiles/      - grouped Make targets
pkg/            - shared utilities (zap logger, helpers)
tests/          - coverage, mocks, setup, testdata
```

Run `make help` to see all available commands. Frequently used ones include:

```bash
make format    # format Go code
make lint      # run static analysis
make test      # execute unit tests
make verify    # run full pipeline before committing
```

---

## Docker Integration

Docker compose files for development and production live under `infrastructure/docker`. The main targets are:

```bash
make build-dev   # build development image
make dev-up      # start development environment (includes Postgres, Redis, Jaeger, Prometheus, Grafana)
make dev-down    # stop and remove development containers
make build-prod  # build production image
make prod-up     # start production environment
make prod-down   # stop production containers
```

---

## API Endpoints

The API exposes REST endpoints for user management, authentication, and health checks, along with GraphQL operations for categories and tags.

### REST Endpoints
- `GET  /aion-api/health-check/` — service status
- `POST /aion-api/user/create` — create a new user
- `GET  /aion-api/user/all` — list users
- `GET  /aion-api/user/{user_id}` — retrieve a user by ID
- `PUT  /aion-api/user/` — update user data
- `PUT  /aion-api/user/password` — update the logged user's password
- `DELETE /aion-api/user/` — soft delete the logged user
- `POST /aion-api/auth/login` — obtain a JWT token
- `POST /aion-api/auth/logout` — invalidate the user session

### GraphQL Operations
Endpoint: `/graphql`
- Queries: `GetAllCategories`, `GetCategoryByID`, `GetCategoryByName`, `GetAllTags`, `GetTagByID`
- Mutations: `CreateCategory`, `CreateTag`, `UpdateCategory`, `SoftDeleteCategory`

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
