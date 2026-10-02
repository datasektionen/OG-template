package config

import (
	"os"
)

type EnvVar struct {
	ServerPort string
	ServerURL  string
	DBConn     string
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetEnv() EnvVar {
	return EnvVar{
		ServerPort: getEnvOrDefault("SERVER_PORT", "3000"),
		ServerURL:  getEnvOrDefault("SERVER_URL", "http://localhost:3000"),
		DBConn:     os.Getenv("DB_CONN"),
	}
}
