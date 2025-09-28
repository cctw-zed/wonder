package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Loader handles configuration loading from multiple sources
type Loader struct {
	viper *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	v := viper.New()
	return &Loader{viper: v}
}

// LoadConfig loads configuration from files and environment variables
func (l *Loader) LoadConfig(configPaths ...string) (*Config, error) {
	// Set default configuration paths if none provided
	if len(configPaths) == 0 {
		configPaths = []string{
			".",
			"./configs",
			"./config",
			"/etc/wonder",
		}
	}

	// Configure viper
	l.setupViper(configPaths)

	// Apply environment variable overrides BEFORE loading config
	l.bindEnvironmentVariables()

	// Attempt to read configuration file
	if err := l.viper.ReadInConfig(); err != nil {
		// If config file is not found, that's OK - we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Start with default configuration and unmarshal from viper to override values
	config := DefaultConfig()
	if err := l.viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// LoadConfigForEnvironment loads configuration for specific environment
func (l *Loader) LoadConfigForEnvironment(env string, configPaths ...string) (*Config, error) {
	// Set default configuration paths if none provided
	if len(configPaths) == 0 {
		configPaths = []string{
			".",
			"./configs",
			"./config",
			"/etc/wonder",
		}
	}

	// Configure viper
	l.setupViperForEnvironment(env, configPaths)

	// Apply environment variable overrides BEFORE loading config
	l.bindEnvironmentVariables()

	// Attempt to read configuration file
	if err := l.viper.ReadInConfig(); err != nil {
		// If config file is not found, that's OK - we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Start with default configuration and unmarshal from viper to override values
	config := DefaultConfig()
	if err := l.viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override environment if specified
	if env != "" {
		config.App.Environment = env
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setupViper configures viper instance
func (l *Loader) setupViper(configPaths []string) {
	// Set config name (without extension)
	l.viper.SetConfigName("config")

	// Set config type
	l.viper.SetConfigType("yaml")

	// Add config paths
	for _, path := range configPaths {
		l.viper.AddConfigPath(path)
	}

	// Enable environment variable support
	l.viper.AutomaticEnv()
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	l.viper.SetEnvPrefix("WONDER")

	// Set defaults
	l.setDefaults()
}

// setupViperForEnvironment configures viper instance for specific environment
func (l *Loader) setupViperForEnvironment(env string, configPaths []string) {
	// Set environment-specific configuration name first (higher priority)
	l.viper.SetConfigName(fmt.Sprintf("config.%s", env))

	// Set config type
	l.viper.SetConfigType("yaml")

	// Add config paths
	for _, path := range configPaths {
		l.viper.AddConfigPath(path)
	}

	// Enable environment variable support
	l.viper.AutomaticEnv()
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	l.viper.SetEnvPrefix("WONDER")

	// Set defaults
	l.setDefaults()
}

// setDefaults sets default values in viper
func (l *Loader) setDefaults() {
	defaults := DefaultConfig()

	// App defaults
	l.viper.SetDefault("app.name", defaults.App.Name)
	l.viper.SetDefault("app.version", defaults.App.Version)
	l.viper.SetDefault("app.environment", defaults.App.Environment)
	l.viper.SetDefault("app.debug", defaults.App.Debug)

	// Server defaults
	l.viper.SetDefault("server.host", defaults.Server.Host)
	l.viper.SetDefault("server.port", defaults.Server.Port)
	l.viper.SetDefault("server.read_timeout", defaults.Server.ReadTimeout)
	l.viper.SetDefault("server.write_timeout", defaults.Server.WriteTimeout)
	l.viper.SetDefault("server.idle_timeout", defaults.Server.IdleTimeout)
	l.viper.SetDefault("server.enable_cors", defaults.Server.EnableCORS)

	// Database defaults
	l.viper.SetDefault("database.host", defaults.Database.Host)
	l.viper.SetDefault("database.port", defaults.Database.Port)
	l.viper.SetDefault("database.username", defaults.Database.Username)
	l.viper.SetDefault("database.password", defaults.Database.Password)
	l.viper.SetDefault("database.database", defaults.Database.Database)
	l.viper.SetDefault("database.ssl_mode", defaults.Database.SSLMode)
	l.viper.SetDefault("database.timezone", defaults.Database.Timezone)
	l.viper.SetDefault("database.max_open_conns", defaults.Database.MaxOpenConns)
	l.viper.SetDefault("database.max_idle_conns", defaults.Database.MaxIdleConns)
	l.viper.SetDefault("database.conn_max_lifetime", defaults.Database.ConnMaxLifetime)
	l.viper.SetDefault("database.conn_max_idle_time", defaults.Database.ConnMaxIdleTime)
	l.viper.SetDefault("database.log_level", defaults.Database.LogLevel)

	// Log defaults
	l.viper.SetDefault("log.level", defaults.Log.Level)
	l.viper.SetDefault("log.format", defaults.Log.Format)
	l.viper.SetDefault("log.output", defaults.Log.Output)
	l.viper.SetDefault("log.enable_file", defaults.Log.EnableFile)
	l.viper.SetDefault("log.file_path", defaults.Log.FilePath)

	// ID defaults
	l.viper.SetDefault("id.service_type", defaults.ID.ServiceType)
	l.viper.SetDefault("id.instance_id", defaults.ID.InstanceID)
	l.viper.SetDefault("id.node_id", defaults.ID.NodeID)

	// External defaults
	if defaults.External.Redis != nil {
		l.viper.SetDefault("external.redis.host", defaults.External.Redis.Host)
		l.viper.SetDefault("external.redis.port", defaults.External.Redis.Port)
		l.viper.SetDefault("external.redis.password", defaults.External.Redis.Password)
		l.viper.SetDefault("external.redis.database", defaults.External.Redis.Database)
	}

	if defaults.External.Email != nil {
		l.viper.SetDefault("external.email.provider", defaults.External.Email.Provider)
		l.viper.SetDefault("external.email.host", defaults.External.Email.Host)
		l.viper.SetDefault("external.email.port", defaults.External.Email.Port)
		l.viper.SetDefault("external.email.username", defaults.External.Email.Username)
		l.viper.SetDefault("external.email.password", defaults.External.Email.Password)
	}
}

// bindEnvironmentVariables binds environment variables to configuration keys
func (l *Loader) bindEnvironmentVariables() {
	// App configuration
	l.viper.BindEnv("app.name", "APP_NAME")
	l.viper.BindEnv("app.version", "APP_VERSION")
	l.viper.BindEnv("app.environment", "APP_ENV")
	l.viper.BindEnv("app.debug", "APP_DEBUG")

	// Server configuration
	l.viper.BindEnv("server.host", "SERVER_HOST")
	l.viper.BindEnv("server.port", "SERVER_PORT")
	l.viper.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	l.viper.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	l.viper.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")
	l.viper.BindEnv("server.enable_cors", "SERVER_ENABLE_CORS")

	// Database configuration
	l.viper.BindEnv("database.host", "DB_HOST")
	l.viper.BindEnv("database.port", "DB_PORT")
	l.viper.BindEnv("database.username", "DB_USERNAME")
	l.viper.BindEnv("database.password", "DB_PASSWORD")
	l.viper.BindEnv("database.database", "DB_DATABASE")
	l.viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	l.viper.BindEnv("database.timezone", "DB_TIMEZONE")
	l.viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	l.viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	l.viper.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	l.viper.BindEnv("database.conn_max_idle_time", "DB_CONN_MAX_IDLE_TIME")
	l.viper.BindEnv("database.log_level", "DB_LOG_LEVEL")

	// Log configuration
	l.viper.BindEnv("log.level", "LOG_LEVEL")
	l.viper.BindEnv("log.format", "LOG_FORMAT")
	l.viper.BindEnv("log.output", "LOG_OUTPUT")
	l.viper.BindEnv("log.enable_file", "LOG_ENABLE_FILE")
	l.viper.BindEnv("log.file_path", "LOG_FILE_PATH")

	// ID configuration
	l.viper.BindEnv("id.service_type", "ID_SERVICE_TYPE", "SERVICE_TYPE")
	l.viper.BindEnv("id.instance_id", "ID_INSTANCE_ID", "INSTANCE_ID")
	l.viper.BindEnv("id.node_id", "ID_NODE_ID", "NODE_ID")

	// Redis configuration
	l.viper.BindEnv("external.redis.host", "REDIS_HOST")
	l.viper.BindEnv("external.redis.port", "REDIS_PORT")
	l.viper.BindEnv("external.redis.password", "REDIS_PASSWORD")
	l.viper.BindEnv("external.redis.database", "REDIS_DATABASE")

	// Email configuration
	l.viper.BindEnv("external.email.provider", "EMAIL_PROVIDER")
	l.viper.BindEnv("external.email.host", "EMAIL_HOST")
	l.viper.BindEnv("external.email.port", "EMAIL_PORT")
	l.viper.BindEnv("external.email.username", "EMAIL_USERNAME")
	l.viper.BindEnv("external.email.password", "EMAIL_PASSWORD")
}

// GetConfigFilePath returns the path of the loaded configuration file
func (l *Loader) GetConfigFilePath() string {
	return l.viper.ConfigFileUsed()
}

// WriteConfigFile writes the current configuration to a file
func (l *Loader) WriteConfigFile(config *Config, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create a new viper instance for writing
	writeViper := viper.New()
	writeViper.SetConfigFile(filePath)
	writeViper.SetConfigType("yaml")

	// Set all configuration values
	l.setConfigValues(writeViper, config)

	// Write configuration file
	if err := writeViper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// setConfigValues sets configuration values in viper for writing
func (l *Loader) setConfigValues(v *viper.Viper, config *Config) {
	// App configuration
	v.Set("app.name", config.App.Name)
	v.Set("app.version", config.App.Version)
	v.Set("app.environment", config.App.Environment)
	v.Set("app.debug", config.App.Debug)

	// Server configuration
	v.Set("server.host", config.Server.Host)
	v.Set("server.port", config.Server.Port)
	v.Set("server.read_timeout", config.Server.ReadTimeout)
	v.Set("server.write_timeout", config.Server.WriteTimeout)
	v.Set("server.idle_timeout", config.Server.IdleTimeout)
	v.Set("server.enable_cors", config.Server.EnableCORS)

	// Database configuration
	v.Set("database.host", config.Database.Host)
	v.Set("database.port", config.Database.Port)
	v.Set("database.username", config.Database.Username)
	v.Set("database.password", config.Database.Password)
	v.Set("database.database", config.Database.Database)
	v.Set("database.ssl_mode", config.Database.SSLMode)
	v.Set("database.timezone", config.Database.Timezone)
	v.Set("database.max_open_conns", config.Database.MaxOpenConns)
	v.Set("database.max_idle_conns", config.Database.MaxIdleConns)
	v.Set("database.conn_max_lifetime", config.Database.ConnMaxLifetime)
	v.Set("database.conn_max_idle_time", config.Database.ConnMaxIdleTime)
	v.Set("database.log_level", config.Database.LogLevel)

	// Log configuration
	v.Set("log.level", config.Log.Level)
	v.Set("log.format", config.Log.Format)
	v.Set("log.output", config.Log.Output)
	v.Set("log.enable_file", config.Log.EnableFile)
	v.Set("log.file_path", config.Log.FilePath)

	// ID configuration
	v.Set("id.service_type", config.ID.ServiceType)
	v.Set("id.instance_id", config.ID.InstanceID)
	v.Set("id.node_id", config.ID.NodeID)

	// External services configuration
	if config.External.Redis != nil {
		v.Set("external.redis.host", config.External.Redis.Host)
		v.Set("external.redis.port", config.External.Redis.Port)
		v.Set("external.redis.password", config.External.Redis.Password)
		v.Set("external.redis.database", config.External.Redis.Database)
	}

	if config.External.Email != nil {
		v.Set("external.email.provider", config.External.Email.Provider)
		v.Set("external.email.host", config.External.Email.Host)
		v.Set("external.email.port", config.External.Email.Port)
		v.Set("external.email.username", config.External.Email.Username)
		v.Set("external.email.password", config.External.Email.Password)
	}
}

// Global configuration loader instance
var globalLoader *Loader

// Load loads configuration using the global loader
func Load(configPaths ...string) (*Config, error) {
	if globalLoader == nil {
		globalLoader = NewLoader()
	}
	return globalLoader.LoadConfig(configPaths...)
}

// LoadForEnvironment loads configuration for specific environment using global loader
func LoadForEnvironment(env string, configPaths ...string) (*Config, error) {
	if globalLoader == nil {
		globalLoader = NewLoader()
	}
	return globalLoader.LoadConfigForEnvironment(env, configPaths...)
}

// WriteConfig writes configuration to file using global loader
func WriteConfig(config *Config, filePath string) error {
	if globalLoader == nil {
		globalLoader = NewLoader()
	}
	return globalLoader.WriteConfigFile(config, filePath)
}
