// pkg/snowflake/id/etcd_allocator.go - ETCD-based dynamic node ID allocation
package id

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	// etcd key前缀
	etcdKeyPrefix = "/snowflake/services"

	// 租约配置
	defaultLeaseTimeout  = 30 * time.Second // 租约过期时间
	defaultRenewInterval = 10 * time.Second // 续租间隔
	defaultRetryInterval = 3 * time.Second  // 重试间隔

	// 分配超时
	allocateTimeout = 10 * time.Second
)

// InstanceInfo 实例信息
type InstanceInfo struct {
	NodeID      int64     `json:"node_id"`
	ServiceType string    `json:"service_type"`
	InstanceID  string    `json:"instance_id"`
	Hostname    string    `json:"hostname"`
	StartTime   time.Time `json:"start_time"`
	LastRenew   time.Time `json:"last_renew"`
}

// EtcdAllocator 基于etcd的节点ID分配器
type EtcdAllocator struct {
	client        *clientv3.Client
	leaseTimeout  time.Duration
	renewInterval time.Duration
	retryInterval time.Duration

	// 当前分配的nodeID和租约
	mu           sync.RWMutex
	nodeID       int64
	leaseID      clientv3.LeaseID
	serviceType  ServiceType
	instanceInfo *InstanceInfo

	// 控制续租的context和cancel
	renewCtx    context.Context
	renewCancel context.CancelFunc

	// 状态
	isAllocated bool
	renewDone   chan struct{}
}

