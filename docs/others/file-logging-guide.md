# File Logging 配置指南

## 概览

Wonder应用支持灵活的日志输出配置，包括输出到控制台、文件或同时输出到两者。本指南将详细说明如何配置和使用文件日志功能。

## 快速开始

### 1. 通过配置文件设置

在配置文件中（如 `configs/config.yaml`）添加以下日志配置：

```yaml
log:
  # 日志级别: debug, info, warn, error
  level: "info"

  # 日志格式: json, text
  format: "json"

  # 日志输出: stdout, file, both
  output: "file"

  # 日志文件路径（当output为file或both时必需）
  file_path: "./logs/wonder.log"

  # 启用文件日志（可选）
  enable_file: true

  # 服务名称（用于日志上下文）
  service_name: "wonder-api"
```

### 2. 通过环境变量设置

也可以通过环境变量来配置日志：

```bash
export LOG_LEVEL=info
export LOG_FORMAT=json
export LOG_OUTPUT=file
export LOG_FILE_PATH=./logs/wonder.log
export LOG_ENABLE_FILE=true
```

### 3. 程序中直接配置

```go
package main

import (
    "context"
    "github.com/cctw-zed/wonder/pkg/logger"
)

func main() {
    // 配置文件日志
    logger.InitializeWithConfig(logger.LogConfig{
        Level:    "info",
        Format:   "json",
        Output:   "file",
        FilePath: "./logs/app.log",
    })

    // 使用全局logger
    ctx := context.Background()
    logger.LogInfo(ctx, "应用程序启动", "version", "1.0.0")
}
```

## 配置选项详解

### 日志级别 (Level)

支持的日志级别（按严重程度从低到高）：

- `debug`: 调试信息，详细的程序执行信息
- `info`: 一般信息，程序正常运行的关键事件
- `warn`: 警告信息，可能的问题但不影响程序运行
- `error`: 错误信息，程序运行中的错误

**示例**：
```yaml
log:
  level: "warn"  # 只输出warn和error级别的日志
```

### 日志格式 (Format)

支持两种日志格式：

#### JSON格式 (`json`)
结构化日志，便于程序解析和分析：

```json
{"level":"info","message":"用户注册成功","timestamp":"2025-09-24T16:15:03.595+08:00","user_id":"123","email":"user@example.com"}
```

#### 文本格式 (`text`)
人类可读的文本格式：

```
time="2025-09-24T16:15:03+08:00" level=info msg="用户注册成功" user_id=123 email=user@example.com
```

### 日志输出 (Output)

支持三种输出模式：

#### 1. 控制台输出 (`stdout`)
```yaml
log:
  output: "stdout"
```

#### 2. 文件输出 (`file`)
```yaml
log:
  output: "file"
  file_path: "./logs/app.log"
```

#### 3. 同时输出 (`both`)
同时输出到控制台和文件：
```yaml
log:
  output: "both"
  file_path: "./logs/app.log"
```

### 文件路径 (FilePath)

指定日志文件的完整路径：

- **相对路径**: `./logs/app.log`
- **绝对路径**: `/var/log/wonder/app.log`
- **自动创建目录**: 如果目录不存在，系统会自动创建

## 实际使用示例

### 示例1：开发环境配置

```yaml
# configs/development.yaml
log:
  level: "debug"
  format: "text"
  output: "both"              # 同时输出到控制台和文件
  file_path: "./logs/dev.log"
```

### 示例2：生产环境配置

```yaml
# configs/production.yaml
log:
  level: "info"
  format: "json"
  output: "file"              # 仅输出到文件
  file_path: "/var/log/wonder/production.log"
```

### 示例3：通过环境变量

```bash
# Docker环境中使用
docker run -e LOG_LEVEL=error \
           -e LOG_FORMAT=json \
           -e LOG_OUTPUT=file \
           -e LOG_FILE_PATH=/app/logs/error.log \
           -v /host/logs:/app/logs \
           wonder:latest
```

## 代码中使用

### 全局Logger使用

```go
import (
    "context"
    "github.com/cctw-zed/wonder/pkg/logger"
)

func someFunction() {
    ctx := context.Background()

    // 不同级别的日志
    logger.LogDebug(ctx, "调试信息", "step", 1)
    logger.LogInfo(ctx, "处理用户请求", "user_id", "123")
    logger.LogWarn(ctx, "检测到潜在问题", "issue", "connection_slow")
    logger.LogError(ctx, "处理失败", "error", err.Error())
}
```

