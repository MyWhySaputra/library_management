# Library Management System

## Overview
This project is a backend system for a library management application. The system is built using Golang with microservice architecture and gRPC for communication between services. Each microservice has its own PostgreSQL database. Authentication and authorization are managed using JWT, and the services are containerized using Docker.

## Microservices
1. **BookService**: Manages books.
2. **AuthorService**: Manages authors.
3. **CategoryService**: Manages book categories.
4. **UserService**: Manages users and authentication.

## Architecture
The system is designed to be scalable and resilient using microservices architecture. Each service communicates with others via gRPC and each has its own separate database to ensure data consistency and isolation.

## Features
1. **Book Management**: Create, read, update, delete books.
2. **Author Management**: Create, read, update, delete authors.
3. **Category Management**: Create, read, update, delete categories.
4. **User Management**: Register, login, get user details.
5. **Search and Recommendation**: (Future enhancement)
6. **Borrowing and Returning Books**: (Future enhancement)

## Technologies Used
- **Programming Language**: Golang
- **Database**: PostgreSQL
- **Communication**: gRPC
- **Authentication**: JWT
- **Containerization**: Docker, Docker Compose
- **CI/CD**: GitHub Actions (example setup)
- **Cache**: Redis (for caching user sessions, etc.) (Future enhancement)

## Setup Instructions
1. **Clone the repository**:
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2. **Create `.env` file** in the root directory with the following environment variables:
    ```
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=your_database_name
    JWT_SECRET=your_secret_key
    ```

3. **Build and run using Docker Compose**:
    ```sh
    docker-compose up --build
    ```

4. **Access the services**:
    - BookService: `localhost:50051`
    - AuthorService: `localhost:50052`
    - CategoryService: `localhost:50053`
    - UserService: `localhost:50054`

## API Documentation
The API endpoints for each service are documented using Protocol Buffers (protobuf). To see the full list of available endpoints and their request/response formats, refer to the `.proto` files located in the `proto` directory.