// NewEtcdAllocator 创建etcd分配器
func NewEtcdAllocator(endpoints []string, opts ...EtcdOption) (*EtcdAllocator, error) {
	config := &EtcdConfig{
		Endpoints:     endpoints,
		LeaseTimeout:  defaultLeaseTimeout,
		RenewInterval: defaultRenewInterval,
		RetryInterval: defaultRetryInterval,
	}

	for _, opt := range opts {
		opt(config)
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	allocator := &EtcdAllocator{
		client:        client,
		leaseTimeout:  config.LeaseTimeout,
		renewInterval: config.RenewInterval,
		retryInterval: config.RetryInterval,
		renewDone:     make(chan struct{}),
	}

	return allocator, nil
}

// EtcdConfig etcd配置
type EtcdConfig struct {
	Endpoints     []string
	LeaseTimeout  time.Duration
	RenewInterval time.Duration
	RetryInterval time.Duration
}

// EtcdOption etcd配置选项
type EtcdOption func(*EtcdConfig)

// WithLeaseTimeout 设置租约超时时间
func WithLeaseTimeout(timeout time.Duration) EtcdOption {
	return func(c *EtcdConfig) {
		c.LeaseTimeout = timeout
	}
}

// WithRenewInterval 设置续租间隔
func WithRenewInterval(interval time.Duration) EtcdOption {
	return func(c *EtcdConfig) {
		c.RenewInterval = interval
	}
}

// AllocateNodeID 分配节点ID
func (e *EtcdAllocator) AllocateNodeID(ctx context.Context, serviceType ServiceType) (int64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.isAllocated {
		return e.nodeID, nil
	}

	// 创建分配上下文
	allocCtx, cancel := context.WithTimeout(ctx, allocateTimeout)
	defer cancel()

	// 使用分布式锁确保原子分配
	lockKey := e.getLockKey(serviceType)
	session, err := concurrency.NewSession(e.client, concurrency.WithTTL(int(e.leaseTimeout.Seconds())))
	if err != nil {
		return 0, fmt.Errorf("failed to create etcd session: %w", err)
	}
	defer session.Close()

	mutex := concurrency.NewMutex(session, lockKey)
	if err := mutex.Lock(allocCtx); err != nil {
		return 0, fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer mutex.Unlock(allocCtx)

	// 查找可用的nodeID
	nodeID, err := e.findAvailableNodeID(allocCtx, serviceType)
	if err != nil {
		return 0, fmt.Errorf("failed to find available node ID: %w", err)
	}

	// 创建租约
	lease, err := e.client.Grant(allocCtx, int64(e.leaseTimeout.Seconds()))
	if err != nil {
		return 0, fmt.Errorf("failed to create lease: %w", err)
	}

	// 注册nodeID
	if err := e.registerNodeID(allocCtx, serviceType, nodeID, lease.ID); err != nil {
		// 撤销租约
		e.client.Revoke(context.Background(), lease.ID)
		return 0, fmt.Errorf("failed to register node ID: %w", err)
	}

	// 保存状态
	e.nodeID = nodeID
	e.leaseID = lease.ID
	e.serviceType = serviceType
	e.isAllocated = true

	// 启动续租goroutine
	e.startRenewLease()

	log.Printf("Successfully allocated node ID %d for service %s", nodeID, serviceType)
	return nodeID, nil
}

// findAvailableNodeID 查找可用的nodeID
func (e *EtcdAllocator) findAvailableNodeID(ctx context.Context, serviceType ServiceType) (int64, error) {
	// 服务的nodeID范围
	startID := int64(serviceType)
	endID := startID + 1023

	// 获取已分配的nodeID列表
	allocatedKey := e.getAllocatedKey(serviceType)
	resp, err := e.client.Get(ctx, allocatedKey, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}

	// 构建已分配nodeID的map
	allocated := make(map[int64]bool)
	for _, kv := range resp.Kvs {
		var info InstanceInfo
		if err := json.Unmarshal(kv.Value, &info); err == nil {
			allocated[info.NodeID] = true
		}
	}

	// 查找第一个可用的nodeID
	for nodeID := startID; nodeID <= endID; nodeID++ {
		if !allocated[nodeID] {
			return nodeID, nil
		}
	}

	return 0, fmt.Errorf("no available node ID in range [%d, %d] for service %s",
		startID, endID, serviceType)
}

// registerNodeID 注册nodeID
func (e *EtcdAllocator) registerNodeID(ctx context.Context, serviceType ServiceType,
	nodeID int64, leaseID clientv3.LeaseID) error {

	// 生成实例信息
	hostname, _ := os.Hostname()
	instanceInfo := &InstanceInfo{
		NodeID:      nodeID,
		ServiceType: serviceType.String(),
		InstanceID:  generateInstanceID(),
		Hostname:    hostname,
		StartTime:   time.Now(),
		LastRenew:   time.Now(),
	}

	data, err := json.Marshal(instanceInfo)
	if err != nil {
		return err
	}

	// 注册到etcd
	key := e.getNodeKey(serviceType, nodeID)
	_, err = e.client.Put(ctx, key, string(data), clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}

	e.instanceInfo = instanceInfo
	return nil
}

// startRenewLease 启动续租
func (e *EtcdAllocator) startRenewLease() {
	e.renewCtx, e.renewCancel = context.WithCancel(context.Background())

	go func() {
		defer close(e.renewDone)

		ticker := time.NewTicker(e.renewInterval)
		defer ticker.Stop()

		for {
			select {
			case <-e.renewCtx.Done():
				return
			case <-ticker.C:
				if err := e.renewLease(); err != nil {
					log.Printf("Failed to renew lease for node %d: %v", e.nodeID, err)
					// 在实际应用中，这里可能需要触发重新分配
				}
			}
		}
	}()
}

// renewLease 续租
func (e *EtcdAllocator) renewLease() error {
	e.mu.RLock()
	leaseID := e.leaseID
	nodeID := e.nodeID
	serviceType := e.serviceType
	e.mu.RUnlock()

	if leaseID == 0 {
		return fmt.Errorf("no lease to renew")
	}

	// 续租
	_, err := e.client.KeepAliveOnce(e.renewCtx, leaseID)
	if err != nil {
		return err
	}

	// 更新实例信息中的续租时间
	if e.instanceInfo != nil {
		e.instanceInfo.LastRenew = time.Now()
		data, _ := json.Marshal(e.instanceInfo)
		key := e.getNodeKey(serviceType, nodeID)
		e.client.Put(e.renewCtx, key, string(data), clientv3.WithLease(leaseID))
	}

	return nil
}

// ReleaseNodeID 释放节点ID
func (e *EtcdAllocator) ReleaseNodeID(ctx context.Context, serviceType ServiceType, nodeID int64) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.isAllocated || e.nodeID != nodeID {
		return nil // 已经释放或不是当前分配的ID
	}

	// 停止续租
	if e.renewCancel != nil {
		e.renewCancel()
		<-e.renewDone // 等待续租goroutine结束
	}

	// 撤销租约
	if e.leaseID != 0 {
		e.client.Revoke(ctx, e.leaseID)
	}

	// 清理状态
	e.nodeID = 0
	e.leaseID = 0
	e.isAllocated = false
	e.instanceInfo = nil

	log.Printf("Successfully released node ID %d for service %s", nodeID, serviceType)
	return nil
}

// RefreshLease 刷新租约（外部调用）
func (e *EtcdAllocator) RefreshLease(ctx context.Context, serviceType ServiceType, nodeID int64) error {
	return e.renewLease()
}

// Close 关闭分配器
func (e *EtcdAllocator) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 释放当前的nodeID
	if e.isAllocated {
		e.ReleaseNodeID(context.Background(), e.serviceType, e.nodeID)
	}

	// 关闭etcd客户端
	return e.client.Close()
}

// 辅助方法
func (e *EtcdAllocator) getLockKey(serviceType ServiceType) string {
	return path.Join(etcdKeyPrefix, serviceType.String(), "lock")
}

func (e *EtcdAllocator) getAllocatedKey(serviceType ServiceType) string {
	return path.Join(etcdKeyPrefix, serviceType.String(), "allocated")
}

func (e *EtcdAllocator) getNodeKey(serviceType ServiceType, nodeID int64) string {
	return path.Join(etcdKeyPrefix, serviceType.String(), "allocated", fmt.Sprintf("%d", nodeID))
}

// generateInstanceID 生成实例ID
func generateInstanceID() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s-%d-%d", hostname, os.Getpid(), time.Now().UnixNano())
}
