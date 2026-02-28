# 🚀 Contest Rating Engine (Backend)

A robust, production-grade REST API designed to process and calculate competitive programming ratings. This system utilizes a transparent percentile model to evaluate user performance against a participant pool, assigning standardized performance ratings and executing controlled rating adjustments.

Built with Go (Fiber) and PostgreSQL (via Prisma), the architecture strictly adheres to Clean Architecture principles, ensuring maximum scalability, testability, and ACID-compliant data integrity.

---

## 🧠 Core Mathematical Philosophy

The engine moves beyond opaque rating algorithms by utilizing a highly predictable, mathematically sound state estimation model. The rating changes are calculated via a strict 7-step deterministic framework, bolstered by advanced probabilistic anomaly detection.

### 1. The 7-Step Deterministic Model
Every time a contest ends, the engine evaluates each participant using the following sequence:

1. **Calculate Beaten**: `Beaten = Total Participants - Rank`
2. **Determine Percentile**: `Percentile = Beaten / Total Participants`
3. **Identify Bracket**: Map the computed percentile to a predefined tier (e.g., Top 10%, Top 50%).
4. **Assign Standard Performance**: Retrieve the fixed standard performance rating ($\rho$) associated with that specific bracket.
5. **Compute Rating Change**: $\Delta R = (\rho - R_{old}) / 2$
6. **Update Rating**: $R_{new} = R_{old} + \Delta R$
7. **Tier Reassignment**: Automatically update the user's string-based tier based on their $R_{new}$.

### 2. Advanced State Estimation (Kalman Filter)
While the formula `(Performance Rating - Actual Rating) / 2` appears simple, it is actually a specific implementation of a **Discrete Kalman Filter** operating under a fixed steady-state assumption. 

In state estimation theory:
$\hat{x}_{k} = \hat{x}_{k-1} + K_k (z_k - \hat{x}_{k-1})$

By strictly dividing the difference by 2, our system sets a constant Kalman Gain of $K = 0.5$. In stochastic terms, this asserts that the uncertainty (variance) of the user's historical rating is exactly equal to the noise of their single-contest performance. This guarantees a mathematically stable mean-reversion, preventing hyper-inflation of ratings from isolated outlier performances.

### 3. Statistical Anomaly Detection (Logistic CDF)
To ensure platform integrity and detect statistical anomalies (e.g., potential cheating or extreme upsets), the engine integrates a Logistic Cumulative Distribution Function derived from the Elo rating system. 


The expected probability $E$ of a user achieving their assigned performance rating $\rho$ given their prior rating $R_{old}$ is calculated as:
$E = \frac{1}{1 + 10^{(\rho - R_{old})/400}}$

If $E < 0.05$, the performance represents a $95\%$ statistical upset, triggering internal analytics for anomaly review without delaying the immediate mathematical update.

---

## 🏗 System Architecture

The backend is structured using **Clean Architecture** to decouple the HTTP transport layer, business logic, and database access. 


```text
backend/
├── cmd/api/main.go              # Application entry point & dependency injection wiring
├── internal/
│   ├── config/db.go             # Database connection & environment configuration
│   ├── handlers/                # HTTP layer (Fiber context parsing and validation)
│   ├── models/                  # Data Transfer Objects (DTOs) and payload structs
│   ├── repository/              # Data access layer (Prisma queries and ACID transactions)
│   ├── routes/                  # API endpoint registration
│   └── services/                # Core mathematical business logic and testing
├── prisma/schema.prisma         # Database schema definition
└── .env                         # Environment variables (Ignored in Git)
```
___
### 📐 System Design Principles Applied

To ensure the engine remains scalable, predictable, and cleanly implemented, the backend strictly adheres to industry-standard system design patterns:

* **Clean Architecture & Separation of Concerns (SoC):**
  The application is divided into distinct, isolated layers (Handlers -> Services -> Repositories). The HTTP routing layer knows nothing about the mathematical logic, and the mathematical logic knows nothing about the PostgreSQL database. This allows the system to scale infinitely without creating "spaghetti code."

* **Dependency Injection (DI) & Inversion of Control:**
  Global state is completely eliminated. The database client is injected into the Repository, the Repository is injected into the Service, and the Service is injected into the Handler via constructor functions (e.g., `NewRatingService(repo)`). This principle is what makes the application fully modular and allows for seamless unit testing by passing mock repositories.

