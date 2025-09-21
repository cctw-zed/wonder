// pkg/snowflake/id/generator.go - Distributed Snowflake ID generator
package id

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
	"strconv"
	"sync"
)

// ServiceType 服务类型枚举
type ServiceType int

const (
	// Service type segments (0-1023 nodes per service)
	ServiceTypeUser    ServiceType = 0    // Node IDs: 0-1023
	ServiceTypeOrder   ServiceType = 1024 // Node IDs: 1024-2047
	ServiceTypePayment ServiceType = 2048 // Node IDs: 2048-3071
	ServiceTypeAuth    ServiceType = 3072 // Node IDs: 3072-4095
	ServiceTypeGateway ServiceType = 4096 // Node IDs: 4096-5119
	// Add more services as needed...
)

// String returns the service type name
func (s ServiceType) String() string {
	switch s {
	case ServiceTypeUser:
		return "user"
	case ServiceTypeOrder:
		return "order"
	case ServiceTypePayment:
		return "payment"
	case ServiceTypeAuth:
		return "auth"
	case ServiceTypeGateway:
		return "gateway"
	default:
		return "unknown"
	}
}

// ParseServiceType parses service type from string
func ParseServiceType(s string) (ServiceType, error) {
	switch s {
	case "user":
		return ServiceTypeUser, nil
	case "order":
		return ServiceTypeOrder, nil
	case "payment":
		return ServiceTypePayment, nil
	case "auth":
		return ServiceTypeAuth, nil
	case "gateway":
		return ServiceTypeGateway, nil
	default:
		return 0, fmt.Errorf("unknown service type: %s", s)
	}
}

// Generator ID生成器接口
type Generator interface {
	Generate() string
	GenerateInt64() int64
	GetNodeID() int64
	GetServiceType() ServiceType
}

// SnowflakeGenerator 分布式雪花ID生成器
type SnowflakeGenerator struct {
	node        *snowflake.Node
	nodeID      int64
	serviceType ServiceType
	instanceID  int64
	allocator   NodeIDAllocator // 可选的分配器，用于生命周期管理
}

var (
	defaultGenerator Generator
	once             sync.Once
)

// NewSnowflakeGenerator 创建基于节点ID的生成器
func NewSnowflakeGenerator(nodeID int64) (Generator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	// 从nodeID反推服务类型
	serviceType := ServiceType((nodeID / 1024) * 1024)
	instanceID := nodeID % 1024

	return &SnowflakeGenerator{
		node:        node,
		nodeID:      nodeID,
		serviceType: serviceType,
		instanceID:  instanceID,
	}, nil
}

// NewSnowflakeGeneratorForService 为指定服务类型创建生成器
func NewSnowflakeGeneratorForService(serviceType ServiceType, instanceID int64) (Generator, error) {
	if instanceID < 0 || instanceID >= 1024 {
		return nil, fmt.Errorf("instance ID must be in range [0, 1023], got: %d", instanceID)
	}

	nodeID := int64(serviceType) + instanceID
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node for service %s, instance %d: %w", serviceType, instanceID, err)
	}

	return &SnowflakeGenerator{
		node:        node,
		nodeID:      nodeID,
		serviceType: serviceType,
		instanceID:  instanceID,
	}, nil
}

// CalculateNodeID 根据服务类型和实例ID计算节点ID
func CalculateNodeID(serviceType ServiceType, instanceID int64) (int64, error) {
	if instanceID < 0 || instanceID >= 1024 {
		return 0, fmt.Errorf("instance ID must be in range [0, 1023], got: %d", instanceID)
	}
	return int64(serviceType) + instanceID, nil
}

