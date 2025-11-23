package distributedlock

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type EtcdLock struct {
	session   *concurrency.Session
	client    *clientv3.Client
	endpoints []string
}

// NewEtcdLock creates a new etcd distributed lock
func NewEtcdLock(session *concurrency.Session) *EtcdLock {
	return &EtcdLock{
		session:   session,
		endpoints: []string{"localhost:2379"}, // Default endpoint
	}
}

// NewEtcdLockWithEndpoints creates a new etcd distributed lock with custom endpoints
func NewEtcdLockWithEndpoints(endpoints []string) (*EtcdLock, error) {
	if len(endpoints) == 0 {
		endpoints = []string{"localhost:2379"}
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	session, err := concurrency.NewSession(client)
	if err != nil {
		client.Close()
		return nil, err
	}

	return &EtcdLock{
		session:   session,
		client:    client,
		endpoints: endpoints,
	}, nil
}

// AcquireLock attempts to acquire a distributed lock
func (e *EtcdLock) AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	if lockInfo.etcdSession == nil {
		// 使用已有的 session 或创建新的
		if e.session != nil {
			lockInfo.etcdSession = e.session
		} else {
			// Create a new session if one doesn't exist
			client, err := clientv3.New(clientv3.Config{
				Endpoints:   e.endpoints,
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				return false, err
			}

			session, err := concurrency.NewSession(client)
			if err != nil {
				client.Close()
				return false, err
			}
			lockInfo.etcdSession = session
		}
	}

	// Create a new mutex for this lock
	mutex := concurrency.NewMutex(lockInfo.etcdSession, lockInfo.key)
	lockInfo.etcdMutex = mutex

	// Try to acquire the lock
	err := mutex.TryLock(ctx)
	if err != nil {
		if err == concurrency.ErrLocked {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// ReleaseLock releases the distributed lock
func (e *EtcdLock) ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	if lockInfo.etcdMutex == nil {
		return false, ErrLockNotHeld
	}

	// Unlock the mutex
	err := lockInfo.etcdMutex.Unlock(ctx)
	if err != nil {
		return false, err
	}

	// Clean up the session if it exists
	if lockInfo.etcdSession != nil {
		err = lockInfo.etcdSession.Close()
		if err != nil {
			return false, err
		}
		lockInfo.etcdSession = nil
	}

	lockInfo.etcdMutex = nil
	return true, nil
}

// RenewLock implements lease keep-alive for the lock
func (e *EtcdLock) RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error {
	// etcd's concurrency.Mutex handles lease renewal automatically
	// as long as the session is kept alive
	return nil
}

// BuildServiceType returns the type of lock service
func (e *EtcdLock) BuildServiceType() string {
	return "etcd"
}
