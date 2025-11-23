package distributedlock

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZookeeperLock struct {
	conn   *zk.Conn
	acl    []zk.ACL
	prefix string
}

func NewZookeeperLock(servers []string, sessionTimeout time.Duration, prefix string) (*ZookeeperLock, error) {
	conn, _, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return nil, err
	}

	// Default ACL gives all permissions to anyone (not suitable for production)
	acl := zk.WorldACL(zk.PermAll)

	return &ZookeeperLock{
		conn:   conn,
		acl:    acl,
		prefix: strings.TrimRight(prefix, "/") + "/",
	}, nil
}

func (z *ZookeeperLock) AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	// Create the lock node
	path := z.prefix + lockInfo.key

	// Create parent nodes if they don't exist
	if err := z.ensurePath(path); err != nil {
		return false, err
	}

	// Create an ephemeral sequential node
	_, err := z.conn.Create(path, []byte(lockInfo.value), zk.FlagEphemeral, z.acl)
	if err != nil {
		if err == zk.ErrNodeExists {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (z *ZookeeperLock) ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	path := z.prefix + lockInfo.key

	// Check if the node exists
	exists, _, err := z.conn.Exists(path)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	// Delete the node to release the lock
	err = z.conn.Delete(path, -1)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (z *ZookeeperLock) RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error {
	path := z.publicPath(lockInfo.key)

	// Check if the lock still exists
	exists, _, err := z.conn.Exists(path)
	if err != nil {
		return err
	}

	if !exists {
		return ErrLockNotHeld
	}

	// In ZooKeeper, ephemeral nodes are automatically removed when the session ends
	// So we just need to check if the node still exists
	return nil
}

func (z *ZookeeperLock) BuildServiceType() string {
	return "zookeeper"
}

// ensurePath creates all nodes in the path if they don't exist
func (z *ZookeeperLock) ensurePath(fullPath string) error {
	nodes := strings.Split(strings.Trim(fullPath, "/"), "/")
	currentPath := ""

	for _, node := range nodes {
		currentPath = path.Join(currentPath, node)
		_, err := z.conn.Create(currentPath, []byte{}, 0, z.acl)
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}

	return nil
}

// publicPath returns the full path for a lock
func (z *ZookeeperLock) publicPath(lockName string) string {
	return path.Join(z.prefix, lockName)
}

// Close closes the ZooKeeper connection
func (z *ZookeeperLock) Close() {
	z.conn.Close()
}
