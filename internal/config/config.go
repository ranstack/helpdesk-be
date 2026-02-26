package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() *Config {
	return &Config{
		AppName: getEnv("APP_NAME", "Helpdesk API"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "helpdesk"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}
	return env
}
