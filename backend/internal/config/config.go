package config

import (
	"os"
)

type Config struct {
    AWSRegion   string
    AWSEndpoint string
    JWTSecret string
    DBConfig postgresConfig
}

type postgresConfig struct {
    Host string
    Port string
    Username string
    Password string
    Dbname string
}

var globalConfig *Config

func LoadConfig() error {
    globalConfig = &Config{
        AWSRegion:   os.Getenv("AWS_REGION"),
        AWSEndpoint: os.Getenv("AWS_ENDPOINT"),
        JWTSecret: "abcdefghijklmno",
        DBConfig: postgresConfig{
            Host: "localhost",
            Port: "5432",
            Username: "postgres",
            Password: "123456",
            Dbname: "go_db",
        },
    }
    return nil
}

func GetConfig() *Config {
    return globalConfig
}
