package logger

import (
	"sync"
)

// LoggerFactory implements the Factory interface
type LoggerFactory struct {
	config *Config
	logger Logger
	mu     sync.RWMutex
}

// NewFactory creates a new logger factory
func NewFactory(config *Config) Factory {
	return &LoggerFactory{
		config: config,
	}
}

// NewLogger creates a new base logger instance
func (f *LoggerFactory) NewLogger() Logger {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.logger == nil {
		f.logger = NewLogrusLogger(f.config)
	}
	return f.logger
}

// NewDomainLogger creates a domain layer specific logger
func (f *LoggerFactory) NewDomainLogger() Logger {
	return f.NewLogger().WithLayer(DomainLayer)
}

// NewApplicationLogger creates an application layer specific logger
func (f *LoggerFactory) NewApplicationLogger() Logger {
	return f.NewLogger().WithLayer(ApplicationLayer)
}

// NewInfrastructureLogger creates an infrastructure layer specific logger
func (f *LoggerFactory) NewInfrastructureLogger() Logger {
	return f.NewLogger().WithLayer(InfrastructureLayer)
}

// NewInterfaceLogger creates an interface layer specific logger
func (f *LoggerFactory) NewInterfaceLogger() Logger {
	return f.NewLogger().WithLayer(InterfaceLayer)
}

// NewComponentLogger creates a component-specific logger
func (f *LoggerFactory) NewComponentLogger(component string) Logger {
	return f.NewLogger().WithComponent(component)
}

// Ensure LoggerFactory implements Factory interface
var _ Factory = (*LoggerFactory)(nil)

// Global factory instance for convenience
var globalFactory Factory

// SetGlobalFactory sets the global logger factory
func SetGlobalFactory(factory Factory) {
	globalFactory = factory
}

// GetGlobalFactory returns the global logger factory
func GetGlobalFactory() Factory {
	return globalFactory
}

// Convenience functions for creating loggers
func NewLogger() Logger {
	if globalFactory == nil {
		panic("global logger factory not initialized")
	}
	return globalFactory.NewLogger()
}

func NewComponentLogger(component string) Logger {
	if globalFactory == nil {
		panic("global logger factory not initialized")
	}
	return globalFactory.NewComponentLogger(component)
}