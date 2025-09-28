package config

import (
	"fmt"
	"time"
)

// Config represents the application configuration structure organized by DDD layers
type Config struct {
	// Application configuration
	App *AppConfig `yaml:"app" mapstructure:"app"`

	// Infrastructure layer configurations
	Database *DatabaseConfig `yaml:"database" mapstructure:"database"`
	Server   *ServerConfig   `yaml:"server" mapstructure:"server"`
	Log      *LogConfig      `yaml:"log" mapstructure:"log"`
	JWT      *JWTConfig      `yaml:"jwt" mapstructure:"jwt"`

	// Domain layer configurations
	ID *IDConfig `yaml:"id" mapstructure:"id"`

	// External services configurations
	External *ExternalConfig `yaml:"external" mapstructure:"external"`
}

// AppConfig represents application-level configuration
type AppConfig struct {
	Name        string `yaml:"name" mapstructure:"name" env:"APP_NAME"`
	Version     string `yaml:"version" mapstructure:"version" env:"APP_VERSION"`
	Environment string `yaml:"environment" mapstructure:"environment" env:"APP_ENV"`
	Debug       bool   `yaml:"debug" mapstructure:"debug" env:"APP_DEBUG"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host         string        `yaml:"host" mapstructure:"host" env:"SERVER_HOST"`
	Port         int           `yaml:"port" mapstructure:"port" env:"SERVER_PORT"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout" env:"SERVER_IDLE_TIMEOUT"`
	EnableCORS   bool          `yaml:"enable_cors" mapstructure:"enable_cors" env:"SERVER_ENABLE_CORS"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level         string `yaml:"level" mapstructure:"level" env:"LOG_LEVEL"`
	Format        string `yaml:"format" mapstructure:"format" env:"LOG_FORMAT"`
	Output        string `yaml:"output" mapstructure:"output" env:"LOG_OUTPUT"`
	EnableFile    bool   `yaml:"enable_file" mapstructure:"enable_file" env:"LOG_ENABLE_FILE"`
	FilePath      string `yaml:"file_path" mapstructure:"file_path" env:"LOG_FILE_PATH"`
	ServiceName   string `yaml:"service_name" mapstructure:"service_name" env:"LOG_SERVICE_NAME"`
	EnableTracing bool   `yaml:"enable_tracing" mapstructure:"enable_tracing" env:"LOG_ENABLE_TRACING"`
	MaxFileSize   int    `yaml:"max_file_size" mapstructure:"max_file_size" env:"LOG_MAX_FILE_SIZE"`
	MaxBackups    int    `yaml:"max_backups" mapstructure:"max_backups" env:"LOG_MAX_BACKUPS"`
	MaxAge        int    `yaml:"max_age" mapstructure:"max_age" env:"LOG_MAX_AGE"`
	Compress      bool   `yaml:"compress" mapstructure:"compress" env:"LOG_COMPRESS"`
}

// IDConfig represents ID generation configuration
type IDConfig struct {
	ServiceType string `yaml:"service_type" mapstructure:"service_type" env:"ID_SERVICE_TYPE"`
	InstanceID  int64  `yaml:"instance_id" mapstructure:"instance_id" env:"ID_INSTANCE_ID"`
	NodeID      int64  `yaml:"node_id" mapstructure:"node_id" env:"ID_NODE_ID"`
}

// ExternalConfig represents external services configuration
type ExternalConfig struct {
	Redis *RedisConfig `yaml:"redis" mapstructure:"redis"`
	Email *EmailConfig `yaml:"email" mapstructure:"email"`
}

// RedisConfig represents Redis configuration (future use)
type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"host" env:"REDIS_HOST"`
	Port     int    `yaml:"port" mapstructure:"port" env:"REDIS_PORT"`
	Password string `yaml:"password" mapstructure:"password" env:"REDIS_PASSWORD"`
	Database int    `yaml:"database" mapstructure:"database" env:"REDIS_DATABASE"`
}

