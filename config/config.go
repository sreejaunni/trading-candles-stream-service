package config

import (
	"os"
	"strconv"
)

const (
	BTCUSDT  = "btcusdt"
	ETHUSDT  = "ethusdt"
	PEPEUSDT = "pepeusdt"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Symbols []string
}

func NewConfig() *Config {

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":50052"),
		},
		DB: DBConfig{
			User:     getEnv("DB_USER", "app_user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "market_db"),
		},
		Symbols: []string{BTCUSDT, ETHUSDT, PEPEUSDT},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
