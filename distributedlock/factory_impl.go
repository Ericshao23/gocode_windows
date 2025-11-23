package distributedlock

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// LockType represents the type of distributed lock
type LockType string

const (
	// RedisLockType represents a Redis-based distributed lock
	RedisLockType LockType = "redis"
	// EtcdLockType represents an etcd-based distributed lock
	EtcdLockType LockType = "etcd"
	// MySQLLockType represents a MySQL-based distributed lock
	MySQLLockType LockType = "mysql"
	// ZookeeperLockType represents a ZooKeeper-based distributed lock
	ZookeeperLockType LockType = "zookeeper"
)

// NewDistributedLock creates a new distributed lock of the specified type
func NewDistributedLock(lockType LockType, config interface{}) (DistributedLockService, error) {
	switch lockType {
	case RedisLockType:
		return newRedisLock(config)
	case EtcdLockType:
		return newEtcdLock(config)
	case MySQLLockType:
		return newMySQLLock(config)
	case ZookeeperLockType:
		return newZookeeperLock(config)
	default:
		return nil, fmt.Errorf("unsupported lock type: %s", lockType)
	}
}

// newRedisLock creates a new Redis distributed lock
func newRedisLock(config interface{}) (*RedisLock, error) {
	cfg, ok := config.(RedisConfig)
	if !ok {
		return nil, errors.New("invalid Redis config")
	}

	// 使用第一个地址，如果没有则使用默认值
	addr := "localhost:6379"
	if len(cfg.Addrs) > 0 {
		addr = cfg.Addrs[0]
	}

	return &RedisLock{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}, nil
}

// newEtcdLock creates a new etcd distributed lock
func newEtcdLock(config interface{}) (*EtcdLock, error) {
	// 尝试多种配置类型
	switch cfg := config.(type) {
	case clientv3.Config:
		client, err := clientv3.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create etcd client: %v", err)
		}

		session, err := concurrency.NewSession(client)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("failed to create etcd session: %v", err)
		}

		return &EtcdLock{
			session:   session,
			client:    client,
			endpoints: cfg.Endpoints,
		}, nil

	case []string:
		// 如果传入的是端点列表
		return NewEtcdLockWithEndpoints(cfg)

	case *concurrency.Session:
		// 如果直接传入 session
		return &EtcdLock{
			session:   cfg,
			endpoints: []string{"localhost:2379"},
		}, nil

	default:
		return nil, errors.New("invalid etcd config: expected clientv3.Config, []string, or *concurrency.Session")
	}
}

// newMySQLLock creates a new MySQL distributed lock
func newMySQLLock(config interface{}) (*MySQLLock, error) {
	dsn, ok := config.(string)
	if !ok {
		dsnCfg, ok := config.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid MySQL config: expected DSN string or config map")
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			dsnCfg["user"],
			dsnCfg["password"],
			dsnCfg["host"],
			dsnCfg["port"],
			dsnCfg["dbname"],
		)
	}

	return NewMySQLLock(dsn)
}

// newZookeeperLock creates a new ZooKeeper distributed lock
func newZookeeperLock(config interface{}) (*ZookeeperLock, error) {
	cfg, ok := config.(ZooKeeperConfig)
	if !ok {
		return nil, errors.New("invalid ZooKeeper config")
	}

	servers := cfg.Servers
	if len(servers) == 0 {
		servers = []string{"127.0.0.1:2181"}
	}

	sessionTimeout := 10 * time.Second
	if cfg.SessionTimeout > 0 {
		sessionTimeout = cfg.SessionTimeout
	}

	prefix := "/locks"
	if cfg.Prefix != "" {
		prefix = cfg.Prefix
	}

	return NewZookeeperLock(servers, sessionTimeout, prefix)
}

// RedisConfig holds the configuration for Redis lock
type RedisConfig struct {
	Addrs    []string
	Password string
	DB       int
}

// ZooKeeperConfig holds the configuration for ZooKeeper lock
type ZooKeeperConfig struct {
	Servers        []string
	SessionTimeout time.Duration
	Prefix         string
}
