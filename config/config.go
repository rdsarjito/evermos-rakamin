package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	External ExternalConfig
}

type AppConfig struct {
	Host      string
	Port      string
	SecretKey string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ExternalConfig struct {
	APILocation string
}

var Cfg *Config

// Load membaca environment variables dan menginisialisasi Config
func Load() error {
	// Load .env file (optional, tidak error jika file tidak ada)
	_ = godotenv.Load()

	Cfg = &Config{
		App: AppConfig{
			Host:      getEnv("APP_HOST", "localhost"),
			Port:      getEnv("APP_PORT", "8080"),
			SecretKey: getEnv("SECRET_KEY", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "evermos_clone"),
		},
		External: ExternalConfig{
			APILocation: getEnv("API_LOCATION", ""),
		},
	}

	// Validasi required fields
	if Cfg.App.SecretKey == "" {
		return fmt.Errorf("SECRET_KEY is required")
	}

	return nil
}

// Get mengembalikan instance Config yang sudah di-load
func Get() *Config {
	if Cfg == nil {
		panic("config not loaded, call config.Load() first")
	}
	return Cfg
}

// getEnv membaca environment variable dengan default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt membaca environment variable sebagai integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// DSN mengembalikan Data Source Name untuk MySQL connection
func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}
