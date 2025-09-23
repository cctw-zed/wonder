package logger

import (
	appConfig "github.com/cctw-zed/wonder/internal/infrastructure/config"
)

// NewConfigFromAppConfig creates a logger config from the application config
func NewConfigFromAppConfig(appCfg *appConfig.Config) *Config {
	logCfg := appCfg.Log

	return &Config{
		Level:         Level(logCfg.Level),
		Format:        logCfg.Format,
		Output:        logCfg.Output,
		EnableFile:    logCfg.EnableFile,
		FilePath:      logCfg.FilePath,
		ServiceName:   logCfg.ServiceName,
		EnableTracing: logCfg.EnableTracing,
		MaxFileSize:   logCfg.MaxFileSize,
		MaxBackups:    logCfg.MaxBackups,
		MaxAge:        logCfg.MaxAge,
		Compress:      logCfg.Compress,
	}
}

// InitializeGlobalLogger initializes the global logger factory from app config
func InitializeGlobalLogger(appCfg *appConfig.Config) Factory {
	loggerConfig := NewConfigFromAppConfig(appCfg)
	factory := NewFactory(loggerConfig)
	SetGlobalFactory(factory)
	return factory
}