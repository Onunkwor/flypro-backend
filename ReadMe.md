<!-- # FlyPro Backend API

## 📌 Overview

This project is a backend API for Travel Expense Management, built for the FlyPro Backend Engineer Assessment. It demonstrates:

- Clean architecture with repository, service, and handler layers
- Request validation using DTOs and custom validators
- Currency conversion integration with USD normalization
- Redis caching for performance
- Goose for database migrations
- Unit testing with mocks

## 🏗 Project Structure

```
cmd/
└── server/main.go          # Entry point
internal/
├── config/                 # DB, Redis, env configs
├── dto/                    # Request/response DTOs
├── handlers/               # HTTP handlers
├── repository/             # Data access layer
├── services/               # Business logic
├── models/                 # GORM models
├── middleware/             # Logging, CORS, rate limiting
└── utils/                  # Helpers (error formatting, etc.)
migrations/                 # Goose migration files
tests/                      # Unit tests
docker-compose.yml          # Postgres + Redis setup
Makefile                    # Migration & test commands
.env                        # Environment variables
```

## ⚙️ Setup & Run

### 1. Clone Repo & Install Dependencies

```bash
git clone https://github.com/yourusername/flypro-backend.git
cd flypro-backend
go mod tidy
```

### 2. Environment Variables

Create a `.env` file in the project root:

```
DB_URL=postgres://[user]:[password]@localhost:5432/flypro?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
PORT=8080
CURRENCY_API_KEY=your_api_key_here
```

### 3. Start Dependencies

```bash
docker-compose up -d
```

### 4. Run Migrations

```bash
make migrate-up
make migrate-status
```

### 5. Run Server

```bash
go run ./cmd/server
```

## 🗄 Database Schema

- **User**: Manages user accounts
- **Expense**: Tracks expenses with original Amount + Currency and computed AmountUSD
- **ExpenseReport**: Groups multiple expenses, stores Total in USD

**Design Decision**: I chose to persist both the original amount + currency and a converted AmountUSD. This preserves data integrity while enabling USD-based reporting.

## 📡 API Endpoints

### Users

- `POST /api/users` – Create user
- `GET /api/users/:id` – Get user details

### Expenses

- `POST /api/expenses` – Create expense
- `GET /api/expenses` – List expenses (pagination, filters)
- `GET /api/expenses/:id` – Get expense details
- `PUT /api/expenses/:id` – Update expense
- `DELETE /api/expenses/:id` – Delete expense

### Reports

- `POST /api/reports` – Create report
- `POST /api/reports/:id/expenses` – Add expenses to report
- `GET /api/reports` – List reports (pagination)
- `PUT /api/reports/:id/submit` – Submit report

## 🔄 Currency Conversion & Caching

- Integrated with a third-party currency API
- All expenses normalized to USD in reports
- Cached exchange rates in Redis (6-hour TTL)
- Report summaries cached in Redis (30-minute TTL)

## 🧪 Testing

Service-layer testing with table-driven tests and mocks:

- **SubmitReport**: Validates ownership & state transitions
- **ListReports**: Converts expenses to USD and calculates totals
- **Error handling**: Handles currency API failures gracefully

**Note**: Focused on reporting tests instead of CreateExpense, as the design stores both original values and computed AmountUSD, showcasing USD conversion at the reporting stage.

Run tests with:

```bash
make test
```

## ❓ Design & Scalability Questions

### Concurrent Expense Approvals

- Use row-level locking (Postgres `FOR UPDATE`) or optimistic locking with versioning
- Prevents race conditions when multiple approvers act on the same report

### Scaling Strategies

- Horizontal scaling with load balancers
- Read replicas for PostgreSQL
- Redis caching for hot queries

### Ensuring Data Consistency

- Database transactions for multi-step operations
- Event-driven architecture (Kafka/RabbitMQ) for eventual consistency if services are decoupled

### Monitoring & Alerting

- Prometheus + Grafana for metrics (API latency, DB health, Redis hit rate)
- Alertmanager for error thresholds (e.g., currency API failures)
- Structured logging with correlation IDs

## ✅ Submission Deliverables

- Clean Git history & commits
- Dockerized local setup (Postgres + Redis)
- Goose migrations with indexes & constraints
- Postman collection for API testing
- Seed script with demo data

## 🚀 Bonus Features Implemented

- Structured logging
- Redis caching for report summaries
- Graceful error handling with custom errors -->
