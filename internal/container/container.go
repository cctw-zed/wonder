package container

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"

	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/internal/infrastructure/config"
	"github.com/cctw-zed/wonder/internal/infrastructure/database"
	"github.com/cctw-zed/wonder/internal/infrastructure/repository"
	"github.com/cctw-zed/wonder/internal/interfaces/http"
	"github.com/cctw-zed/wonder/pkg/snowflake/id"
)

type Container struct {
	Config        *config.Config
	UserHandler   *http.UserHandler
	Database      *database.Connection
	nodeAllocator id.NodeIDAllocator // 节点ID分配器，用于优雅关闭时释放资源
}

func NewContainer() (*Container, error) {
	return NewContainerWithContext(context.Background())
}

// NewContainerWithContext 使用上下文创建容器，支持动态nodeID分配
func NewContainerWithContext(ctx context.Context) (*Container, error) {
	return NewContainerWithConfig(ctx, "")
}

// NewContainerForEnvironment 为指定环境创建容器
func NewContainerForEnvironment(ctx context.Context, environment string) (*Container, error) {
	// Load environment-specific configuration from configs directory
	cfg, err := config.LoadForEnvironment(environment, "./configs")
	if err != nil {
		return nil, fmt.Errorf("failed to load config for environment %s: %w", environment, err)
	}

	// 检测ID分配策略
	allocator := createNodeIDAllocator(ctx, cfg)

	// Initialize database connection using config
	dbConn, err := database.NewConnection(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run database migrations
	migrator := database.NewMigrator(dbConn.DB())
	if err := migrator.MigrateAll(); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// 根据分配器类型初始化ID生成器
	if allocator != nil {
		// 使用动态分配器
		serviceType := getServiceTypeFromConfig(cfg)
		if err := id.InitDefaultWithAllocator(ctx, serviceType, allocator); err != nil {
			return nil, fmt.Errorf("failed to init ID generator with allocator: %w", err)
		}
	} else {
		// 使用配置文件中的ID配置
		if err := id.InitDefaultForService(getServiceTypeFromConfig(cfg), cfg.ID.InstanceID); err != nil {
			return nil, fmt.Errorf("failed to init distributed ID generator: %w", err)
		}
	}

	// 后续组件可以直接使用 id.Generate()
	userRepo := repository.NewUserRepository(dbConn.DB())
	idGen := id.GetDefault()
	userService := service.NewUserService(userRepo, idGen)
	userHandler := http.NewUserHandler(userService)

	return &Container{
		Config:        cfg,
		UserHandler:   userHandler,
		Database:      dbConn,
		nodeAllocator: allocator,
	}, nil
}

// NewContainerWithConfig 使用配置文件路径创建容器
func NewContainerWithConfig(ctx context.Context, configPath string) (*Container, error) {
	// Load configuration
	var cfg *config.Config
	var err error
	if configPath != "" {
		cfg, err = config.Load(configPath)
	} else {
		// Use configs directory by default
		cfg, err = config.Load("./configs")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 检测ID分配策略
	allocator := createNodeIDAllocator(ctx, cfg)

	// Initialize database connection using config
	dbConn, err := database.NewConnection(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run database migrations
	migrator := database.NewMigrator(dbConn.DB())
	if err := migrator.MigrateAll(); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// 根据分配器类型初始化ID生成器
	if allocator != nil {
		// 使用动态分配器
		serviceType := getServiceTypeFromConfig(cfg)
		if err := id.InitDefaultWithAllocator(ctx, serviceType, allocator); err != nil {
			return nil, fmt.Errorf("failed to init ID generator with allocator: %w", err)
		}
	} else {
		// 使用配置文件中的ID配置
		if err := id.InitDefaultForService(getServiceTypeFromConfig(cfg), cfg.ID.InstanceID); err != nil {
			return nil, fmt.Errorf("failed to init distributed ID generator: %w", err)
		}
	}

	// 后续组件可以直接使用 id.Generate()
	userRepo := repository.NewUserRepository(dbConn.DB())
	idGen := id.GetDefault()
	userService := service.NewUserService(userRepo, idGen)
	userHandler := http.NewUserHandler(userService)

	return &Container{
		Config:        cfg,
		UserHandler:   userHandler,
		Database:      dbConn,
		nodeAllocator: allocator,
	}, nil
}

// createNodeIDAllocator 创建节点ID分配器
func createNodeIDAllocator(ctx context.Context, cfg *config.Config) id.NodeIDAllocator {
	// 检查是否配置了etcd
	etcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpoints != "" {
		endpoints := strings.Split(etcdEndpoints, ",")

		// 创建etcd分配器
		allocator, err := id.NewEtcdAllocator(endpoints)
		if err != nil {
			fmt.Printf("Failed to create etcd allocator: %v, falling back to static allocation\n", err)
			return nil
		}

		fmt.Println("Using etcd-based dynamic node ID allocation")
		return allocator
	}

	// 检查是否使用机器特征分配
	if os.Getenv("USE_MACHINE_BASED_ID") == "true" {
		allocator, err := id.NewMachineBasedAllocator()
		if err != nil {
			fmt.Printf("Failed to create machine-based allocator: %v, falling back to static allocation\n", err)
			return nil
		}

		fmt.Println("Using machine-based node ID allocation")
		return allocator
	}

	// 检查是否使用后备分配器
	if os.Getenv("USE_FALLBACK_ALLOCATOR") == "true" {
		allocator, err := id.NewFallbackAllocator()
		if err != nil {
			fmt.Printf("Failed to create fallback allocator: %v, falling back to static allocation\n", err)
			return nil
		}

		fmt.Println("Using fallback node ID allocation")
		return allocator
	}

	fmt.Println("Using static node ID allocation from environment variables")
	return nil
}

// getServiceTypeFromConfig 从配置获取服务类型
func getServiceTypeFromConfig(cfg *config.Config) id.ServiceType {
	serviceType, err := id.ParseServiceType(cfg.ID.ServiceType)
	if err != nil {
		fmt.Printf("Invalid service type in config: %v, using default (user)\n", err)
		return id.ServiceTypeUser
	}
	return serviceType
}

// getServiceTypeFromEnv 从环境变量获取服务类型（保留向后兼容）
func getServiceTypeFromEnv() id.ServiceType {
	serviceTypeStr := os.Getenv("SERVICE_TYPE")
	if serviceTypeStr == "" {
		return id.ServiceTypeUser // 默认为用户服务
	}

	serviceType, err := id.ParseServiceType(serviceTypeStr)
	if err != nil {
		fmt.Printf("Invalid SERVICE_TYPE: %v, using default (user)\n", err)
		return id.ServiceTypeUser
	}

	return serviceType
}

// Close 优雅关闭容器，释放资源
func (c *Container) Close() error {
	if c.nodeAllocator != nil {
		// 如果是etcd分配器，需要关闭连接
		if etcdAllocator, ok := c.nodeAllocator.(*id.EtcdAllocator); ok {
			return etcdAllocator.Close()
		}
		// 其他分配器也可以添加关闭逻辑
	}
	return nil
}

// NewContainerForService 为指定服务类型创建容器（静态分配方式）
func NewContainerForService(db *gorm.DB, serviceType id.ServiceType, instanceID int64) *Container {
	// 为指定服务类型初始化ID生成器
	if err := id.InitDefaultForService(serviceType, instanceID); err != nil {
		panic(fmt.Sprintf("Failed to init ID generator for service %s: %v", serviceType, err))
	}

	// 后续组件可以直接使用 id.Generate()
	userRepo := repository.NewUserRepository(db)
	idGen := id.GetDefault()
	userService := service.NewUserService(userRepo, idGen)
	userHandler := http.NewUserHandler(userService)

	return &Container{
		UserHandler:   userHandler,
		nodeAllocator: nil, // 静态分配不需要分配器
	}
}

// getNodeID 兼容旧版本的节点ID获取（已废弃）
// Deprecated: 使用 id.ParseNodeIDFromEnv() 替代
func getNodeID() int64 {
	nodeID, _, err := id.ParseNodeIDFromEnv()
	if err != nil {
		panic(fmt.Sprintf("Failed to parse node ID: %v", err))
	}
	return nodeID
}
