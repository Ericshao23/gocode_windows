package distributedlock

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLLock struct {
	db *sql.DB
}

func NewMySQLLock(dsn string) (*MySQLLock, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &MySQLLock{db: db}, nil
}

func (m *MySQLLock) AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	// Use MySQL's GET_LOCK function
	var result int
	err := m.db.QueryRowContext(ctx, "SELECT GET_LOCK(?, ?)", lockInfo.key, int(lockInfo.expiration.Seconds())).Scan(&result)
	if err != nil {
		return false, err
	}

	// GET_LOCK returns 1 if the lock was obtained, 0 if it timed out
	return result == 1, nil
}

func (m *MySQLLock) ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	// Use MySQL's RELEASE_LOCK function
	var result int
	err := m.db.QueryRowContext(ctx, "SELECT RELEASE_LOCK(?)", lockInfo.key).Scan(&result)
	if err != nil {
		return false, err
	}

	// RELEASE_LOCK returns 1 if the lock was released, 0 if the lock wasn't held by this thread, or NULL if the lock didn't exist
	return result == 1, nil
}

func (m *MySQLLock) RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error {
	// MySQL's GET_LOCK doesn't support renewing the lock directly
	// We need to release and re-acquire the lock
	released, err := m.ReleaseLock(ctx, lockInfo)
	if err != nil {
		return err
	}
	if !released {
		return ErrLockNotHeld
	}

	acquired, err := m.AcquireLock(ctx, lockInfo)
	if err != nil {
		return err
	}
	if !acquired {
		return ErrLockNotAcquired
	}

	return nil
}

func (m *MySQLLock) BuildServiceType() string {
	return "mysql"
}

// Close closes the database connection
func (m *MySQLLock) Close() error {
	return m.db.Close()
}
