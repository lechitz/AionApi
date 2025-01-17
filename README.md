# AionApi

## Aion: Empowering you to take control of your time, habits, and aspirations

> _Aion is an innovative habit management system designed to help you organize, track, and analyze your daily routine for improved physical, mental, and emotional well-being. It combines cutting-edge technology with a user-centered approach to make your productivity and self-improvement journey seamless and insightful._
>
> Whether you’re focusing on fitness, learning, or personal growth, Aion is your companion in building the discipline you need to achieve sustainable success.

## **Table of Contents**

- [Overview](#overview)
- [Current and Upcoming Features](#Current-and-Upcoming-Features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development](#development)
- [Docker Integration](#docker-integration)
- [API Endpoints](#api-endpoints)
    - [User Management](#user-management)
    - [Authentication](#authentication)
- [License](#license)

---

## **Overview**

The Aion API is a RESTful backend solution designed to empower developers to build productivity applications with robust user and activity management features. Built with **Go**, powered by **PostgreSQL**, and based on the **Ports & Adapters architecture** for scalability and maintainability.

---

## **Current and Upcoming Features**

- **Streamlined Habit Management:** Organize and track your habits effortlessly.
- **Data-Driven Insights:** Visualize your progress and analyze behavior patterns.
- **Modern Integrations:** Sync with tools and platforms for extended usability.
- **Developer-Friendly API:** Clear, scalable, and extensible endpoints for all your needs.

---

## **Installation**

1. **Clone the repository:**
   ```bash
   git clone https://github.com/lechitz/AionApi.git
   cd AionApi
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

---

## **Configuration**

1. **Set up environment variables:**
   ```bash
   cp .env.example .env
   ```

2. **Update the `.env` file** with your database credentials and configurations:

<div style="text-align: center;">

| Variable         | Description                     |
|------------------|---------------------------------|
| `SERVER_CONTEXT` | Path for the server context     |
| `PORT`           | The port the server will run on |
| `DB_USER`        | Database username               |
| `DB_PASSWORD`    | Database password               |
| `DB_NAME`        | Database name                   |
| `DB_HOST`        | Database host address           |
| `DB_PORT`        | Database port                   |
| `DB_TYPE`        | Database type                   |
| `SECRET_KEY`     | Secret key for JWT              |

</div>

3. **Set up the database:**
    - Create a new database in PostgreSQL.
    - Connect to the database by configuring the connection settings in the `.env` file.
    - Run all migration files located in `infra/db/migrations`.

4. **Generate JWT Secret Key (optional):**
    - Use the `GenerateJWTKey()` function to produce a random 64-byte secret key (base64-encoded).
    - Run the key generator:
      ```bash
      go run pkg/utils/generate_jwt_key.go
      ```
    - Copy the generated key into your `.env` file under `SECRET_KEY`.

---

## **Development**

### **Folder Structure**

The Aion API adopts the Ports and Adapters Architecture (Hexagonal Architecture) to ensure flexibility, scalability, and testability. Below is an overview of the folder structure:

<div style="text-align: center;">

| Directory       | Description                                                 |
|-----------------|-------------------------------------------------------------|
| `adapters`      | Contains input/output adapters for external communications. |
| `cmd`           | Main application entry point.                               |
| `config`        | Environment variables and configurations.                   |
| `infra`         | Infrastructure setup including database connections.        |
| `internal`      | Core business logic and domain entities.                    |
| `pkg`           | Shared utilities and helper functions.                      |
| `ports`         | Input/output interfaces for the application.                |

</div>



The Aion API is organized following the **Ports and Adapters Architecture (Hexagonal Architecture)** to promote clear separation between the core business logic and external systems. Here is the folder structure:

<details>
<summary>
Complete Folder Structure
</summary>

```plaintext
.
├── adapters
│   ├── input
│   │   └── http
│   │       ├── dto
│   │       ├── handlers
│   │       └── server

│   └── output
│       ├── cache
│       │   └── redis
│       └── db
│           └── postgres
├── app
│   ├── bootstrap
│   ├── config
│   ├── logger
│   ├── logs
│   └── middlewares
│       └── auth
├── cmd
│   └── aion-api
│       └── main.go
├── core
│   ├── domain
│   │   ├── entities
│   │   ├── events
│   │   └── exceptions
│   ├── msg
│   └── service
├── infra
│   ├── cache
│   ├── db
│   │   ├── migrations
│   │   ├── postgres
│   │   │   └── migrations
│   ├── messaging
│   ├── observability
├── pkg
│   ├── contextkeys
│   ├── errors
│   └── utils
├── ports
│   ├── input
│   │   └── http
│   └── output
│       ├── cache
│       └── db
├── .env
├── .env.example
├── .gitignore
├── docker-compose-dev.yaml
├── docker-compose-prod.yaml
├── Dockerfile
├── go.mod
├── LICENSE
├── Makefile
└──  README.md
```

</details>

---

## **Docker Integration**

### **Prerequisites**

- Install **Docker**. Follow the [official Docker documentation](https://docs.docker.com/get-docker/) for installation instructions.
- Optionally, install **Docker Compose** for managing multi-container setups.

### **Running the Project with Docker**

---

<details>
<summary> 
 The development environment setup
</summary>

1. **Build the development image:**
   ```bash
   make docker-build-dev
   ```

2. **Start the development environment:**
   ```bash
   make docker-compose-dev-up
   ```

3. **Stop the development environment:**
   ```bash
   make docker-compose-dev-down
   ```
</details>

----

<details>
<summary> 
 The production environment setup
</summary>

1. **Build the production image:**
   ```bash
   make docker-build-prod
   ```

2. **Start the production environment:**
   ```bash
   make docker-compose-prod-up
   ```

3. **Stop the production environment:**
   ```bash
   make docker-compose-prod-down
   ```

</details>

---

## **API Endpoints**

### **User Management**

<details>
<summary> 
 The user management endpoints here
</summary>

#### **Create User**

- **Method:** `POST`
- **Endpoint:** `localhost:5001/aion-api/user/create`
- **Request Body:**

  ```json
  {
    "name": "John Doe",
    "username": "johndoe",
    "email": "johndoe@example.com",
    "password": "securePassword123"
  }
  ```
- **Response:**

  ```json
  {
    "message": "user created successfully",
    "result": {
        "id": 1,
        "name": "John Doe",
        "username": "johndoe",
        "email": "johndoe@example.com"
    },
    "date": "2025-01-07T15:41:50.803251738Z"
  }
  ```

#### **Get All Users**

- **Method:** `GET`
- **Endpoint:** `localhost:5001/aion-api/user/all`
- **Headers:**
    - `Authorization: Bearer <token>`
- **Response:**

  ```json
  {
    "message": "users get successfully",
    "result": [
        {
            "id": 1,
            "name": "John Doe",
            "username": "johndoe",
            "email": "johndoe@example.com",
            "created_at": "2025-01-07T15:41:50.800147Z"
        },
        {
            "id": 2,
            "name": "Alice Smith",
            "username": "alicesmith",
            "email": "alice.smith@example.com",
            "created_at": "2025-01-07T15:56:32.174753Z"
        }
    ],
    "date": "2025-01-07T15:56:35.028477172Z"
  }
  ```

#### **Get User by ID**

- **Method:** `GET`
- **Endpoint:** `localhost:5001/aion-api/user/{id}`
- **Headers:**
    - `Authorization: Bearer <token>`
- **Response:**

  ```json
  {
    "message": "user get successfully",
    "result": [
        {
            "id": 1,
            "name": "John Doe",
            "username": "johndoe",
            "email": "johndoe@example.com",
            "created_at": "2025-01-07T15:41:50.800147Z"
        }
    ],
    "date": "2025-01-07T15:59:05.406681717Z"
  }
  ```

#### **Update User**

- **Method:** `PUT`
- **Endpoint:** `localhost:5001/aion-api/user/{id}`
- **Request Body:**

  ```json
  {
      "name": "Mark Taylor",
      "username": "markt89",
      "email": "mark.taylor@example.com"
  }
  ```
- **Response:**

  ```json
  {
    "message": "user updated successfully",
    "result": {
        "id": 2,
        "name": "Mark Taylor",
        "username": "markt89",
        "email": "mark.taylor@example.com",
        "updated_at": "2025-01-07T16:01:47.919084Z"
    },
    "date": "2025-01-07T16:01:47.929188372Z"
  }
  ```

#### **Soft Delete User**

- **Method:** `DELETE`
- **Endpoint:** `localhost:5001/aion-api/user/{id}`
- **Headers:**
    - `Authorization: Bearer <token>`

</details>

### **Authentication**

<details>
<summary> 
 The authentication endpoints here
</summary>

#### **Login**

- **Method:** `POST`
- **Endpoint:** `localhost:5001/aion-api/login`
- **Request Body:**

  ```json
  {
      "username": "johndoe",
      "password": "securePassword123"
  }
  ```
- **Response:**

  ```json
  {
    "message": "success to login",
    "result": {
        "username": "johndoe",
        "token": "eyJhbUHaiHL9AS6IkpXVCJ9.eyJhdXRob3JKAPkSVnS"
    },
    "date": "2025-01-07T15:50:48.751092612Z"
  }
  ```
</details>


---

## **License**

- This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
