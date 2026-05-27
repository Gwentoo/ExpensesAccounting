package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database     PostgresConfig
	Redis        RedisConfig
	SuperSet     SuperSetConfig
	Auth         ServiceConfig
	Expenses     ServiceConfig
	GatewayPort  string
	Google       GoogleConfig
	JwtSecretKey string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type GoogleConfig struct {
	GooglePass string
	From       string
}
type RedisConfig struct {
	Addr string
}

type SuperSetConfig struct {
	DashboardID string
	Username    string
	Password    string
	Port        int
}

type ServiceConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func LoadConfig() *Config {

	cfg := &Config{
		Database: PostgresConfig{
			Host:     getEnvString("DB_HOST", "postgres"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnvString("POSTGRES_USER", "postgres"),
			Password: getEnvString("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnvString("POSTGRES_DB", "expense"),
			SSLMode:  getEnvString("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr: getEnvString("REDIS_ADDR", "redis:6379"),
		},
		SuperSet: SuperSetConfig{
			Username: os.Getenv("SUPERSET_USERNAME"),
			Password: os.Getenv("SUPERSET_PASSWORD"),
			Port:     getEnvInt("SUPERSET_EXTERNAL_PORT", 8088),
		},
		Auth: ServiceConfig{
			Host:    getEnvString("AUTH_HOST", "auth-service"),
			Port:    getEnvInt("AUTH_SERVICE_PORT", 50051),
			Timeout: 10 * time.Second,
		},
		Expenses: ServiceConfig{
			Host:    getEnvString("EXPENSES_HOST", "expenses-service"),
			Port:    getEnvInt("EXPENSES_SERVICE_PORT", 50052),
			Timeout: 10 * time.Second,
		},
		Google: GoogleConfig{
			GooglePass: os.Getenv("GOOGLE_PASS"),
			From:       os.Getenv("GOOGLE_FROM"),
		},
		GatewayPort:  getEnvString("GATEWAY_EXTERNAL_PORT", "8080"),
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	return cfg
}

func getEnvString(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func (p *PostgresConfig) GetPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}
