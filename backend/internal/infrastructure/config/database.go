package config

import (
	"fmt"
	"time"
)

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string        `yaml:"host" mapstructure:"host" env:"DB_HOST"`
	Port            int           `yaml:"port" mapstructure:"port" env:"DB_PORT"`
	Username        string        `yaml:"username" mapstructure:"username" env:"DB_USERNAME"`
	Password        string        `yaml:"password" mapstructure:"password" env:"DB_PASSWORD"`
	Database        string        `yaml:"database" mapstructure:"database" env:"DB_DATABASE"`
	SSLMode         string        `yaml:"ssl_mode" mapstructure:"ssl_mode" env:"DB_SSL_MODE"`
	Timezone        string        `yaml:"timezone" mapstructure:"timezone" env:"DB_TIMEZONE"`
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME"`
	LogLevel        string        `yaml:"log_level" mapstructure:"log_level" env:"DB_LOG_LEVEL"`
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		Username:        "dev",
		Password:        "dev",
		Database:        "wonder_dev",
		SSLMode:         "disable",
		Timezone:        "UTC",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 30,
		LogLevel:        "info",
	}
}

// DSN builds PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode, c.Timezone)
}

// Validate validates database configuration
func (c *DatabaseConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if c.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("max_open_conns must be positive")
	}
	if c.MaxIdleConns <= 0 {
		return fmt.Errorf("max_idle_conns must be positive")
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("max_idle_conns cannot be greater than max_open_conns")
	}
	return nil
}