// EmailConfig represents email service configuration (future use)
type EmailConfig struct {
	Provider string `yaml:"provider" mapstructure:"provider" env:"EMAIL_PROVIDER"`
	Host     string `yaml:"host" mapstructure:"host" env:"EMAIL_HOST"`
	Port     int    `yaml:"port" mapstructure:"port" env:"EMAIL_PORT"`
	Username string `yaml:"username" mapstructure:"username" env:"EMAIL_USERNAME"`
	Password string `yaml:"password" mapstructure:"password" env:"EMAIL_PASSWORD"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	SigningKey string        `yaml:"signing_key" mapstructure:"signing_key" env:"JWT_SIGNING_KEY"`
	Expiry     time.Duration `yaml:"expiry" mapstructure:"expiry" env:"JWT_EXPIRY"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		App: &AppConfig{
			Name:        "wonder",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Server: &ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			EnableCORS:   true,
		},
		Database: DefaultDatabaseConfig(),
		Log: &LogConfig{
			Level:         "info",
			Format:        "json",
			Output:        "stdout",
			EnableFile:    false,
			FilePath:      "logs/app.log",
			ServiceName:   "wonder",
			EnableTracing: true,
			MaxFileSize:   100, // MB
			MaxBackups:    3,
			MaxAge:        28, // days
			Compress:      true,
		},
		JWT: &JWTConfig{
			SigningKey: "your-secret-signing-key-change-this-in-production",
			Expiry:     24 * time.Hour,
		},
		ID: &IDConfig{
			ServiceType: "user",
			InstanceID:  0,
			NodeID:      1,
		},
		External: &ExternalConfig{
			Redis: &RedisConfig{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				Database: 0,
			},
			Email: &EmailConfig{
				Provider: "smtp",
				Host:     "smtp.gmail.com",
				Port:     587,
				Username: "",
				Password: "",
			},
		},
	}
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return fmt.Errorf("app config validation failed: %w", err)
	}

	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log config validation failed: %w", err)
	}

	if err := c.ID.Validate(); err != nil {
		return fmt.Errorf("id config validation failed: %w", err)
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config validation failed: %w", err)
	}

	return nil
}

// Validate validates app configuration
func (c *AppConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("app name is required")
	}
	if c.Version == "" {
		return fmt.Errorf("app version is required")
	}
	if c.Environment == "" {
		return fmt.Errorf("app environment is required")
	}
	if c.Environment != "development" && c.Environment != "testing" && c.Environment != "production" {
		return fmt.Errorf("app environment must be one of: development, testing, production")
	}
	return nil
}

// Validate validates server configuration
func (c *ServerConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if c.ReadTimeout <= 0 {
		return fmt.Errorf("server read_timeout must be positive")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("server write_timeout must be positive")
	}
	if c.IdleTimeout <= 0 {
		return fmt.Errorf("server idle_timeout must be positive")
	}
	return nil
}

// Validate validates log configuration
func (c *LogConfig) Validate() error {
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	valid := false
	for _, level := range validLevels {
		if c.Level == level {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("log level must be one of: %v", validLevels)
	}

	validFormats := []string{"json", "text", "console"}
	valid = false
	for _, format := range validFormats {
		if c.Format == format {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("log format must be one of: %v", validFormats)
	}

	validOutputs := []string{"stdout", "stderr", "file", "both"}
	valid = false
	for _, output := range validOutputs {
		if c.Output == output {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("log output must be one of: %v", validOutputs)
	}

	if c.EnableFile && c.FilePath == "" {
		return fmt.Errorf("log file_path is required when enable_file is true")
	}

	if c.ServiceName == "" {
		return fmt.Errorf("log service_name is required")
	}

	if c.MaxFileSize <= 0 {
		return fmt.Errorf("log max_file_size must be positive")
	}

	if c.MaxBackups < 0 {
		return fmt.Errorf("log max_backups must be non-negative")
	}

	if c.MaxAge < 0 {
		return fmt.Errorf("log max_age must be non-negative")
	}

	return nil
}

// Validate validates ID configuration
func (c *IDConfig) Validate() error {
	validServiceTypes := []string{"user", "order", "payment", "auth"}
	valid := false
	for _, serviceType := range validServiceTypes {
		if c.ServiceType == serviceType {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("id service_type must be one of: %v", validServiceTypes)
	}

	if c.InstanceID < 0 {
		return fmt.Errorf("id instance_id must be non-negative")
	}

	if c.NodeID < 0 {
		return fmt.Errorf("id node_id must be non-negative")
	}

	return nil
}

// Validate validates JWT configuration
func (c *JWTConfig) Validate() error {
	if c.SigningKey == "" {
		return fmt.Errorf("jwt signing_key is required")
	}
	if len(c.SigningKey) < 32 {
		return fmt.Errorf("jwt signing_key must be at least 32 characters long")
	}
	if c.Expiry <= 0 {
		return fmt.Errorf("jwt expiry must be positive")
	}
	return nil
}

// GetEnvironment returns the current environment
func (c *Config) GetEnvironment() string {
	return c.App.Environment
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsTesting returns true if running in testing environment
func (c *Config) IsTesting() bool {
	return c.App.Environment == "testing"
}
