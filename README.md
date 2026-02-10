# Digital Wallet App

A simple and solid digital wallet application built with Go following **Clean Architecture** principles.

## Features

- **Withdraw API**: Withdraw funds from user's wallet with balance validation and concurrency safety.
- **Balance Inquiry API**: Check current wallet balance.
- **Data Consistency**: Prevents negative balances using database transactions and row-level locking (`FOR UPDATE`).
- **Clean Architecture**: Decoupled layers (Domain, Usecase, Repository, Delivery) for better maintainability and testability.

## Technology Stack

- **Go** (Golang)
- **PostgreSQL** (Database)
- **Gin Web Framework** (HTTP Router & Middleware)

## Project Structure

```
api/             # Server entry point
internal/
  domain/        # Entities and Interface definitions (Core Logic)
  usecase/       # Business Logic (Application Layer)
  handler/       # HTTP Handlers (Web Layer)
  repository/    # Database Operations (Data Layer)
  config/        # Configuration management
  database/      # Database connection setup
migrations/      # Database schema
```

## How to Run

### 1. Prerequisites

- Go 1.23+

### 2. Setup Database


The database will be available at `localhost:5432`.

### 3. Setup Schema

Connect to the database and run the SQL script in `migrations/000001_init_schema.up.sql`.
If you have `psql` installed:

```bash
psql -h localhost -U postgres -d digital_wallet -f migrations/000001_init_schema.up.sql
```

_(Default password is `postgres`)_

### 4. Configuration

Review `.env` file :

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=digital_wallet
DB_SSLMODE=disable
APP_PORT=8080
```

### 5. Run Application

```bash
go run cmd/api/main.go
```

## API Documentation

### 1. Balance Inquiry

- **Endpoint**: `GET /api/balance/:user_id`
- **Example**: `GET http://localhost:8080/api/balance/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11`
- **Response**:

```json
{
  "message": "success",
  "balance": 1000.0,
  "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
}
```

### 2. Withdraw

- **Endpoint**: `POST /api/withdraw`
- **Body**:

```json
{
  "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "amount": 50.0
}
```

- **Response**:

```json
{
  "message": "withdrawal successful",
  "new_balance": 950.0,
  "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
}
```

### 3. Error Response

All error responses will have the same format:

```json
{
  "message": "error description"
}
```

## Key Considerations

- **Concurrency**: Uses `SELECT ... FOR UPDATE` to ensure that parallel withdrawal requests for the same user don't cause race conditions.
- **Error Handling**: Proper mapping of internal domain errors to HTTP status codes.
- **Dependency Injection**: Dependencies are injected in `main.go`, making it easy to swap implementations (e.g., Mock Repository for testing).
