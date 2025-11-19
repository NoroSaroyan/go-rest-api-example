package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		wantErr     bool
		validate    func(*Config) bool
		description string
	}{
		{
			name: "default values",
			env:  map[string]string{},
			validate: func(c *Config) bool {
				return c.App.Port == "8080" &&
					c.DB.Host == "localhost" &&
					c.DB.Port == 5432 &&
					c.Log.Level == "info"
			},
			description: "should load with default values when no env vars set",
		},
		{
			name: "custom values",
			env: map[string]string{
				"APP_PORT":  "3000",
				"DB_HOST":   "db.example.com",
				"DB_PORT":   "5433",
				"LOG_LEVEL": "debug",
			},
			validate: func(c *Config) bool {
				return c.App.Port == "3000" &&
					c.DB.Host == "db.example.com" &&
					c.DB.Port == 5433 &&
					c.Log.Level == "debug"
			},
			description: "should load custom values from environment",
		},
		{
			name: "invalid port",
			env: map[string]string{
				"APP_PORT": "99999",
			},
			wantErr:     true,
			description: "should fail validation with invalid port",
		},
		{
			name: "invalid db port",
			env: map[string]string{
				"DB_PORT": "invalid",
			},
			wantErr:     true,
			description: "should fail with invalid database port",
		},
		{
			name: "invalid log level",
			env: map[string]string{
				"LOG_LEVEL": "invalid",
			},
			wantErr:     true,
			description: "should fail validation with invalid log level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.env {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error but got nil - %s", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error = %v - %s", err, tt.description)
				return
			}

			if cfg == nil {
				t.Errorf("Load() returned nil config - %s", tt.description)
				return
			}

			if tt.validate != nil && !tt.validate(cfg) {
				t.Errorf("Load() validation failed - %s", tt.description)
			}
		})
	}
}

func TestDatabaseURL(t *testing.T) {
	cfg := &Config{
		DB: DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
		},
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	actual := cfg.DatabaseURL()

	if actual != expected {
		t.Errorf("DatabaseURL() = %v, want %v", actual, expected)
	}
}

func TestDatabaseURLEscaping(t *testing.T) {
	cfg := &Config{
		DB: DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "user@domain",
			Password: "pass word",
			Name:     "testdb",
		},
	}

	url := cfg.DatabaseURL()

	// Should contain escaped values
	if !contains(url, "user%40domain") {
		t.Error("DatabaseURL() should escape @ in username")
	}

	if !contains(url, "pass+word") {
		t.Error("DatabaseURL() should escape spaces in password")
	}
}

func TestConfigTimeouts(t *testing.T) {
	os.Setenv("APP_READ_TIMEOUT", "30s")
	os.Setenv("APP_WRITE_TIMEOUT", "45s")
	defer func() {
		os.Unsetenv("APP_READ_TIMEOUT")
		os.Unsetenv("APP_WRITE_TIMEOUT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.App.ReadTimeout != 30*time.Second {
		t.Errorf("ReadTimeout = %v, want %v", cfg.App.ReadTimeout, 30*time.Second)
	}

	if cfg.App.WriteTimeout != 45*time.Second {
		t.Errorf("WriteTimeout = %v, want %v", cfg.App.WriteTimeout, 45*time.Second)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