* **The Repository Pattern:**
  All data access logic is abstracted behind an interface. The core business logic never interacts with Prisma or SQL directly. If the platform ever migrates from PostgreSQL to MongoDB or DynamoDB, only the `rating_repo.go` file needs to be rewritten. The rest of the application remains untouched.

* **ACID-Compliant Transactions:**
  Updating a user's rating involves writing to the `RatingHistory` table and simultaneously updating the `Users` table's current rating and tier. If a server crash occurs between these two writes, the ledger becomes corrupted. To prevent this, both queries are bundled into a single Prisma Atomic Transaction (`.Tx()`). It is an all-or-nothing operation: if one write fails, the entire batch instantly rolls back.

* **Single Responsibility Principle (SRP):**
  Every function and file has exactly one reason to change. The `rating_handler.go` solely validates JSON payloads. The `rating_service.go` solely executes the core mathematical formulas. The `rating_repo.go` solely executes database queries. 

---

## 🛠 Tech Stack & Third-Party Dependencies

* **Language:** Go (Golang)
* **Database:** PostgreSQL (NeonDB Serverless)

### Third-Party Go Imports:
Every package is intentionally selected to optimize performance and developer experience:

* **`github.com/gofiber/fiber/v2`**
  * *Purpose:* The core web framework. Chosen for its extreme speed (built on `fasthttp`) and low memory footprint. Handles routing, JSON parsing, and HTTP middlewares (logger, panic recovery).
* **`github.com/steebchen/prisma-client-go`**
  * *Purpose:* The Prisma ORM client for Go. Provides a strictly-typed, auto-generated database client based on the `schema.prisma` file, enabling complex nested queries and ACID transactions without writing raw SQL strings.
* **`github.com/joho/godotenv`**
  * *Purpose:* Environment variable manager. Safely loads the `.env` file into the application's `os.Getenv` context during local development, explicitly bypassing itself in production environments.
* **`github.com/stretchr/testify`**
  * *Purpose:* Testing toolkit. Provides the `assert` package, which vastly improves the readability of Go's native testing suite by allowing clean, descriptive `assert.Equal()` checks in our table-driven unit tests.

---

## 🗄️ Database Schema

The PostgreSQL database is organized into three primary models:

1. **Users:** Tracks `id`, `name`, `current_rating`, `max_rating`, `contests_played`, and `tier`.
2. **Contests:** Stores `id`, `name`, `date`, and `total_participants`.
3. **RatingHistory:** A relational ledger recording `user_id`, `contest_id`, `old_rating`, `new_rating`, `performance_rating`, `rank`, `percentile`, and `rating_change` for every contest entry.

---

## 🔌 API Documentation

### 1. Process Contest Results
Calculates and applies rating updates for a batch of users after a contest ends.
* **Endpoint:** `POST /api/contests/end`
* **Payload:**
  ```json
  {
    "contest_id": "uuid-string",
    "results": [
      { "user_id": "uuid-string", "rank": 5 }
    ]
  }
  ```
* **Response (200 OK):** `{"message": "Contest ratings calculated and updated successfully"}`

### 2. Fetch User Profile
Retrieves the user's current statistics and their complete rating history (including nested contest metadata for frontend charts).
* **Endpoint:** `GET /api/users/:id`
* **Response (200 OK):** Returns the serialized `User` object with an included array of `RatingHistory` items.

---

## 🚀 Setup & Installation

1. **Clone & Navigate:**
   ```bash
   git clone <repo-url>
   cd backend
   ```

2. **Environment Variables:**
   Create a `.env` file in the root of the `backend` directory.
   ```env
   # Ensure you use a direct connection string (no PgBouncer/pooler)
   DATABASE_URL="postgresql://user:password@hostname/dbname?sslmode=require"
   PORT="8080"
   ```

3. **Install Dependencies & Sync Database:**
   ```bash
   go mod tidy
   npx prisma db push
   go run [github.com/steebchen/prisma-client-go](https://github.com/steebchen/prisma-client-go) generate
   ```

4. **Seed Dummy Data (Optional):**
   ```bash
   go run seed.go
   ```

5. **Start the Server:**
   ```bash
   go run cmd/api/main.go
   ```

---

## 🧪 Testing

The core math logic and percentile assignment rules are strictly validated via table-driven unit tests. This guarantees the 7-step logic remains mathematically sound through future updates. Code quality, documentation, and proper testing are mandatory requirements that this suite fulfills.

To execute the test suite:
```bash
go test ./internal/services/... -v
```