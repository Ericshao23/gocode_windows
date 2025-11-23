package distributedlock

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	// ErrLockNotHeld is returned when trying to release or renew a lock that is not held
	ErrLockNotHeld = errors.New("lock not held")
	// ErrLockNotAcquired is returned when a lock cannot be acquired
	ErrLockNotAcquired = errors.New("failed to acquire lock")
)

type DistributedLockInfo struct {
	key         string
	value       string
	expiration  time.Duration
	mutex       sync.Mutex
	locked      bool
	stopChan    chan struct{}
	failTrys    int
	failDelay   time.Duration
	etcdSession *concurrency.Session
	etcdMutex   *concurrency.Mutex
}

type DistributedLockService interface {
	AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error)
	ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error)
	RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error
	BuildServiceType() string
}

var (
	serviceContainer = make(map[string]DistributedLockService)
	serviceMutex     sync.RWMutex
)

// NewDistributedLockInfo creates a new distributed lock info with proper initialization
func NewDistributedLockInfo(key, value string, expiration time.Duration) *DistributedLockInfo {
	return &DistributedLockInfo{
		key:        key,
		value:      value,
		expiration: expiration,
		stopChan:   make(chan struct{}),
		failTrys:   3,                      // 默认重试 3 次
		failDelay:  100 * time.Millisecond, // 默认延迟 100ms
	}
}

// SetRetry sets the number of retry attempts and delay between retries
func (dl *DistributedLockInfo) SetRetry(tries int, delay time.Duration) {
	dl.failTrys = tries
	dl.failDelay = delay
}

// RegisterService registers a new distributed lock service
func RegisterService(service DistributedLockService) {
	serviceMutex.Lock()
	defer serviceMutex.Unlock()
	serviceContainer[service.BuildServiceType()] = service
}

// GetService returns a registered service by its type
func GetService(serviceType string) (DistributedLockService, error) {
	serviceMutex.RLock()
	defer serviceMutex.RUnlock()

	if service, ok := serviceContainer[serviceType]; ok {
		return service, nil
	}
	return nil, errors.New("service not found")
}
