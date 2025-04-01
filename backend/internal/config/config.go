package config

import (
	"log"
	"os"
)

type Config struct {
	AWSRegion   string
	AWSEndpoint string
	S3Bucket    string
	JWTSecret   string
	TestMode    bool
	DBConfig    postgresConfig
	CORSConfig  corsConfig
}

type postgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
}

type corsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

var globalConfig *Config

func LoadConfig() error {
	// Load test mode from environment
	testMode := os.Getenv("TEST_MODE") == "true" || os.Getenv("S3_TEST_MODE") == "true"

	globalConfig = &Config{
		AWSRegion:   getEnvOrDefault("AWS_REGION", "us-east-1"),
		AWSEndpoint: getEnvOrDefault("AWS_ENDPOINT", ""),
		S3Bucket:    getEnvOrDefault("S3_BUCKET", "hireeasy-resumes"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "abcdefghijklmno"),
		TestMode:    testMode,
		DBConfig: postgresConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			Username: getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "123456"),
			Dbname:   getEnvOrDefault("DB_NAME", "go_db"),
		},
		CORSConfig: corsConfig{
			AllowedOrigins:   []string{getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: getEnvOrDefault("CORS_ALLOW_CREDENTIALS", "true") == "true",
		},
	}

	// Debugging: Print loaded configuration
	log.Printf("✅ Config Loaded: Host=%s, Port=%s, User=%s, DB=%s",
		globalConfig.DBConfig.Host,
		globalConfig.DBConfig.Port,
		globalConfig.DBConfig.Username,
		globalConfig.DBConfig.Dbname)

	log.Printf("🌍 CORS Settings: AllowedOrigins=%v, AllowedMethods=%v, AllowedHeaders=%v, AllowCredentials=%v",
		globalConfig.CORSConfig.AllowedOrigins,
		globalConfig.CORSConfig.AllowedMethods,
		globalConfig.CORSConfig.AllowedHeaders,
		globalConfig.CORSConfig.AllowCredentials,
	)

	return nil
}

func GetConfig() *Config {
	if globalConfig == nil {
		panic("❌ ERROR: Config is not initialized. Did you forget to call LoadConfig()?")
	}
	return globalConfig
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
