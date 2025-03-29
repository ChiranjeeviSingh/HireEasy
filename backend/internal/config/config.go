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
}

type postgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
}

var globalConfig *Config

// func LoadConfig() error {
// 	globalConfig = &Config{
// 		AWSRegion:   os.Getenv("AWS_REGION"),
// 		AWSEndpoint: os.Getenv("AWS_ENDPOINT"),
// 		JWTSecret:   "abcdefghijklmno",
// 		DBConfig: postgresConfig{
// 			Host:     "localhost",
// 			Port:     "5432",
// 			Username: "reshma",
// 			Password: "postgres",
// 			Dbname:   "app_db",
// 		},
// 	}
// 	return nil
// }

// func GetConfig() *Config {
//     return globalConfig
// }

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
			Username: getEnvOrDefault("DB_USER", "reshma"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			Dbname:   getEnvOrDefault("DB_NAME", "app_db"),
		},
	}

	// Debugging: Print loaded configuration
	log.Printf("✅ Config Loaded: Host=%s, Port=%s, User=%s, DB=%s",
		globalConfig.DBConfig.Host,
		globalConfig.DBConfig.Port,
		globalConfig.DBConfig.Username,
		globalConfig.DBConfig.Dbname)

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
