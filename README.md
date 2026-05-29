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
---

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
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE `coaches` (
  `id`         INT NOT NULL AUTO_INCREMENT,
  `uuid`       VARCHAR(255) DEFAULT (UUID()),
  `full_name`  VARCHAR(255) NOT NULL,
  `email`      VARCHAR(255) NOT NULL UNIQUE,
  `status`     ENUM('active', 'inactive') DEFAULT 'active',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `income_transactions` (
  `id`               INT NOT NULL AUTO_INCREMENT,
  `uuid`             VARCHAR(255) DEFAULT (UUID()),
  `coach_id`         INT NOT NULL,
  `amount`           DECIMAL(15,2) NOT NULL DEFAULT 0,
  `status`           ENUM('pending','completed','failed','refunded') DEFAULT 'completed',
  `transaction_date` DATE NOT NULL,
  `created_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`       TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_income_coach_date` (`coach_id`, `transaction_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

### 5. Seed sample data

```sql
INSERT INTO `coaches` (`uuid`, `full_name`, `email`, `status`) VALUES
(UUID(), 'Rahul Sharma', 'rahul.sharma@gmail.com', 'active'),
(UUID(), 'Priya Mehta',  'priya.mehta@gmail.com',  'active'),
(UUID(), 'Amit Verma',   'amit.verma@gmail.com',   'active'),
(UUID(), 'Sneha Patil',  'sneha.patil@gmail.com',  'inactive'),
(UUID(), 'Vikas Joshi',  'vikas.joshi@gmail.com',  'active');

INSERT INTO `income_transactions` (`uuid`, `coach_id`, `amount`, `status`, `transaction_date`) VALUES
(UUID(), 1, 80000,  'completed', CURDATE() - INTERVAL 6 DAY),
(UUID(), 1, 120000, 'completed', CURDATE() - INTERVAL 5 DAY),
(UUID(), 1, 95000,  'completed', CURDATE() - INTERVAL 4 DAY),
(UUID(), 1, 150000, 'completed', CURDATE() - INTERVAL 2 DAY),
(UUID(), 1, 200000, 'completed', CURDATE() - INTERVAL 1 DAY),
(UUID(), 1, 175000, 'completed', CURDATE() - INTERVAL 0 DAY),
(UUID(), 1, 60000,  'completed', CURDATE() - INTERVAL 13 DAY),
(UUID(), 1, 90000,  'completed', CURDATE() - INTERVAL 12 DAY),
(UUID(), 1, 75000,  'completed', CURDATE() - INTERVAL 11 DAY),
(UUID(), 1, 110000, 'completed', CURDATE() - INTERVAL 10 DAY),
(UUID(), 1, 85000,  'completed', CURDATE() - INTERVAL 9 DAY),
(UUID(), 1, 130000, 'completed', CURDATE() - INTERVAL 8 DAY),
(UUID(), 1, 100000, 'completed', CURDATE() - INTERVAL 7 DAY);
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

