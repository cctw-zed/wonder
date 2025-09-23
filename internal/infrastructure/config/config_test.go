package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	require.NotNil(t, config)

	// Test app configuration
	assert.Equal(t, "wonder", config.App.Name)
	assert.Equal(t, "1.0.0", config.App.Version)
	assert.Equal(t, "development", config.App.Environment)
	assert.True(t, config.App.Debug)

	// Test server configuration
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, 30*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.Server.IdleTimeout)
	assert.True(t, config.Server.EnableCORS)

	// Test database configuration
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "dev", config.Database.Username)
	assert.Equal(t, "dev", config.Database.Password)
	assert.Equal(t, "wonder_dev", config.Database.Database)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "UTC", config.Database.Timezone)
	assert.Equal(t, 25, config.Database.MaxOpenConns)
	assert.Equal(t, 10, config.Database.MaxIdleConns)
	assert.Equal(t, time.Hour, config.Database.ConnMaxLifetime)
	assert.Equal(t, 30*time.Minute, config.Database.ConnMaxIdleTime)
	assert.Equal(t, "info", config.Database.LogLevel)

	// Test log configuration
	assert.Equal(t, "info", config.Log.Level)
	assert.Equal(t, "json", config.Log.Format)
	assert.Equal(t, "stdout", config.Log.Output)
	assert.False(t, config.Log.EnableFile)
	assert.Equal(t, "logs/app.log", config.Log.FilePath)

	// Test ID configuration
	assert.Equal(t, "user", config.ID.ServiceType)
	assert.Equal(t, int64(0), config.ID.InstanceID)
	assert.Equal(t, int64(1), config.ID.NodeID)

	// Test external configuration
	require.NotNil(t, config.External.Redis)
	assert.Equal(t, "localhost", config.External.Redis.Host)
	assert.Equal(t, 6379, config.External.Redis.Port)
	assert.Equal(t, "", config.External.Redis.Password)
	assert.Equal(t, 0, config.External.Redis.Database)

	require.NotNil(t, config.External.Email)
	assert.Equal(t, "smtp", config.External.Email.Provider)
	assert.Equal(t, "smtp.gmail.com", config.External.Email.Host)
	assert.Equal(t, 587, config.External.Email.Port)
}

func TestConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := DefaultConfig()
		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid app config", func(t *testing.T) {
		config := DefaultConfig()
		config.App.Name = ""
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app config validation failed")
	})

	t.Run("invalid server config", func(t *testing.T) {
		config := DefaultConfig()
		config.Server.Port = 0
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "server config validation failed")
	})

	t.Run("invalid database config", func(t *testing.T) {
		config := DefaultConfig()
		config.Database.Host = ""
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database config validation failed")
	})

	t.Run("invalid log config", func(t *testing.T) {
		config := DefaultConfig()
		config.Log.Level = "invalid"
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "log config validation failed")
	})

	t.Run("invalid id config", func(t *testing.T) {
		config := DefaultConfig()
		config.ID.ServiceType = "invalid"
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "id config validation failed")
	})
}

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AppConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &AppConfig{
				Name:        "wonder",
				Version:     "1.0.0",
				Environment: "development",
				Debug:       true,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			config: &AppConfig{
				Name:        "",
				Version:     "1.0.0",
				Environment: "development",
				Debug:       true,
			},
			wantErr: true,
			errMsg:  "app name is required",
		},
		{
			name: "empty version",
			config: &AppConfig{
				Name:        "wonder",
				Version:     "",
				Environment: "development",
				Debug:       true,
			},
			wantErr: true,
			errMsg:  "app version is required",
		},
		{
			name: "empty environment",
			config: &AppConfig{
				Name:        "wonder",
				Version:     "1.0.0",
				Environment: "",
				Debug:       true,
			},
			wantErr: true,
			errMsg:  "app environment is required",
		},
		{
			name: "invalid environment",
			config: &AppConfig{
				Name:        "wonder",
				Version:     "1.0.0",
				Environment: "invalid",
				Debug:       true,
			},
			wantErr: true,
			errMsg:  "app environment must be one of: development, testing, production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ServerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &ServerConfig{
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				EnableCORS:   true,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			config: &ServerConfig{
				Host:         "",
				Port:         8080,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				EnableCORS:   true,
			},
			wantErr: true,
			errMsg:  "server host is required",
		},
		{
			name: "invalid port zero",
			config: &ServerConfig{
				Host:         "localhost",
				Port:         0,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				EnableCORS:   true,
			},
			wantErr: true,
			errMsg:  "server port must be between 1 and 65535",
		},
		{
			name: "invalid port too high",
			config: &ServerConfig{
				Host:         "localhost",
				Port:         65536,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				EnableCORS:   true,
			},
			wantErr: true,
			errMsg:  "server port must be between 1 and 65535",
		},
		{
			name: "invalid read timeout",
			config: &ServerConfig{
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  0,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  60 * time.Second,
				EnableCORS:   true,
			},
			wantErr: true,
			errMsg:  "server read_timeout must be positive",
		},
		{
			name: "invalid write timeout",
			config: &ServerConfig{
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 0,
				IdleTimeout:  60 * time.Second,
				EnableCORS:   true,
			},
			wantErr: true,
			errMsg:  "server write_timeout must be positive",
		},
		{
			name: "invalid idle timeout",
			config: &ServerConfig{
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
				IdleTimeout:  0,
				EnableCORS:   true,
			},
			wantErr: true,
			errMsg:  "server idle_timeout must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *LogConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &LogConfig{
				Level:       "info",
				Format:      "json",
				Output:      "stdout",
				EnableFile:  false,
				FilePath:    "logs/app.log",
				ServiceName: "test-service",
				MaxFileSize: 100,
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: &LogConfig{
				Level:      "invalid",
				Format:     "json",
				Output:     "stdout",
				EnableFile: false,
				FilePath:   "logs/app.log",
			},
			wantErr: true,
			errMsg:  "log level must be one of",
		},
		{
			name: "invalid format",
			config: &LogConfig{
				Level:      "info",
				Format:     "invalid",
				Output:     "stdout",
				EnableFile: false,
				FilePath:   "logs/app.log",
			},
			wantErr: true,
			errMsg:  "log format must be one of",
		},
		{
			name: "invalid output",
			config: &LogConfig{
				Level:      "info",
				Format:     "json",
				Output:     "invalid",
				EnableFile: false,
				FilePath:   "logs/app.log",
			},
			wantErr: true,
			errMsg:  "log output must be one of",
		},
		{
			name: "file enabled but no path",
			config: &LogConfig{
				Level:      "info",
				Format:     "json",
				Output:     "file",
				EnableFile: true,
				FilePath:   "",
			},
			wantErr: true,
			errMsg:  "log file_path is required when enable_file is true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIDConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *IDConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &IDConfig{
				ServiceType: "user",
				InstanceID:  0,
				NodeID:      1,
			},
			wantErr: false,
		},
		{
			name: "invalid service type",
			config: &IDConfig{
				ServiceType: "invalid",
				InstanceID:  0,
				NodeID:      1,
			},
			wantErr: true,
			errMsg:  "id service_type must be one of",
		},
		{
			name: "negative instance id",
			config: &IDConfig{
				ServiceType: "user",
				InstanceID:  -1,
				NodeID:      1,
			},
			wantErr: true,
			errMsg:  "id instance_id must be non-negative",
		},
		{
			name: "negative node id",
			config: &IDConfig{
				ServiceType: "user",
				InstanceID:  0,
				NodeID:      -1,
			},
			wantErr: true,
			errMsg:  "id node_id must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_EnvironmentHelpers(t *testing.T) {
	tests := []struct {
		name          string
		environment   string
		isDevelopment bool
		isProduction  bool
		isTesting     bool
	}{
		{
			name:          "development",
			environment:   "development",
			isDevelopment: true,
			isProduction:  false,
			isTesting:     false,
		},
		{
			name:          "production",
			environment:   "production",
			isDevelopment: false,
			isProduction:  true,
			isTesting:     false,
		},
		{
			name:          "testing",
			environment:   "testing",
			isDevelopment: false,
			isProduction:  false,
			isTesting:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.App.Environment = tt.environment

			assert.Equal(t, tt.environment, config.GetEnvironment())
			assert.Equal(t, tt.isDevelopment, config.IsDevelopment())
			assert.Equal(t, tt.isProduction, config.IsProduction())
			assert.Equal(t, tt.isTesting, config.IsTesting())
		})
	}
}
