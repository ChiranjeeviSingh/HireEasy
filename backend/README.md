# README.md for Hiring Easy Backend

## Features

- User authentication for hiring managers
- Job Management
- Form Template Management
- Application Form System

## Setup Instructions

1. Clone the repository:
   ```
   git clone <repository-url>
   ```

2. Navigate to the backend directory:
   ```
   cd HireEasy/backend
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Set up the PostgreSQL database using the provided schema:
   ```
   psql -U <username> -d <database> -f internal/database/migrations/schema.sql
   ```

5. Run the application:
   ```
   go run cmd/server/main.go
   ```

6. Run unittests:
   ```
   go test ./test/handlers/ -v
   ```