package test

import (
	"backend/internal/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"backend/internal/database"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
	"math/rand"
)

// SetupTestDB initializes a test database
func SetupTestDB() *sql.DB {
	// Load config first
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.GetConfig().DBConfig

	// First connect to default database to create test database if it doesn't exist
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Use a separate test database
	testDbName := "app_db_test"
	
	// Create test database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDbName))
	if err != nil && !isDatabaseExistsError(err) {
		log.Printf("Warning: Could not create test database: %v", err)
	}

	// Set the global database connection
	os.Setenv("DB_NAME", testDbName)
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	database.Connect()

	// Run migrations on test database
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations on test database: %v", err)
	}

	return database.GetDB().DB
}

// dropExistingTables drops all tables in the correct order to handle foreign key constraints
func dropExistingTables(db *sqlx.DB) error {
	// Drop tables in reverse order of dependencies
	dropStatements := []string{
		"DROP TABLE IF EXISTS application_form CASCADE;",
		"DROP TABLE IF EXISTS form_templates CASCADE;",
		"DROP TABLE IF EXISTS jobs CASCADE;",
		"DROP TABLE IF EXISTS users CASCADE;",
	}

	for _, stmt := range dropStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to drop table: %v\nStatement: %s", err, stmt)
		}
	}

	return nil
}

// runMigrations executes the schema.sql file on the test database
func runMigrations(db *sqlx.DB) error {
	// Drop existing tables first
	if err := dropExistingTables(db); err != nil {
		return fmt.Errorf("failed to drop existing tables: %v", err)
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	// Go up two levels from the test directory to reach the backend root
	backendRoot := filepath.Dir(filepath.Dir(cwd))

	// Read the schema file using absolute path
	schemaPath := filepath.Join(backendRoot, "internal", "database", "migrations", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// Split the SQL file into individual statements
	statements := strings.Split(string(schemaSQL), ";")

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %v\nStatement: %s", err, stmt)
		}
	}

	return nil
}

// isDatabaseExistsError checks if the error is due to database already existing
func isDatabaseExistsError(err error) bool {
	return err.Error() == "pq: database \"app_db_test\" already exists"
}

// InsertTestUser inserts a dummy user for authentication tests
func InsertTestUser(db *sql.DB) (userID int, token string) {
	// Hash the password 'password123'
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	query := `INSERT INTO users (email, password_hash, username, role, company_name) VALUES ('test@example.com', $1, 'testuser', 'HR', 'Test Company') RETURNING id`
	err = db.QueryRow(query, string(hashedPassword)).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to insert test user: %v", err)
	}

	token = "dummy-hr-jwt-token" // Simulate a valid token
	return
}

// InsertTestInterviewerUser inserts a dummy interviewer user for authentication tests
func InsertTestInterviewerUser(db *sql.DB) (userID int, token string) {
	// Hash the password 'password123'
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	rand.Seed(time.Now().UnixNano())
    r := rand.Intn(100000) 

	query := fmt.Sprintf(`INSERT INTO users (email, password_hash, username, role, company_name) VALUES ('interviewer%d@example.com', $1, 'interviewer_%d', 'INTERVIEWER', 'Test Company') RETURNING id`,  r, r)
	err = db.QueryRow(query, string(hashedPassword)).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to insert test interviewer: %v", err)
	}

	token = "dummy-interviewer-jwt-token" // Simulate a valid token
	return
}

// InsertTestData inserts test data before a test
func InsertTestData(db *sql.DB) string {
	formUUID := "986c54c7-c3c5-4f97-a952-5dd7013ced3f"

	_, err := db.Exec(`
        INSERT INTO application_form (form_uuid, job_id, form_id, status, date_created)
        VALUES ($1, 1, 1, 'active', NOW())`, formUUID)
	if err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	return formUUID
}

// CleanupTestDB removes all test data after a test
func CleanupTestDB(db *sql.DB) {
	_, err := db.Exec("DELETE FROM application_form; DELETE FROM jobs; DELETE FROM form_templates; DELETE FROM users; TRUNCATE users, jobs, job_submissions, availabilities, interviews CASCADE;")
	if err != nil {
		log.Fatalf("Failed to clean up test DB: %v", err)
	}
}

