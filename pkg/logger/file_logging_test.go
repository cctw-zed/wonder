package logger

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileLogging(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	t.Run("log to file only", func(t *testing.T) {
		config := LogConfig{
			Level:    "info",
			Format:   "json",
			Output:   "file",
			FilePath: logFile,
		}

		logger := NewLoggerWithConfig(config)
		ctx := context.Background()

		// Log some messages
		logger.Info(ctx, "test message", "key", "value")
		logger.Warn(ctx, "warning message", "level", "warn")
		logger.Error(ctx, "error message", "error", "test error")

		// Give time for writes to complete
		time.Sleep(100 * time.Millisecond)

		// Verify file exists and contains logs
		assert.FileExists(t, logFile)

		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		logContent := string(content)
		assert.Contains(t, logContent, "test message")
		assert.Contains(t, logContent, "warning message")
		assert.Contains(t, logContent, "error message")
		assert.Contains(t, logContent, `"key":"value"`)
		assert.Contains(t, logContent, `"level":"info"`)
		assert.Contains(t, logContent, `"level":"warning"`)
		assert.Contains(t, logContent, `"level":"error"`)
	})

	t.Run("log to both stdout and file", func(t *testing.T) {
		bothLogFile := filepath.Join(tempDir, "both.log")
		config := LogConfig{
			Level:    "debug",
			Format:   "json",
			Output:   "both",
			FilePath: bothLogFile,
		}

		logger := NewLoggerWithConfig(config)
		ctx := context.Background()

		// Log some messages
		logger.Debug(ctx, "debug message", "component", "test")
		logger.Info(ctx, "info message", "action", "both_logging")

		// Give time for writes to complete
		time.Sleep(100 * time.Millisecond)

		// Verify file exists and contains logs
		assert.FileExists(t, bothLogFile)

		content, err := os.ReadFile(bothLogFile)
		require.NoError(t, err)

		logContent := string(content)
		assert.Contains(t, logContent, "debug message")
		assert.Contains(t, logContent, "info message")
		assert.Contains(t, logContent, `"component":"test"`)
		assert.Contains(t, logContent, `"action":"both_logging"`)
	})

	t.Run("text format logging", func(t *testing.T) {
		textLogFile := filepath.Join(tempDir, "text.log")
		config := LogConfig{
			Level:    "info",
			Format:   "text",
			Output:   "file",
			FilePath: textLogFile,
		}

		logger := NewLoggerWithConfig(config)
		ctx := context.Background()

		// Log some messages
		logger.Info(ctx, "text format message", "format", "text")

		// Give time for writes to complete
		time.Sleep(100 * time.Millisecond)

		// Verify file exists and contains logs in text format
		assert.FileExists(t, textLogFile)

		content, err := os.ReadFile(textLogFile)
		require.NoError(t, err)

		logContent := string(content)
		assert.Contains(t, logContent, "text format message")
		assert.Contains(t, logContent, "format=text")
		assert.Contains(t, logContent, "level=info")
		// Text format should not contain JSON structure
		assert.NotContains(t, logContent, `"level":"info"`)
	})

	t.Run("directory creation", func(t *testing.T) {
		// Test with nested directories
		nestedLogFile := filepath.Join(tempDir, "nested", "deeper", "test.log")
		config := LogConfig{
			Level:    "info",
			Format:   "json",
			Output:   "file",
			FilePath: nestedLogFile,
		}

		logger := NewLoggerWithConfig(config)
		ctx := context.Background()

		// Log a message
		logger.Info(ctx, "nested directory test", "nested", true)

		// Give time for writes to complete
		time.Sleep(100 * time.Millisecond)

		// Verify directory and file were created
		assert.FileExists(t, nestedLogFile)

		content, err := os.ReadFile(nestedLogFile)
		require.NoError(t, err)

		assert.Contains(t, string(content), "nested directory test")
	})

	t.Run("invalid file path fallback", func(t *testing.T) {
		// Use an invalid path that cannot be created
		invalidPath := "/root/cannot/create/this.log"
		config := LogConfig{
			Level:    "info",
			Format:   "json",
			Output:   "file",
			FilePath: invalidPath,
		}

		// This should not panic and should fallback to stdout
		logger := NewLoggerWithConfig(config)
		ctx := context.Background()

		// Should not panic
		logger.Info(ctx, "should fallback to stdout", "fallback", true)

		// Should succeed without error
		assert.NotNil(t, logger)
	})
}

func TestLogLevels(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "levels.log")

	t.Run("log level filtering", func(t *testing.T) {
		config := LogConfig{
			Level:    "warn", // Only warn and error should be logged
			Format:   "json",
			Output:   "file",
			FilePath: logFile,
		}

		logger := NewLoggerWithConfig(config)
		ctx := context.Background()

		// Log messages at different levels
		logger.Debug(ctx, "debug message") // Should not appear
		logger.Info(ctx, "info message")   // Should not appear
		logger.Warn(ctx, "warn message")   // Should appear
		logger.Error(ctx, "error message") // Should appear

		// Give time for writes to complete
		time.Sleep(100 * time.Millisecond)

		content, err := os.ReadFile(logFile)
		require.NoError(t, err)

		logContent := string(content)

		// Debug and info should not be present
		assert.NotContains(t, logContent, "debug message")
		assert.NotContains(t, logContent, "info message")

		// Warn and error should be present
		assert.Contains(t, logContent, "warn message")
		assert.Contains(t, logContent, "error message")

		// Count the number of log lines (should be 2)
		lines := strings.Count(logContent, "\n")
		assert.Equal(t, 2, lines)
	})
}

func TestGlobalLoggerWithConfig(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "global.log")

	// Initialize global logger with file output
	InitializeWithConfig(LogConfig{
		Level:    "info",
		Format:   "json",
		Output:   "file",
		FilePath: logFile,
	})

	ctx := context.Background()

	// Use global logger functions
	LogInfo(ctx, "global info message", "global", true)
	LogWarn(ctx, "global warn message", "type", "warning")
	LogError(ctx, "global error message", "error", "test")

	// Give time for writes to complete
	time.Sleep(100 * time.Millisecond)

	// Verify file exists and contains logs
	assert.FileExists(t, logFile)

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	logContent := string(content)
	assert.Contains(t, logContent, "global info message")
	assert.Contains(t, logContent, "global warn message")
	assert.Contains(t, logContent, "global error message")
	assert.Contains(t, logContent, `"global":true`)
}