### 创建专用Logger

```go
// 为特定组件创建logger
appLogger := logger.Get().WithLayer("application").WithComponent("user_service")

// 使用专用logger
appLogger.Info(ctx, "用户服务初始化", "version", "2.0.0")
```

### 与不同框架集成

#### 在HTTP handler中使用

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    log := logger.Get().WithLayer("interface").WithComponent("user_handler")
    log.Info(ctx, "接收到用户请求", "method", r.Method, "path", r.URL.Path)

    // 处理业务逻辑...
}
```

#### 在数据库操作中使用

```go
func (r *userRepository) Create(ctx context.Context, user *User) error {
    r.log.Debug(ctx, "创建用户", "user_id", user.ID, "email", user.Email)

    if err := r.db.Create(user).Error; err != nil {
        r.log.Error(ctx, "数据库创建失败", "error", err.Error(), "user_id", user.ID)
        return err
    }

    r.log.Info(ctx, "用户创建成功", "user_id", user.ID)
    return nil
}
```

## 日志轮转和管理

### 使用logrotate（Linux系统）

创建 `/etc/logrotate.d/wonder` 配置文件：

```
/var/log/wonder/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
    postrotate
        # 可选：向应用发送信号重新加载
        # killall -USR1 wonder
    endscript
}
```

### Docker环境中的日志管理

```dockerfile
# Dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY .. .
RUN go build -o wonder ./cmd/server

# 创建日志目录
RUN mkdir -p /app/logs

# 使用卷挂载日志目录
VOLUME ["/app/logs"]

CMD ["./wonder"]
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  wonder:
    build: .
    environment:
      - LOG_OUTPUT=file
      - LOG_FILE_PATH=/app/logs/wonder.log
    volumes:
      - ./logs:/app/logs
    ports:
      - "8080:8080"
```

## 监控和告警

### 使用ELK Stack

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/wonder/*.log
  fields:
    service: wonder
    environment: production

output.elasticsearch:
  hosts: ["elasticsearch:9200"]

processors:
- decode_json_fields:
    fields: ["message"]
    target: ""
```

### 错误日志告警

可以配置监控系统监听ERROR级别的日志：

```bash
# 使用grep监控错误日志
tail -f /var/log/wonder/app.log | grep -i error | while read line; do
    # 发送告警
    echo "错误发生: $line" | mail -s "Wonder应用错误" admin@company.com
done
```

## 性能考虑

### 异步日志

对于高并发应用，可以考虑使用缓冲写入：

```go
// 在logger配置中可以添加缓冲区大小配置
// 这需要在未来版本中实现
config := logger.LogConfig{
    Level:      "info",
    Output:     "file",
    FilePath:   "./logs/app.log",
    // BufferSize: 4096,  // 未来可能的配置
}
```

### 日志级别控制

在生产环境中，避免使用DEBUG级别：

```yaml
# 生产环境
log:
  level: "info"  # 而不是 "debug"
```

## 故障排除

### 常见问题

1. **日志文件无法创建**
   - 检查文件路径的目录权限
   - 确保磁盘空间充足

2. **日志未输出到文件**
   - 检查配置中的`output`和`file_path`设置
   - 验证应用有写入权限

3. **日志格式不正确**
   - 检查`format`配置（`json` or `text`）
   - 验证JSON格式的有效性

### 调试配置

可以临时启用详细日志来调试配置问题：

```go
// 临时启用debug级别
logger.InitializeWithConfig(logger.LogConfig{
    Level:  "debug",
    Format: "text",
    Output: "both",
    FilePath: "./debug.log",
})
```

## 最佳实践

1. **结构化日志**: 使用JSON格式便于分析
2. **合适的级别**: 生产环境使用INFO级别
3. **日志轮转**: 配置日志轮转避免磁盘满
4. **监控告警**: 监控ERROR级别日志
5. **性能考虑**: 避免在高频代码中使用DEBUG日志
6. **上下文信息**: 包含足够的上下文信息便于问题定位

通过以上配置，您可以灵活地控制Wonder应用的日志输出，满足不同环境的需求。