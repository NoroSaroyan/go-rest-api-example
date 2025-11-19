// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	Log LogConfig
}

type AppConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type LogConfig struct {
	Level string
}

// Load reads configuration from environment variables and validates them.
func Load() (*Config, error) {
	// Load .env file if it exists (silently ignore if it doesn't)
	_ = godotenv.Load()

	cfg := &Config{}

	if err := cfg.loadAppConfig(); err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	if err := cfg.loadDBConfig(); err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	if err := cfg.loadLogConfig(); err != nil {
		return nil, fmt.Errorf("failed to load log config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) loadAppConfig() error {
	var err error

	c.App.Port = getEnv("APP_PORT", "8080")

	if c.App.ReadTimeout, err = parseDuration("APP_READ_TIMEOUT", "10s"); err != nil {
		return err
	}

	if c.App.WriteTimeout, err = parseDuration("APP_WRITE_TIMEOUT", "10s"); err != nil {
		return err
	}

	if c.App.IdleTimeout, err = parseDuration("APP_IDLE_TIMEOUT", "60s"); err != nil {
		return err
	}

	return nil
}

func (c *Config) loadDBConfig() error {
	var err error

	c.DB.Host = getEnv("DB_HOST", "localhost")
	c.DB.User = getEnv("DB_USER", "todo")
	c.DB.Password = getEnv("DB_PASSWORD", "todo")
	c.DB.Name = getEnv("DB_NAME", "todo_db")

	if c.DB.Port, err = parseInt("DB_PORT", "5432"); err != nil {
		return err
	}

	if c.DB.MaxOpenConns, err = parseInt("DB_MAX_OPEN_CONNS", "25"); err != nil {
		return err
	}

	if c.DB.MaxIdleConns, err = parseInt("DB_MAX_IDLE_CONNS", "25"); err != nil {
		return err
	}

	if c.DB.ConnMaxLifetime, err = parseDuration("DB_CONN_MAX_LIFETIME", "5m"); err != nil {
		return err
	}

	if c.DB.ConnMaxIdleTime, err = parseDuration("DB_CONN_MAX_IDLE_TIME", "5m"); err != nil {
		return err
	}

	return nil
}

func (c *Config) loadLogConfig() error {
	c.Log.Level = getEnv("LOG_LEVEL", "info")
	return nil
}

func (c *Config) validate() error {
	// Validate app port
	if port, err := strconv.Atoi(c.App.Port); err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid APP_PORT: must be a number between 1 and 65535")
	}

	// Validate database port
	if c.DB.Port < 1 || c.DB.Port > 65535 {
		return fmt.Errorf("invalid DB_PORT: must be between 1 and 65535")
	}

	// Validate required database fields
	if strings.TrimSpace(c.DB.Host) == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if strings.TrimSpace(c.DB.User) == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if strings.TrimSpace(c.DB.Name) == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}

	if !validLogLevels[strings.ToLower(c.Log.Level)] {
		return fmt.Errorf("invalid LOG_LEVEL: must be one of debug, info, warn, error, fatal")
	}

	return nil
}

// DatabaseURL returns a properly formatted PostgreSQL connection string.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(c.DB.User),
		url.QueryEscape(c.DB.Password),
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func parseInt(key, defaultValue string) (int, error) {
	val := getEnv(key, defaultValue)
	result, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return result, nil
}

func parseDuration(key, defaultValue string) (time.Duration, error) {
	val := getEnv(key, defaultValue)
	result, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return result, nil
}
