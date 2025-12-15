package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	_        struct{}
	Host     string
	Port     uint16
	User     string
	Password string
	Database string // Database name
	Params   map[string]string

	MaxOpenConnections int
	MaxIdleConnections int

	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewConfig(opts ...ConnectionConfig) (*Config, error) {
	cfg := &Config{
		Params: make(map[string]string),
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// cfg := Config{} // forbidden
// cfg.WithHost().With_ConnMaxLifetime().With_ConnMaxIdleTime()...().FormatDSN()

type ConnectionConfig func(*Config) error

func WithHost(host string) ConnectionConfig {
	return func(config *Config) error {
		config.Host = host
		return nil
	}
}

func WithPort(port uint16) ConnectionConfig {
	return func(config *Config) error {
		config.Port = port
		return nil
	}
}

func WithUser(user string) ConnectionConfig {
	return func(config *Config) error {
		config.User = user
		return nil
	}
}

func WithPassword(password string) ConnectionConfig {
	return func(config *Config) error {
		config.Password = password
		return nil
	}
}

func WithDatabase(database string) ConnectionConfig {
	return func(config *Config) error {
		config.Database = database
		return nil
	}
}

func WithMaxOpenConnections(maxOpenConnections int) ConnectionConfig {
	return func(config *Config) error {
		config.MaxOpenConnections = maxOpenConnections
		return nil
	}
}

func WithMaxIdleConnections(maxIdleConnections int) ConnectionConfig {
	return func(config *Config) error {
		config.MaxIdleConnections = maxIdleConnections
		return nil
	}
}

func WithConnMaxLifetime(connMaxLifetime time.Duration) ConnectionConfig {
	return func(config *Config) error {
		config.ConnMaxLifetime = connMaxLifetime
		return nil
	}
}

func WithConnMaxIdleTime(connMaxIdleTime time.Duration) ConnectionConfig {
	return func(config *Config) error {
		config.ConnMaxIdleTime = connMaxIdleTime
		return nil
	}
}

func WithParams(m map[string]string) ConnectionConfig {
	return func(config *Config) error {
		config.Params = m
		return nil
	}
}

func (cfg Config) FormatDSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
	)

	for key, value := range cfg.Params {
		dsn += fmt.Sprintf(" %s=%s", key, value)
	}

	return dsn
}

func MustConfig(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

func (cfg *Config) validate() error {
	if cfg.Host == "" {
		return fmt.Errorf("postgres host is required")
	}
	if cfg.User == "" {
		return fmt.Errorf("postgres user is required")
	}
	if cfg.Database == "" {
		return fmt.Errorf("postgres database is required")
	}
	if cfg.MaxOpenConnections < 0 {
		return fmt.Errorf("max open connections must be >= 0")
	}
	if cfg.MaxIdleConnections < 0 {
		return fmt.Errorf("max idle connections must be >= 0")
	}
	return nil
}
