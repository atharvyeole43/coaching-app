# Coaching App — Daily Income Chart API

A backend REST API built with **Go (Gin)** for the Daily Income Chart feature in a coaching application. It compares the last 7 days of income versus the previous 7 days for a specific coach.

---

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** MySQL
- **ORM:** GORM
- **Logger:** Logrus

---

## Project Structure

coaching-app/
├── app/
│   └── app.go
├── cmd/
│   └── app/
│       └── run.go
├── config/
│   ├── config.go
│   └── db.go
├── constant/
├── dto/
│   └── income_dto.go
├── internal/
│   ├── apperrors/
│   │   └── apperrors.go
│   ├── controller/
│   │   └── inc/
│   │       ├── income_controller.go
│   │       └── init_income_controller.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   ├── rate_limiter.go
│   │   └── timeout_middleware.go
│   ├── repository/
│   │   └── inc/
│   │       ├── income_repository.go
│   │       └── init_income_repository.go
│   ├── routes/
│   │   └── routes.go
│   └── service/
│       └── inc/
│           ├── income_service.go
│           └── init_income_service.go
├── models/
│   ├── coach.go
│   └── income_transactions.go
├── pkg/
│   └── logger/
│       └── logger.go
├── response/
│   └── response.go
├── storage/
│   └── db/
│       └── client.go
├── utils/
│   ├── env_utils.go
│   ├── errors.go
│   ├── mysql_error_map.go
│   ├── response_utils.go
│   └── validator_utils.go
├── .env.example
├── go.mod
├── go.sum
└── main.go

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/atharvyeole43/coaching-app.git
cd coaching-app
```

### 2. Setup environment variables


APP_ENV="LOCAL"
APP_PORT=8082
DB_DRIVER="mysql"
DB_HOST="localhost"
DB_PORT="3306"
DB_USERNAME="root"

DB_DATABASE_NAME="coaching_app"    


DB_PASSWORD="Shriram@2001"


COACHING_APP_DB_MAX_IDLE_CONN=5
COACHING_APP_DB_MAX_OPEN_CONN=10
COACHING_APP_DB_MAX_IDLE_TIME=300    # 5 minutes
COACHING_APP_DB_MAX_LIFE_TIME=600    # 10 minutes```

### 3. Create the database

```sql
CREATE DATABASE coaching_db;
```

### 4. Run DDL — create tables

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE coach_status AS ENUM ('active', 'inactive');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed', 'refunded');

CREATE TABLE coaches (
    id          SERIAL PRIMARY KEY,
    uuid        UUID NOT NULL DEFAULT gen_random_uuid(),
    full_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    status      coach_status NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT uq_coaches_email UNIQUE (email)
);

CREATE TABLE income_transactions (
    id                SERIAL PRIMARY KEY,
    uuid              UUID NOT NULL DEFAULT gen_random_uuid(),
    coach_id          INT NOT NULL,
    amount            NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    status            transaction_status NOT NULL DEFAULT 'completed',
    transaction_date  DATE NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_income_coach FOREIGN KEY (coach_id) REFERENCES coaches(id)
);

CREATE INDEX idx_income_coach_date ON income_transactions (coach_id, transaction_date);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_coaches_updated_at
    BEFORE UPDATE ON coaches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_income_updated_at
    BEFORE UPDATE ON income_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 5. Seed sample data

