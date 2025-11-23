# 分布式锁
## 1. 什么是锁
锁是一种同步机制，用于在多任务、多线程或多进程的环境中，控制对共享资源的并发访问。它可以防止多个执行流（如线程）同时访问或修改同一份数据，从而避免由此引发的数据竞争和数据不一致问题。锁是协调并发环境中共享资源访问、保证数据一致性的核心工具。

## 2. 分布式锁
定义：一把锁协调各独立主机节点。
释义：基于同步读写区分唯一key能力的通用组件（mysql/redis/etcd/zookeeperdeng等）实现具有加锁和解锁方法的一个服务，此服务针对key可以保证互斥性、超时释放、解锁只能解自己加的锁。

### QA
```
1. 如果刚加完的锁没有执行到释放锁，系统重启了，这样就不会释放锁了怎么办？
答：增加锁超时，超过超时时间后自动释放key。

2. 超时时间设置多少合适？
答：根据业务确定，没有固定时间，一般为业务执行时间P99的1.5倍。

3. 如果部分业务耗时抖动导致耗时增加超过超时怎么办？
答：自动续约，watchdog机制。

4. 如果执行超时了，最后执行到释放锁，释放的是别人加的锁怎么办？
答: value设置为线程id，或某次请求的唯一值，释放锁的时候比较value，value一致时释放。。

5. 需要查询value，并进行比较，两步不是原子的怎么办，会有并发问题。
答： 使用lua脚本写，保证原子性。

6. 如果查询的value存在主从同步的延迟？
答：使用redlock
```

# 分布式锁的实现
##  分布式锁

一个 Go 语言实现的分布式锁库，支持多种后端（Redis、etcd、MySQL、ZooKeeper）。

### 功能特性

- 多后端支持（Redis、etcd、MySQL、ZooKeeper）
- 支持通过环境变量或 YAML 配置文件进行配置
- 自动续期（看门狗机制）
- 获取锁的重试机制
- 线程安全实现
- 简洁易用的 API

### 安装

```bash
go get github.com/yourusername/distributedlock
```

### 配置

#### 环境变量

所有配置都可以通过 `DLOCK_` 前缀的环境变量进行设置：

```bash
# Redis 配置
DLOCK_REDIS_ENABLED=true
DLOCK_REDIS_ADDRS=localhost:6379
DLOCK_REDIS_PASSWORD=yourpassword
DLOCK_REDIS_DB=0

# etcd 配置
DLOCK_ETCD_ENABLED=false
DLOCK_ETCD_ENDPOINTS=localhost:2379
DLOCK_ETCD_USERNAME=user
DLOCK_ETCD_PASSWORD=password

# MySQL 配置
DLOCK_MYSQL_ENABLED=false
DLOCK_MYSQL_USERNAME=root
DLOCK_MYSQL_PASSWORD=password
DLOCK_MYSQL_HOST=localhost
DLOCK_MYSQL_PORT=3306
DLOCK_MYSQL_DBNAME=distributed_locks

# ZooKeeper 配置
DLOCK_ZOOKEEPER_ENABLED=false
DLOCK_ZOOKEEPER_SERVERS=localhost:2181
DLOCK_ZOOKEEPER_SESSION_TIMEOUT=10s
DLOCK_ZOOKEEPER_PREFIX=/locks
```

#### 配置文件

在项目根目录或 `config` 目录下创建 `config.yaml` 文件：

```yaml
# config/config.yaml
redis:
  enabled: true
  addrs: ["localhost:6379"]
  password: ""
  db: 0

etcd:
  enabled: false
  endpoints: ["localhost:2379"]
  dial_timeout: "5s"

mysql:
  enabled: false
  username: "root"
  password: "password"
  host: "localhost"
  port: 3306
  dbname: "distributed_locks"

zookeeper:
  enabled: false
  servers: ["localhost:2181"]
  session_timeout: "10s"
  prefix: "/locks"
```

### 使用指南

#### 基础用法

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/distributedlock"
	"github.com/yourusername/distributedlock/config"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 创建一个 Redis 后端的锁
	lock := distributedlock.NewDistributedLockInfo(
		"my-resource",  // 锁的键
		"unique-value",  // 客户端的唯一标识
		30*time.Second,  // 锁的超时时间
	)

	// 设置重试选项
	lock.SetRetry(3, 100*time.Millisecond)

	// 获取锁
	ctx := context.Background()
	acquired, err := lock.AcquireLock(ctx, "redis")
	if err != nil {
		panic(fmt.Sprintf("获取锁失败: %v", err))
	}

	if !acquired {
		fmt.Println("获取锁失败")
		return
	}

	// 确保在函数退出时释放锁
	defer func() {
		err := lock.ReleaseLock(ctx, "redis")
		if err != nil {
			fmt.Printf("释放锁失败: %v\n", err)
		}
	}()

	// 临界区代码
	fmt.Println("成功获取锁，执行任务中...")
	time.Sleep(5 * time.Second)
	fmt.Println("任务完成")
}
```

#### 使用依赖注入

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/distributedlock"
)

func main() {
	// 初始化 Redis 锁服务
	redisLock, err := distributedlock.NewDistributedLock(
		distributedlock.RedisLockType,
		distributedlock.RedisConfig{
			Addrs:    []string{"localhost:6379"},
			Password: "",
			DB:       0,
		},
	)
	if err != nil {
		panic(fmt.Sprintf("创建 Redis 锁失败: %v", err))
	}

	// 注册服务
	distributedlock.RegisterService(redisLock)

	// 创建锁
	lock := distributedlock.NewDistributedLockInfo(
		"my-resource",
		"unique-value",
		30*time.Second,
	)

	// 使用锁
	ctx := context.Background()
	acquired, err := lock.AcquireLock(ctx, "redis")
	if err != nil {
		panic(fmt.Sprintf("获取锁失败: %v", err))
	}

	if !acquired {
		fmt.Println("获取锁失败")
		return
	}

	defer func() {
		err := lock.ReleaseLock(ctx, "redis")
		if err != nil {
			fmt.Printf("释放锁失败: %v\n", err)
		}
	}()

	// 临界区代码
	fmt.Println("成功获取锁，执行任务中...")
	time.Sleep(5 * time.Second)
}
```

### API 参考

#### DistributedLockInfo

```go
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
```

##### 方法

- `NewDistributedLockInfo(key, value string, expiration time.Duration) *DistributedLockInfo`
  - 使用指定的键、值和过期时间创建一个新的分布式锁。

- `SetRetry(tries int, delay time.Duration)`
  - 设置重试次数和重试间隔。

- `AcquireLock(ctx context.Context, serviceType string) (bool, error)`
  - 尝试使用指定的服务类型获取锁。

- `ReleaseLock(ctx context.Context, serviceType string) error`
  - 释放锁。

#### DistributedLockService

```go
type DistributedLockService interface {
	AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error)
	ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error)
	RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error
	BuildServiceType() string
}
```

### 后端特定说明

#### Redis
- 使用 `SET key value NX PX milliseconds` 命令实现原子性获取锁
- 锁会在过期后自动释放
- 通过看门狗机制支持锁的自动续期

#### etcd
- 使用 etcd 的分布式互斥锁实现
- 自动处理租约续期
- 需要 etcd v3 API

#### MySQL
- 使用 `GET_LOCK()` 和 `RELEASE_LOCK()` 函数实现
- 需要 MySQL 5.7.5 或更高版本
- 会话结束时锁会自动释放

#### ZooKeeper
- 使用临时节点实现锁
- 会话结束时锁会自动释放
- 支持通过前缀实现分层锁

