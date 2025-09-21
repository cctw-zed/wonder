package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	require.NotNil(t, loader)
	require.NotNil(t, loader.viper)
}

func TestLoader_LoadConfig_DefaultConfig(t *testing.T) {
	loader := NewLoader()

	// Load config without any config file (should use defaults)
	config, err := loader.LoadConfig("/nonexistent/path")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Verify default values
	assert.Equal(t, "wonder", config.App.Name)
	assert.Equal(t, "1.0.0", config.App.Version)
	assert.Equal(t, "development", config.App.Environment)
	assert.True(t, config.App.Debug)
}

func TestLoader_LoadConfig_WithConfigFile(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	configContent := `
app:
  name: "test-app"
  version: "2.0.0"
  environment: "testing"
  debug: false

server:
  host: "0.0.0.0"
  port: 9090
  read_timeout: "60s"
  write_timeout: "60s"
  idle_timeout: "120s"
  enable_cors: false

database:
  host: "testdb"
  port: 5433
  username: "testuser"
  password: "testpass"
  database: "testdb"
  ssl_mode: "require"
  timezone: "Asia/Shanghai"
  max_open_conns: 50
  max_idle_conns: 20
  conn_max_lifetime: "2h"
  conn_max_idle_time: "1h"
  log_level: "debug"

log:
  level: "debug"
  format: "text"
  output: "file"
  enable_file: true
  file_path: "/tmp/test.log"

id:
  service_type: "order"
  instance_id: 5
  node_id: 10
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	loader := NewLoader()
	config, err := loader.LoadConfig(tempDir)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Verify loaded values
	assert.Equal(t, "test-app", config.App.Name)
	assert.Equal(t, "2.0.0", config.App.Version)
	assert.Equal(t, "testing", config.App.Environment)
	assert.False(t, config.App.Debug)

	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, 60*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 60*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.Server.IdleTimeout)
	assert.False(t, config.Server.EnableCORS)

	assert.Equal(t, "testdb", config.Database.Host)
	assert.Equal(t, 5433, config.Database.Port)
	assert.Equal(t, "testuser", config.Database.Username)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "testdb", config.Database.Database)
	assert.Equal(t, "require", config.Database.SSLMode)
	assert.Equal(t, "Asia/Shanghai", config.Database.Timezone)
	assert.Equal(t, 50, config.Database.MaxOpenConns)
	assert.Equal(t, 20, config.Database.MaxIdleConns)
	assert.Equal(t, 2*time.Hour, config.Database.ConnMaxLifetime)
	assert.Equal(t, 1*time.Hour, config.Database.ConnMaxIdleTime)
	assert.Equal(t, "debug", config.Database.LogLevel)

	assert.Equal(t, "debug", config.Log.Level)
	assert.Equal(t, "text", config.Log.Format)
	assert.Equal(t, "file", config.Log.Output)
	assert.True(t, config.Log.EnableFile)
	assert.Equal(t, "/tmp/test.log", config.Log.FilePath)

	assert.Equal(t, "order", config.ID.ServiceType)
	assert.Equal(t, int64(5), config.ID.InstanceID)
	assert.Equal(t, int64(10), config.ID.NodeID)
}

func TestLoader_LoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"APP_NAME":        "env-app",
		"APP_VERSION":     "3.0.0",
		"APP_ENV":         "production",
		"APP_DEBUG":       "false",
		"SERVER_HOST":     "prod.example.com",
		"SERVER_PORT":     "443",
		"DB_HOST":         "prod-db.example.com",
		"DB_PORT":         "5432",
		"DB_USERNAME":     "prod_user",
		"DB_PASSWORD":     "prod_password",
		"DB_DATABASE":     "prod_db",
		"LOG_LEVEL":       "error",
		"ID_SERVICE_TYPE": "payment",
		"ID_INSTANCE_ID":  "100",
		"ID_NODE_ID":      "200",
	}

	// Set environment variables
	for key, value := range envVars {
		t.Setenv(key, value)
	}

	loader := NewLoader()
	config, err := loader.LoadConfig("/nonexistent/path")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Verify environment variable overrides
	assert.Equal(t, "env-app", config.App.Name)
	assert.Equal(t, "3.0.0", config.App.Version)
	assert.Equal(t, "production", config.App.Environment)
	assert.False(t, config.App.Debug)

	assert.Equal(t, "prod.example.com", config.Server.Host)
	assert.Equal(t, 443, config.Server.Port)

	assert.Equal(t, "prod-db.example.com", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "prod_user", config.Database.Username)
	assert.Equal(t, "prod_password", config.Database.Password)
	assert.Equal(t, "prod_db", config.Database.Database)

	assert.Equal(t, "error", config.Log.Level)

	assert.Equal(t, "payment", config.ID.ServiceType)
	assert.Equal(t, int64(100), config.ID.InstanceID)
	assert.Equal(t, int64(200), config.ID.NodeID)
}

func TestLoader_LoadConfigForEnvironment(t *testing.T) {
	tempDir := t.TempDir()

	// Create environment-specific config file
	configContent := `
app:
  name: "prod-app"
  version: "1.0.0"
  environment: "production"
  debug: false

server:
  host: "0.0.0.0"
  port: 80
`

	configFile := filepath.Join(tempDir, "config.production.yaml")
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	loader := NewLoader()
	config, err := loader.LoadConfigForEnvironment("production", tempDir)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "prod-app", config.App.Name)
	assert.Equal(t, "production", config.App.Environment)
	assert.False(t, config.App.Debug)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 80, config.Server.Port)
}

func TestLoader_LoadConfig_InvalidConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create invalid YAML file
	invalidYAML := `
app:
  name: "test-app"
  invalid: [
server:
  port: "invalid_port"
`

	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	loader := NewLoader()
	_, err = loader.LoadConfig(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoader_LoadConfig_ValidationFailure(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create config with invalid values that will fail validation
	invalidConfig := `
app:
  name: ""
  version: "1.0.0"
  environment: "development"

server:
  host: "localhost"
  port: 0
`

	err := os.WriteFile(configFile, []byte(invalidConfig), 0644)
	require.NoError(t, err)

	loader := NewLoader()
	_, err = loader.LoadConfig(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config validation failed")
}

func TestLoader_WriteConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "output", "config.yaml")

	config := DefaultConfig()
	config.App.Name = "written-app"
	config.App.Version = "4.0.0"
	config.Server.Port = 9000

	loader := NewLoader()
	err := loader.WriteConfigFile(config, configFile)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configFile)
	require.NoError(t, err)

	// Load the written config and verify
	loadedConfig, err := loader.LoadConfig(filepath.Dir(configFile))
	require.NoError(t, err)

	assert.Equal(t, "written-app", loadedConfig.App.Name)
	assert.Equal(t, "4.0.0", loadedConfig.App.Version)
	assert.Equal(t, 9000, loadedConfig.Server.Port)
}

func TestLoad_GlobalFunction(t *testing.T) {
	// Test the global Load function
	config, err := Load("/nonexistent/path")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should get default config
	assert.Equal(t, "wonder", config.App.Name)
	assert.Equal(t, "development", config.App.Environment)
}

func TestLoadForEnvironment_GlobalFunction(t *testing.T) {
	// Test the global LoadForEnvironment function
	config, err := LoadForEnvironment("testing", "/nonexistent/path")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Environment should be overridden
	assert.Equal(t, "testing", config.App.Environment)
}

func TestWriteConfig_GlobalFunction(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "global-config.yaml")

	config := DefaultConfig()
	config.App.Name = "global-app"

	err := WriteConfig(config, configFile)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configFile)
	require.NoError(t, err)
}

func TestLoader_GetConfigFilePath(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	configContent := `
app:
  name: "path-test"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	loader := NewLoader()
	_, err = loader.LoadConfig(tempDir)
	require.NoError(t, err)

	usedPath := loader.GetConfigFilePath()
	assert.Equal(t, configFile, usedPath)
}

func TestLoader_EnvironmentVariableBinding(t *testing.T) {
	// Test that specific environment variables are properly bound
	testCases := []struct {
		envVar      string
		configPath  string
		setValue    string
		expectValue interface{}
	}{
		{"SERVICE_TYPE", "id.service_type", "auth", "auth"},
		{"INSTANCE_ID", "id.instance_id", "42", int64(42)},
		{"NODE_ID", "id.node_id", "99", int64(99)},
		{"WONDER_APP_NAME", "app.name", "wonder-env", "wonder-env"},
	}

	for _, tc := range testCases {
		t.Run(tc.envVar, func(t *testing.T) {
			t.Setenv(tc.envVar, tc.setValue)

			loader := NewLoader()
			config, err := loader.LoadConfig("/nonexistent/path")
			require.NoError(t, err)

			switch tc.configPath {
			case "id.service_type":
				assert.Equal(t, tc.expectValue, config.ID.ServiceType)
			case "id.instance_id":
				assert.Equal(t, tc.expectValue, config.ID.InstanceID)
			case "id.node_id":
				assert.Equal(t, tc.expectValue, config.ID.NodeID)
			case "app.name":
				assert.Equal(t, tc.expectValue, config.App.Name)
			}
		})
	}
}