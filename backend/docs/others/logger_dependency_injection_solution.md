# Logger Dependency Injection Solution

## 问题描述

在Go项目中，如果在初始化service或repository时忘记传入logger组件，会导致运行时nil pointer panic。这个问题在测试环境中尤其容易出现，因为测试代码经常需要手动构造结构体。

## 问题的根本原因

1. **手动创建结构体的风险**: 在测试中直接使用结构体字面量创建实例时，容易忘记初始化某些字段
2. **无编译时检查**: Go编译器无法检测到未初始化的结构体字段，只有运行时才会发生panic
3. **隐式依赖**: logger作为内部字段，依赖关系不够明显

## 经典项目的解决方案

通过研究Kubernetes、Docker、etcd等经典Go项目，发现以下几种常见模式：

### 1. Constructor Pattern (构造器模式)
强制使用构造函数创建实例，确保所有必需字段都被正确初始化。

### 2. Explicit Logger Injection (显式Logger注入)
将logger作为构造函数的显式参数，使依赖关系更加明确。

### 3. Validation in Constructors (构造函数中的验证)
在构造函数中添加参数验证，防止nil值传入。

## 我们的解决方案

### 实现的改进

1. **双构造函数模式**: 提供两个构造函数
   - `NewUserService()`: 使用默认logger配置的便捷构造函数
   - `NewUserServiceWithLogger()`: 接受显式logger参数的完整构造函数

2. **参数验证**: 在构造函数中验证所有必需参数不为nil，如果为nil则panic并提供清晰的错误信息

3. **向后兼容**: 保持原有API不变，同时提供更安全的替代方案

### 代码示例

#### UserService构造器
```go
func NewUserService(repo user.UserRepository, idGen id.Generator) user.UserService {
    return NewUserServiceWithLogger(repo, idGen, logger.Get().WithLayer("application").WithComponent("user_service"))
}

func NewUserServiceWithLogger(repo user.UserRepository, idGen id.Generator, log logger.Logger) user.UserService {
    if repo == nil {
        panic("user repository cannot be nil")
    }
    if idGen == nil {
        panic("ID generator cannot be nil")
    }
    if log == nil {
        panic("logger cannot be nil")
    }

    return &userService{
        repo:  repo,
        idGen: idGen,
        log:   log,
    }
}
```

#### UserRepository构造器
```go
func NewUserRepository(db *gorm.DB) user.UserRepository {
    return NewUserRepositoryWithLogger(db, logger.Get().WithLayer("infrastructure").WithComponent("user_repository"))
}

func NewUserRepositoryWithLogger(db *gorm.DB, log logger.Logger) user.UserRepository {
    if db == nil {
        panic("database connection cannot be nil")
    }
    if log == nil {
        panic("logger cannot be nil")
    }

    return &userRepository{
        db:  db,
        log: log,
    }
}
```

### 测试改进

更新测试代码使用安全的构造函数：

```go
// 之前的不安全做法
service := &userService{
    repo:  mockRepo,
    idGen: mockIDGen,
    // 容易忘记初始化log字段 -> panic
}

// 现在的安全做法
service := NewUserServiceWithLogger(mockRepo, mockIDGen, logger.Get().WithLayer("application").WithComponent("user_service"))
```

## 解决方案的优势

### 1. **Fail Fast原则**
- 在对象创建时就发现问题，而不是在使用时才panic
- 提供清晰的错误信息，便于快速定位问题

### 2. **编译时安全**
- 强制使用构造函数，避免遗漏字段初始化
- 通过类型系统确保依赖的正确性

### 3. **测试友好**
- 为测试提供显式logger注入的接口
- 防止测试中因忘记初始化logger而panic

### 4. **向后兼容**
- 保持现有API不变
- 提供渐进式迁移路径

### 5. **清晰的依赖关系**
- 显式声明所有依赖
- 便于理解和维护代码

## 最佳实践

### 1. 构造函数验证
```go
func NewService(deps ...interface{}) *Service {
    // 验证所有依赖不为nil
    for i, dep := range deps {
        if dep == nil {
            panic(fmt.Sprintf("dependency %d cannot be nil", i))
        }
    }
    // ... 创建服务
}
```

### 2. 接口优于实现
```go
// 好：接受接口
func NewService(logger logger.Logger) *Service

// 不推荐：接受具体实现
func NewService(logger *zap.Logger) *Service
```

### 3. 分层初始化
```go
// 在每层使用适当的logger配置
appLogger := logger.Get().WithLayer("application")
infraLogger := logger.Get().WithLayer("infrastructure")
```

### 4. 测试中的明确依赖
```go
func TestService(t *testing.T) {
    logger.Initialize() // 确保logger已初始化
    service := NewServiceWithLogger(mockRepo, testLogger)
    // ... 测试逻辑
}
```

## 验证结果

所有测试都通过，包括：
- ✅ 构造函数验证测试
- ✅ 正常功能测试
- ✅ 错误处理测试
- ✅ 并发测试

这个解决方案有效防止了nil logger panic的问题，同时保持了代码的可维护性和向后兼容性。