```sql
INSERT INTO coaches (uuid, full_name, email, status) VALUES
(gen_random_uuid(), 'Rahul Sharma', 'rahul.sharma@gmail.com', 'active'),
(gen_random_uuid(), 'Priya Mehta',  'priya.mehta@gmail.com',  'active'),
(gen_random_uuid(), 'Amit Verma',   'amit.verma@gmail.com',   'active'),
(gen_random_uuid(), 'Sneha Patil',  'sneha.patil@gmail.com',  'inactive'),
(gen_random_uuid(), 'Vikas Joshi',  'vikas.joshi@gmail.com',  'active');

INSERT INTO income_transactions (uuid, coach_id, amount, status, transaction_date) VALUES
(gen_random_uuid(), 1, 80000,  'completed', CURRENT_DATE - INTERVAL '6 days'),
(gen_random_uuid(), 1, 120000, 'completed', CURRENT_DATE - INTERVAL '5 days'),
(gen_random_uuid(), 1, 95000,  'completed', CURRENT_DATE - INTERVAL '4 days'),
(gen_random_uuid(), 1, 150000, 'completed', CURRENT_DATE - INTERVAL '2 days'),
(gen_random_uuid(), 1, 200000, 'completed', CURRENT_DATE - INTERVAL '1 day'),
(gen_random_uuid(), 1, 175000, 'completed', CURRENT_DATE),
(gen_random_uuid(), 1, 60000,  'completed', CURRENT_DATE - INTERVAL '13 days'),
(gen_random_uuid(), 1, 90000,  'completed', CURRENT_DATE - INTERVAL '12 days'),
(gen_random_uuid(), 1, 75000,  'completed', CURRENT_DATE - INTERVAL '11 days'),
(gen_random_uuid(), 1, 110000, 'completed', CURRENT_DATE - INTERVAL '10 days'),
(gen_random_uuid(), 1, 85000,  'completed', CURRENT_DATE - INTERVAL '9 days'),
(gen_random_uuid(), 1, 130000, 'completed', CURRENT_DATE - INTERVAL '8 days'),
(gen_random_uuid(), 1, 100000, 'completed', CURRENT_DATE - INTERVAL '7 days');
```
### 6. Run the application

```bash
go run main.go
# or
make run
```

Server starts at `http://localhost:8080`

---

## API Reference

### GET /api/v1/income/daily

Returns daily income chart data comparing the last 7 days vs previous 7 days for a specific coach.

**Query Parameters**

| Parameter | Type | Required | Description |
|---|---|---|---|
| `coach_id` | integer | Yes | ID of the coach |
| `period` | string | No | Period to compare. Default: `7d` |

**Example Request**

GET /api/v1/income/daily?period=7d&coach_id=1

**Success Response — 200**

```json
{
  "status": true,
  "message": "Daily income fetched successfully",
  "data": {
    "summary": {
      "total_income_current": 820000,
      "total_income_previous": 650000,
      "delta_amount": 170000,
      "delta_percent": 26.15,
      "trend": "increased"
    },
    "chart": {
      "current": [
        { "date": "2024-04-01", "amount": 80000 },
        { "date": "2024-04-02", "amount": 120000 },
        { "date": "2024-04-03", "amount": 95000 },
        { "date": "2024-04-05", "amount": 150000 },
        { "date": "2024-04-06", "amount": 200000 },
        { "date": "2024-04-07", "amount": 175000 }
      ],
      "previous": [
        { "date": "2024-03-25", "amount": 60000 },
        { "date": "2024-03-26", "amount": 90000 },
        { "date": "2024-03-27", "amount": 75000 },
        { "date": "2024-03-28", "amount": 110000 },
        { "date": "2024-03-29", "amount": 85000 },
        { "date": "2024-03-30", "amount": 130000 },
        { "date": "2024-03-31", "amount": 100000 }
      ]
    },
    "meta": {
      "current_range": {
        "from": "2024-04-01",
        "to": "2024-04-07"
      },
      "previous_range": {
        "from": "2024-03-25",
        "to": "2024-03-31"
      }
    }
  }
}
```

**Error Responses**

| Status | Scenario | Message |
|---|---|---|
| 400 | Missing `coach_id` | `coach_id is required` |
| 400 | Invalid `coach_id` | `coach_id must be a valid number` |
| 404 | Coach not found | `Coach not found` |
| 500 | Server error | `Unexpected error occurred` |

---

## Test Scenarios

| Coach ID | Scenario | Expected Result |
|---|---|---|
| `1` | Data in both periods | Full chart with delta and trend |
| `2` | Current period only | `delta_percent: null`, trend `increased` |
| `3` | Mixed statuses | Only `completed` rows counted |
| `4` | Soft deleted coach | `404 Not Found` |
| `5` | No transactions | All zeros, trend `unchanged` |
| `999` | Does not exist | `404 Not Found` |

---

## Postman Collection

---
curl --location 'http://localhost:8082/api/v1/income/daily?period=7d&coach_id=3' \
--header 'Content-Type: application/json'

## Author

**Atharv Yeole**
GitHub: [@atharvyeole43](https://github.com/atharvyeole43)

