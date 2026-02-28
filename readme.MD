# 🚀 Contest Rating Engine (Backend)

A robust, production-grade REST API designed to process and calculate competitive programming ratings[cite: 23]. This system utilizes a transparent percentile model to evaluate user performance against a participant pool, assigning standardized performance ratings and executing controlled rating adjustments[cite: 4, 5, 6].

Built with Go (Fiber) and PostgreSQL (via Prisma), the architecture strictly adheres to Clean Architecture principles, ensuring maximum scalability, testability, and ACID-compliant data integrity.

---

## 🧠 Core Mathematical Philosophy

The engine calculates rating changes via a strict 7-step deterministic model [cite: 8-14]:

1. **Calculate Beaten**: `Total Participants - Rank` [cite: 8]
2. **Determine Percentile**: `Beaten / Total Participants` [cite: 9]
3. **Identify Bracket**: Map the percentile to a predefined tier (e.g., Top 10%, Top 50%)[cite: 10].
4. **Assign Standard Performance**: Retrieve the fixed performance rating for that bracket[cite: 11].
5. **Compute Rating Change**: `(Performance Rating - Actual Rating) / 2`[cite: 12].
6. **Update Rating**: `Old Rating + Rating Change`[cite: 13].
7. **Tier Reassignment**: Automatically update the user's tier based on the new rating[cite: 14].

### Advanced Analytics Integration
While adhering to the 7-step ruleset, the service layer also models the rating adjustment as a **Discrete Kalman Filter** with a fixed gain of `K=0.5`. Additionally, the system calculates a **Logistic Upset Probability (Elo CDF)** to detect statistical anomalies in user performance, logging outlier jumps (probability < 5%) for potential anti-cheat review.

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
---

### 📐 System Design Principles Applied

To ensure the engine remains scalable, predictable, and cleanly implemented[cite: 45], the backend strictly adheres to industry-standard system design patterns:

* **Clean Architecture & Separation of Concerns (SoC):**
  The application is divided into distinct, isolated layers (Handlers $\rightarrow$ Services $\rightarrow$ Repositories). The HTTP routing layer knows nothing about the mathematical logic, and the mathematical logic knows nothing about the PostgreSQL database. This allows the system to scale infinitely without creating "spaghetti code."

* **Dependency Injection (DI) & Inversion of Control:**
  Global state is completely eliminated. The database client is injected into the Repository, the Repository is injected into the Service, and the Service is injected into the Handler via constructor functions (e.g., `NewRatingService(repo)`). This principle is what makes the application fully modular and allows for seamless unit testing by passing mock repositories.

* **The Repository Pattern:**
  All data access logic is abstracted behind an interface. The core business logic never interacts with Prisma or SQL directly. If the platform ever migrates from PostgreSQL to MongoDB or DynamoDB, only the `rating_repo.go` file needs to be rewritten. The rest of the application remains untouched.

* **ACID-Compliant Transactions:**
  Updating a user's rating involves writing to the `RatingHistory` table [cite: 31] and simultaneously updating the `Users` table's current rating and tier[cite: 32]. If a server crash occurs between these two writes, the ledger becomes corrupted. To prevent this, both queries are bundled into a single Prisma Atomic Transaction (`.Tx()`). It is an all-or-nothing operation: if one write fails, the entire batch instantly rolls back.

* **Single Responsibility Principle (SRP):**
  Every function and file has exactly one reason to change. The `rating_handler.go` solely validates JSON payloads. The `rating_service.go` solely executes the 7-step mathematical logic [cite: 8-14]. The `rating_repo.go` solely executes database queries. 

---

## 🛠 Tech Stack & Third-Party Dependencies

* **Language:** Go (Golang)
* **Database:** PostgreSQL (NeonDB Serverless) [cite: 22]

### Third-Party Go Imports:
Every package is intentionally selected to optimize performance and developer experience:

* **`github.com/gofiber/fiber/v2`**
  * *Purpose:* The core web framework[cite: 21]. Chosen for its extreme speed (built on `fasthttp`) and low memory footprint. Handles routing, JSON parsing, and HTTP middlewares (logger, panic recovery).
* **`github.com/steebchen/prisma-client-go`**
  * *Purpose:* The Prisma ORM client for Go. Provides a strictly-typed, auto-generated database client based on the `schema.prisma` file, enabling complex nested queries and ACID transactions without writing raw SQL strings.
* **`github.com/joho/godotenv`**
  * *Purpose:* Environment variable manager. Safely loads the `.env` file into the application's `os.Getenv` context during local development, explicitly bypassing itself in production environments.
* **`github.com/stretchr/testify`**
  * *Purpose:* Testing toolkit. Provides the `assert` package, which vastly improves the readability of Go's native testing suite by allowing clean, descriptive `assert.Equal()` checks in our table-driven unit tests.

---

## 🗄️ Database Schema

The PostgreSQL database is organized into three primary models[cite: 22, 33]:

1. **Users:** Tracks `id`, `name`, `current_rating`, `max_rating`, `contests_played`, and `tier`[cite: 34].
2. **Contests:** Stores `id`, `name`, `date`, and `total_participants`[cite: 35].
3. **RatingHistory:** A relational ledger recording `user_id`, `contest_id`, `old_rating`, `new_rating`, `performance_rating`, `rank`, `percentile`, and `rating_change` for every contest entry[cite: 36].

---

## 🔌 API Documentation

### 1. Process Contest Results
Calculates and applies rating updates for a batch of users after a contest ends [cite: 25-32].
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

The core math logic and percentile assignment rules are strictly validated via table-driven unit tests. This guarantees the 7-step logic remains mathematically sound through future updates. Code quality, documentation, and proper testing are mandatory requirements that this suite fulfills[cite: 46].

To execute the test suite:
```bash
go test ./internal/services/... -v
```