# FlyPro Backend API (Go + Gin + GORM + Goose)

## Overview

This project is a backend API for **travel expense management**, built as part of the **FlyPro Backend Engineer Assessment**.

It demonstrates:

- Clean architecture with repository, service, and handler layers
- DTO-based request validation
- Goose for database migrations (instead of GORM AutoMigrate)
- Unit testing with mocks for the service layer
- RESTful API endpoints

---

## 🏗 Project Structure

cmd/
└── server/main.go # Entry point
internal/
├── config/ # DB, Redis, env configs
├── dto/ # Request/response DTOs
├── handlers/ # HTTP handlers
├── repository/ # DB access layer
├── services/ # Business logic
├── models/ # GORM models
└── utils/ # Helpers (error formatting, etc.)
migrations/ # Goose migrations
tests/ # Unit tests
docker-compose.yml # Postgres + Redis
Makefile # Migrations, tests, coverage
.env # Environment variables

---

## ⚙️ Setup & Run

### 1. Clone Repo & Install Deps

```bash
git clone https://github.com/yourusername/flypro-backend.git
cd flypro-backend
go mod tidy
```

## Setup & Run

### 2. Environment Variables

Create a `.env` file:

```env
DB_URL=postgres://[user]:[password]@localhost:5432/flypro?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
PORT=8080
```

### 3. Start Dependencies (Docker)

```
docker-compose up -d

```

### 4. Run Migrations

```
make migrate-up
make migrate-status


```

### 5. Run Server

```
go run ./cmd/server



```
