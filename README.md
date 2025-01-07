# AionApi

<h2>Aion: Empowering you to take control of your time, habits, and aspirations.</h2>

> _Aion is an innovative habit management system designed to help you organize, track, and analyze your daily routine for improved physical, mental, and emotional well-being. It combines cutting-edge technology with a user-centered approach to make your productivity and self-improvement journey seamless and insightful._
>
> Whether you’re focusing on fitness, learning, or personal growth, Aion is your companion in building the discipline you need to achieve sustainable success.


**With Aion, you can:**
  
  - Streamline your habits with an intuitive system that adapts to your needs.
  - Track and analyze your activities to uncover insights about your behavior.
  - Visualize progress in real-time to stay motivated and on track with your goals.
  - Integrate effortlessly with modern tools and platforms for enhanced usability.
    

## **Table of Contents**

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development](#development)
- [API Endpoints](#api-endpoints)
    - [User Management](#user-management)
    - [Authentication](#authentication)
    - [Activity Management](#activity-management)
- [Docker Integration](#docker-integration)
- [License](#license)

---

## **Overview**

The Aion API is a RESTful backend solution designed to help developers build productivity applications with robust user and activity management capabilities. Built with **Go**, powered by **PostgreSQL**, and based on the **Ports & Adapters architecture** for scalability.

Aion API is your go-to platform for managing users, activities, and personal workflows.

## **Features**

- **User Management:**
    - User creation, retrieval, update, and soft deletion.
    - Authentication with secure JWT-based mechanisms.
- **Activity Tracking:**
    - Organize and track daily activities using custom tags.
    - Flexible structure to support personal or organizational workflows.
- **Developer-Friendly:**
    - Comprehensive documentation.
    - Easy-to-use API endpoints.
    - Architecture designed for extensibility, scalability, and maintainability.
- **Future-Ready:**
    - Prepared for frontend integration and advanced features like analytics.

---

## **Installation**

To set up and run the Aion API, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/lechitz/AionApi.git
   cd AionApi
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## **Configuration**

1. Set up the environment variables:

   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your database credentials and other configurations.

<div align="center">

| Variable            | Description                        |
|---------------------|------------------------------------|
| `SERVER_CONTEXT`     | Path for the server context        |
| `PORT`               | The port the server will run on    |
| `DB_USER`            | Database username                  |
| `DB_PASSWORD`        | Database password                  |
| `DB_NAME`            | Database name                      |
| `DB_HOST`            | Database host address              |
| `DB_PORT`            | Database port                      |
| `DB_TYPE`            | Database type (e.g., postgres)     |
| `SECRET_KEY`         | Secret key for JWT                 |
</div>

3. Set up the database:

    - Create a new database in PostgreSQL.
    - Connect to the database by configuring the connection settings in the .env file.
    - Run all the migration files in the `infra/db/migrations` folder.


4. Generate JWT Secret Key:

     - The `GenerateJWTKey()` function generates a unique secret key for JWT authentication and saves it to the `.env` file. This function is commented out by default in the `init()` function.

     - **Steps to Set Up the JWT Key:**
          1. **Uncomment the `GenerateJWTKey()` call in the `init()` function** during the first setup.  
          2. Run the application once to generate the key.  
          3. After the key is generated and stored in the `SECRET_KEY` variable in the `.env` file, **comment out the `GenerateJWTKey()` function again** to prevent overwriting the key in future runs.

     - **Example Code:**

		```go
		func init() {
		    // Uncomment the line below for the first setup
		    middlewares.GenerateJWTKey()
		}
		```


## **Development**

### **Folder Structure**

  - The Ports and Adapters Architecture (Hexagonal Architecture) was chosen for the Aion API to ensure flexibility, scalability, and testability. This design approach promotes a clear separation of concerns, dividing the application into core business logic (internal) and external systems (adapters).

<div align="center">
	
| **Directory** | **Description**                                                        |
|---------------|------------------------------------------------------------------------|
| `adapters`    | Contains the input and output adapters for the application.            |
| `cmd`         | Contains the main application entry point.                             |
| `config`      | Contains environment variables and configurations.                     |
| `infra`       | Contains the database connection and migration files.                  |
| `internal`    | Contains the core business logic of the application.                   |
| `pkg`         | Contains utility functions and shared code.                            |
| `ports`       | Contains the input and output interfaces for the application.          |
</div>

<details>
<summary> 
The complete folder structure here
</summary>

```plaintext
.
├── adapters
│   ├── input
│   │   └── http
│   │       ├── handlers
│   │       │   ├── generic.go
│   │       │   ├── models.go
│   │       │   └── user.go
│   │       └── router.go
│   ├── middlewares
│   │   ├── auth.go
│   │   ├── jwt.go
│   │   └── security.go
│   ├── output
│   │   └── repository
│   │       └── user.go
│   └── security
│       └── token.go
├── cmd
│   └── aion-api
│       └── main.go
├── config
│   └── environments.go
├── docker-compose-dev.yaml
├── docker-compose-prod.yaml
├── Dockerfile
├── docs
├── go.mod
├── go.sum
├── infra
│   └── db
│       ├── migrations
│       │   ├── 01-users.sql
│       │   ├── 02-tags.sql
│       │   ├── 03_days.sql
│       │   ├── 04_personal_diaries.sql
│       │   ├── 05_professional_diaries.sql
│       │   ├── 06_day_tag_summary.sql
│       │   ├── 07_day_moods.sql
│       │   ├── 08_day_energy.sql
│       │   ├── 09_day_water_intake.sql
│       │   └── 10_day_intentions.sql
│       └── postgres.go
├── internal
│   └── core
│       ├── domain
│       │   ├── context.go
│       │   └── user.go
│       └── service
│           └── user.go
├── LICENSE
├── Makefile
├── pkg
│   └── utils
│       ├── error_handler.go
│       └── response.go
├── ports
│   ├── input
│   │   └── user.go
│   └── output
│       └── user.go
└── README.md
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
- **Endpoint:** `/user/create`
- **Request Body:**

  ```json
  {
      "name": "string",
      "username": "string",
      "email": "string",
      "password": "string"
  }
  ```
- **Response:**

  ```json
  {
      "id": "uuid",
      "name": "string",
      "username": "string",
      "email": "string"
  }
  ```

#### **Get All Users**
- **Method:** `GET`
- **Endpoint:** `/user/all`
- **Headers:**
    - `Authorization: Bearer <token>`

#### **Get User by ID**
- **Method:** `GET`
- **Endpoint:** `/user/{id}`
- **Headers:**
    - `Authorization: Bearer <token>`

#### **Update User**
- **Method:** `PUT`
- **Endpoint:** `/user/{id}`
- **Request Body:**

  ```json
  {
      "name": "string",
      "username": "string",
      "email": "string"
  }
  ```

#### **Soft Delete User**
- **Method:** `DELETE`
- **Endpoint:** `/user/{id}`
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
- **Endpoint:** `/login`
- **Request Body:**

  ```json
  {
      "username": "string",
      "password": "string"
  }
  ```
- **Response:**

  ```json
  {
      "token": "string"
  }
  ```
</details>

### **Activity Management**

<details>
<summary> 
 The activity management endpoints here
</summary>

#### **Create Activity**
- **Method:** `POST`
- **Endpoint:** `/activity/create`
- **Request Body:**

  ```json
  {
      "title": "string",
      "description": "string",
      "tags": ["string"]
  }
  ```

#### **Get All Activities**
- **Method:** `GET`
- **Endpoint:** `/activity/all`
- **Headers:**
    - `Authorization: Bearer <token>`

#### **Update Activity**
- **Method:** `PUT`
- **Endpoint:** `/activity/{id}`
- **Request Body:**

  ```json
  {
      "title": "string",
      "description": "string",
      "tags": ["string"]
  }
  ```

#### **Delete Activity**
- **Method:** `DELETE`
- **Endpoint:** `/activity/{id}`
- **Headers:**
    - `Authorization: Bearer <token>`

---
</details>

---

## **Docker Integration**

### Prerequisites

1. Download and install Docker for your system:

    - [Docker Desktop for Windows](https://docs.docker.com/docker-for-windows/install/)
    - [Docker Desktop for Mac](https://docs.docker.com/docker-for-mac/install/)
    - [Docker Engine for Linux](https://docs.docker.com/engine/install/)

2. Verify the installation:

   ```bash
   docker --version
   ```

3. Choose the appropriate Docker Compose file:

    - `docker-compose-dev.yaml`: for development.
    - `docker-compose-prod.yaml`: for production.

4. Build the Docker image:
   ```bash
   make docker-build-prod
   ```

5. Start the Docker container:

   ```bash
   make docker-compose-prod-up
   ```

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
