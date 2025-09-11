// pkg/snowflake/id/allocator.go - Dynamic node ID allocation
package id

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// NodeIDAllocator 节点ID分配器接口
type NodeIDAllocator interface {
	// AllocateNodeID 为指定服务分配节点ID
	AllocateNodeID(ctx context.Context, serviceType ServiceType) (int64, error)
	// ReleaseNodeID 释放节点ID
	ReleaseNodeID(ctx context.Context, serviceType ServiceType, nodeID int64) error
	// RefreshLease 刷新租约
	RefreshLease(ctx context.Context, serviceType ServiceType, nodeID int64) error
}

// MachineBasedAllocator 基于机器特征的分配器
type MachineBasedAllocator struct {
	machineID string
}

// NewMachineBasedAllocator 创建基于机器特征的分配器
func NewMachineBasedAllocator() (*MachineBasedAllocator, error) {
	machineID, err := generateMachineID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate machine ID: %w", err)
	}

	return &MachineBasedAllocator{
		machineID: machineID,
	}, nil
}

// AllocateNodeID 基于机器特征生成稳定的节点ID
func (m *MachineBasedAllocator) AllocateNodeID(ctx context.Context, serviceType ServiceType) (int64, error) {
	// 使用服务类型 + 机器ID生成哈希
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%s", serviceType, m.machineID)))

	// 取哈希的前8字节转换为int64，然后模1024得到实例ID
	instanceID := int64(binary.BigEndian.Uint64(hash[:8])) % 1024
	if instanceID < 0 {
		instanceID = -instanceID
	}

	nodeID := int64(serviceType) + instanceID
	return nodeID, nil
}

// ReleaseNodeID 基于机器特征的分配器不需要释放
func (m *MachineBasedAllocator) ReleaseNodeID(ctx context.Context, serviceType ServiceType, nodeID int64) error {
	// 基于机器特征的分配器不需要释放操作
	return nil
}

// RefreshLease 基于机器特征的分配器不需要刷新租约
func (m *MachineBasedAllocator) RefreshLease(ctx context.Context, serviceType ServiceType, nodeID int64) error {
	// 基于机器特征的分配器不需要租约刷新
	return nil
}

// generateMachineID 生成机器唯一标识
func generateMachineID() (string, error) {
	var parts []string

	// 1. 尝试获取MAC地址
	if macAddr, err := getMACAddress(); err == nil && macAddr != "" {
		parts = append(parts, macAddr)
	}

	// 2. 获取主机名
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		parts = append(parts, hostname)
	}

	// 3. 获取本机IP（作为备份）
	if ip, err := getLocalIP(); err == nil && ip != "" {
		parts = append(parts, ip)
	}

	// 4. 如果都失败，使用进程ID + 时间戳
	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("pid-%d-time-%d", os.Getpid(), time.Now().UnixNano()))
	}

	return strings.Join(parts, "-"), nil
}

// getMACAddress 获取第一个非回环网卡的MAC地址
func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 跳过回环和虚拟接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 跳过无MAC地址的接口
		if iface.HardwareAddr == nil {
			continue
		}

		mac := iface.HardwareAddr.String()
		if mac != "" && mac != "00:00:00:00:00:00" {
			return mac, nil
		}
	}

	return "", fmt.Errorf("no valid MAC address found")
}

// getLocalIP 获取本机IP地址
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// FallbackAllocator 后备分配器（环境变量 + 机器特征）
type FallbackAllocator struct {
	machineAllocator *MachineBasedAllocator
}

// NewFallbackAllocator 创建后备分配器
func NewFallbackAllocator() (*FallbackAllocator, error) {
	machineAllocator, err := NewMachineBasedAllocator()
	if err != nil {
		return nil, err
	}

	return &FallbackAllocator{
		machineAllocator: machineAllocator,
	}, nil
}

// AllocateNodeID 优先使用环境变量，失败则使用机器特征
func (f *FallbackAllocator) AllocateNodeID(ctx context.Context, serviceType ServiceType) (int64, error) {
	// 1. 尝试从环境变量获取
	if nodeID, _, err := ParseNodeIDFromEnv(); err == nil {
		// 验证nodeID是否在正确的服务范围内
		expectedStart := int64(serviceType)
		expectedEnd := expectedStart + 1023
		if nodeID >= expectedStart && nodeID <= expectedEnd {
			return nodeID, nil
		}
	}

	// 2. 使用机器特征生成
	return f.machineAllocator.AllocateNodeID(ctx, serviceType)
}

// ReleaseNodeID 释放节点ID
func (f *FallbackAllocator) ReleaseNodeID(ctx context.Context, serviceType ServiceType, nodeID int64) error {
	return f.machineAllocator.ReleaseNodeID(ctx, serviceType, nodeID)
}

// RefreshLease 刷新租约
func (f *FallbackAllocator) RefreshLease(ctx context.Context, serviceType ServiceType, nodeID int64) error {
	return f.machineAllocator.RefreshLease(ctx, serviceType, nodeID)
}
