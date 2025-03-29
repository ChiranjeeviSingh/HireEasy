package database

import (
	"fmt"
	"log"

	"backend/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import PostgreSQL driver
)

var db *sqlx.DB

// Connect initializes the database connection
func Connect() {
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatal("ERROR: Config is nil. Database configuration is not loaded properly.")
	}

	// Debugging: Print database config values
	dbConfig := cfg.DBConfig
	log.Printf("🔍 DB Config -> Host: %s, Port: %s, User: %s, DB: %s", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Dbname)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Dbname,
	)

	var err error
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to reach database: %v", err)
	}

	log.Println("✅ Successfully connected to the database")
}

// GetDB returns the active database connection
func GetDB() *sqlx.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	return db
}