// ParseNodeIDFromEnv 从环境变量解析节点ID配置
func ParseNodeIDFromEnv() (int64, ServiceType, error) {
	// 优先使用 SERVICE_TYPE 和 INSTANCE_ID
	serviceTypeStr := os.Getenv("SERVICE_TYPE")
	instanceIDStr := os.Getenv("INSTANCE_ID")

	if serviceTypeStr != "" && instanceIDStr != "" {
		serviceType, err := ParseServiceType(serviceTypeStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid SERVICE_TYPE: %w", err)
		}

		instanceID, err := strconv.ParseInt(instanceIDStr, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid INSTANCE_ID: %w", err)
		}

		nodeID, err := CalculateNodeID(serviceType, instanceID)
		if err != nil {
			return 0, 0, err
		}

		return nodeID, serviceType, nil
	}

	// 回退到直接使用 NODE_ID
	nodeIDStr := os.Getenv("NODE_ID")
	if nodeIDStr == "" {
		// 默认为用户服务的第一个实例
		return int64(ServiceTypeUser), ServiceTypeUser, nil
	}

	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid NODE_ID: %w", err)
	}

	serviceType := ServiceType((nodeID / 1024) * 1024)
	return nodeID, serviceType, nil
}

// Generate 生成字符串ID
func (s *SnowflakeGenerator) Generate() string {
	return s.node.Generate().String()
}

// GenerateInt64 生成int64 ID
func (s *SnowflakeGenerator) GenerateInt64() int64 {
	return s.node.Generate().Int64()
}

// GetNodeID 获取节点ID
func (s *SnowflakeGenerator) GetNodeID() int64 {
	return s.nodeID
}

// GetServiceType 获取服务类型
func (s *SnowflakeGenerator) GetServiceType() ServiceType {
	return s.serviceType
}

// Close 关闭生成器并释放资源
func (s *SnowflakeGenerator) Close() error {
	if s.allocator != nil {
		return s.allocator.ReleaseNodeID(context.Background(), s.serviceType, s.nodeID)
	}
	return nil
}

// InitDefault 初始化默认生成器（兼容旧版本）
func InitDefault(nodeID int64) error {
	var err error
	once.Do(func() {
		defaultGenerator, err = NewSnowflakeGenerator(nodeID)
	})
	return err
}

// InitDefaultFromEnv 从环境变量初始化默认生成器
func InitDefaultFromEnv() error {
	nodeID, _, err := ParseNodeIDFromEnv()
	if err != nil {
		return fmt.Errorf("failed to parse node ID from environment: %w", err)
	}
	return InitDefault(nodeID)
}

// InitDefaultForService 为指定服务初始化默认生成器
func InitDefaultForService(serviceType ServiceType, instanceID int64) error {
	var err error
	once.Do(func() {
		defaultGenerator, err = NewSnowflakeGeneratorForService(serviceType, instanceID)
	})
	return err
}

// InitDefaultWithAllocator 使用指定分配器初始化默认生成器
func InitDefaultWithAllocator(ctx context.Context, serviceType ServiceType, allocator NodeIDAllocator) error {
	var err error
	once.Do(func() {
		defaultGenerator, err = NewSnowflakeGeneratorWithAllocator(ctx, serviceType, allocator)
	})
	return err
}

// NewSnowflakeGeneratorWithAllocator 使用分配器创建生成器
func NewSnowflakeGeneratorWithAllocator(ctx context.Context, serviceType ServiceType, allocator NodeIDAllocator) (Generator, error) {
	// 通过分配器获取nodeID
	nodeID, err := allocator.AllocateNodeID(ctx, serviceType)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate node ID: %w", err)
	}

	// 创建snowflake节点
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		// 分配失败时释放nodeID
		allocator.ReleaseNodeID(context.Background(), serviceType, nodeID)
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}

	instanceID := nodeID % 1024
	generator := &SnowflakeGenerator{
		node:        node,
		nodeID:      nodeID,
		serviceType: serviceType,
		instanceID:  instanceID,
		allocator:   allocator,
	}

	return generator, nil
}

// Generate 使用默认生成器生成ID
func Generate() string {
	if defaultGenerator == nil {
		panic("default generator not initialized")
	}
	return defaultGenerator.Generate()
}

// GenerateInt64 使用默认生成器生成int64 ID
func GenerateInt64() int64 {
	if defaultGenerator == nil {
		panic("default generator not initialized")
	}
	return defaultGenerator.GenerateInt64()
}

// GetDefault 获取默认生成器
func GetDefault() Generator {
	if defaultGenerator == nil {
		panic("default generator not initialized")
	}
	return defaultGenerator
}
