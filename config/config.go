package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DefaultMaxPlayers       int
	HashCost                int
	JwtDurationInMinutes    int
	JwtRenewDueInMinutes    int
	JwtSecret               string
	ServerAddress           string
	SQLiteDatabasePath      string
	WebSocketAllowedOrigins string
}

var config *Config = nil

func GetConfig() Config {
	return *config
}

func getStringEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)

	if err != nil {
		panic("failed to parse int from " + key + " env: " + err.Error())
	}

	return parsed
}

func init() {
	err := godotenv.Load()

	if err != nil {
		panic("failed to load .env file")
	}

	config = &Config{
		DefaultMaxPlayers:       getIntEnv("DEFAULT_MAX_PLAYERS", 6),
		HashCost:                getIntEnv("HASH_COST", 14),
		JwtDurationInMinutes:    getIntEnv("JWT_DURATION_IN_MINUTES", 60),
		JwtRenewDueInMinutes:    getIntEnv("JWT_RENEW_DUE_IN_MINUTES", 300),
		JwtSecret:               getStringEnv("JWT_SECRET", "secret"),
		ServerAddress:           getStringEnv("SERVER_ADDRESS", ":8080"),
		SQLiteDatabasePath:      getStringEnv("SQLITE_DATABASE_PATH", "database.sqlite3"),
		WebSocketAllowedOrigins: getStringEnv("WEBSOCKET_ALLOWED_ORIGINS", "*"),
	}
}
