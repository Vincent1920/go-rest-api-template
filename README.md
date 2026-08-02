# Go Todo REST API

A REST API for managing Todo data, built with Go, Gin, PostgreSQL, GORM, Docker, JWT Authentication, and Clean Architecture.

## Features

- User registration
- User login
- JWT authentication
- Protected user profile endpoint
- Password hashing with bcrypt
- PostgreSQL database
- GORM Auto Migration
- Docker Compose
- pgAdmin
- Clean Architecture
- Configuration using YAML
- Todo CRUD
- Pagination
- Search
- Swagger documentation

## Tech Stack

- Go
- Gin
- PostgreSQL
- GORM
- Docker
- Docker Compose
- JWT
- bcrypt
- Viper
- Swagger

## Project Structure

```text
todo-api/
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   ├── config.go
│   └── config.yaml
├── docs/
├── internal/
│   ├── database/
│   ├── domain/
│   ├── dto/
│   ├── handler/
│   ├── middleware/
│   ├── repository/
│   ├── routes/
│   ├── service/
│   └── utils/
├── migrations/
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md