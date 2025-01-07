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
    - [Prerequisites](#prerequisites)
    - [Running the Project with Docker](#running-the-project-with-docker)
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

---

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
---
     
## **Docker Integration**

### **Prerequisites**

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

### **Running the Project with Docker**

#### For Development:
   - To start the application in a development environment, use the following commands:
     - First, build the Docker image for development:
       
        ```bash
         make docker-build-dev
        ```
        
     - Then, start the application with the development Docker Compose file:
       
        ```bash
        make docker-compose-dev-up
        ```

    
#### For Production:
   - To start the application in a production environment, follow these steps:
     - First, build the Docker image for development:
       
        ```bash
         make docker-build-prod
        ```
        
     - Then, start the application with the development Docker Compose file:
       
        ```bash
        make docker-compose-prod-up
        ```
        

### Other Available Commands in Makefile:
   - In addition to the commands for starting the development and production environments, the Makefile also includes other useful commands, such as those for cleaning up Docker containers, volumes, and images, as well as specific commands for running tests and generating coverage reports. To view all available commands, simply run `make help`.
     

---

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

### **Activity Management**

<details>
<summary> 
 The activity management endpoints here
</summary>

- This section is still under development.


</details>

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
