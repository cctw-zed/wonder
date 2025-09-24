package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cctw-zed/wonder/pkg/logger"
)

func main() {
	// 确保日志目录存在
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	fmt.Println("=== File Logging Demo ===\n")

	// Demo 1: JSON格式输出到文件
	fmt.Println("1. JSON格式日志输出到文件...")
	jsonLogFile := filepath.Join(logDir, "demo_json.log")
	logger.InitializeWithConfig(logger.LogConfig{
		Level:    "debug",
		Format:   "json",
		Output:   "file",
		FilePath: jsonLogFile,
	})

	ctx := context.Background()
	appLogger := logger.Get().WithLayer("demo").WithComponent("json_logger")

	appLogger.Debug(ctx, "这是DEBUG级别的日志", "step", 1, "format", "json")
	appLogger.Info(ctx, "应用启动", "version", "1.0.0", "environment", "demo")
	appLogger.Warn(ctx, "检测到潜在问题", "issue", "memory_usage", "usage_percent", 85)
	appLogger.Error(ctx, "模拟错误", "error", "connection timeout", "retry_count", 3)

	fmt.Printf("   JSON日志已写入: %s\n", jsonLogFile)

	// Demo 2: 文本格式输出到不同文件
	fmt.Println("\n2. 文本格式日志输出到文件...")
	textLogFile := filepath.Join(logDir, "demo_text.log")
	logger.InitializeWithConfig(logger.LogConfig{
		Level:    "info",
		Format:   "text",
		Output:   "file",
		FilePath: textLogFile,
	})

	textLogger := logger.Get().WithLayer("demo").WithComponent("text_logger")
	textLogger.Info(ctx, "用户注册", "user_id", "user_123", "email", "demo@example.com")
	textLogger.Warn(ctx, "密码强度不够", "user_id", "user_123", "strength", "weak")
	textLogger.Info(ctx, "邮件验证发送", "user_id", "user_123", "email", "demo@example.com")

	fmt.Printf("   文本日志已写入: %s\n", textLogFile)

	// Demo 3: 同时输出到控制台和文件
	fmt.Println("\n3. 同时输出到控制台和文件...")
	bothLogFile := filepath.Join(logDir, "demo_both.log")
	logger.InitializeWithConfig(logger.LogConfig{
		Level:    "info",
		Format:   "json",
		Output:   "both",
		FilePath: bothLogFile,
	})

	bothLogger := logger.Get().WithLayer("demo").WithComponent("both_logger")
	bothLogger.Info(ctx, "这条日志同时显示在控制台和文件中", "demo", "both_output", "timestamp", time.Now().Format(time.RFC3339))

	fmt.Printf("   同时输出日志已写入: %s\n", bothLogFile)
	fmt.Println("   (并且也显示在上方控制台中)")

	// Demo 4: 级别过滤演示
	fmt.Println("\n4. 日志级别过滤演示 (只记录WARN和ERROR)...")
	filterLogFile := filepath.Join(logDir, "demo_filtered.log")
	logger.InitializeWithConfig(logger.LogConfig{
		Level:    "warn", // 只记录warn和error
		Format:   "json",
		Output:   "file",
		FilePath: filterLogFile,
	})

	filterLogger := logger.Get().WithLayer("demo").WithComponent("filter_logger")
	filterLogger.Debug(ctx, "这条DEBUG不会被记录", "filtered", true)
	filterLogger.Info(ctx, "这条INFO也不会被记录", "filtered", true)
	filterLogger.Warn(ctx, "这条WARN会被记录", "recorded", true)
	filterLogger.Error(ctx, "这条ERROR也会被记录", "recorded", true, "severity", "high")

	fmt.Printf("   过滤日志已写入: %s (只包含WARN和ERROR)\n", filterLogFile)

	// 等待文件写入完成
	time.Sleep(100 * time.Millisecond)

	// 显示文件内容
	fmt.Println("\n=== 日志文件内容预览 ===")
	showLogFileContent("JSON格式日志", jsonLogFile)
	showLogFileContent("文本格式日志", textLogFile)
	showLogFileContent("过滤后的日志", filterLogFile)

	fmt.Println("\n=== Demo 完成 ===")
	fmt.Printf("所有日志文件保存在: %s\n", logDir)
	fmt.Println("你可以查看这些文件来了解不同的日志配置效果。")
}

func showLogFileContent(title, filePath string) {
	fmt.Printf("\n--- %s (%s) ---\n", title, filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	if len(content) == 0 {
		fmt.Println("(文件为空)")
		return
	}

	// 显示前500字符或全部内容（如果较短）
	contentStr := string(content)
	if len(contentStr) > 500 {
		fmt.Printf("%s...\n(内容被截断，完整内容请查看文件)\n", contentStr[:500])
	} else {
		fmt.Print(contentStr)
	}
}