package db

import "os"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "admin"),
		Password: getEnv("POSTGRES_PASSWORD", "admin123"),
		Name:     getEnv("POSTGRES_DB", "filedb"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
