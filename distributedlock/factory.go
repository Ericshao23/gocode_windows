package distributedlock

import (
	"context"
	"log"
	"time"
)

func (dl *DistributedLockInfo) AcquireLock(ctx context.Context, serviceType string) (bool, error) {
	service, err := GetService(serviceType)
	if err != nil {
		return false, err
	}
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	if dl.locked {
		return true, nil
	}

	var timer *time.Timer
	var acquireLock = false
	var lockErr error = nil

	for i := range dl.failTrys {
		if i != 0 {
			// 创建或重置定时器
			if timer == nil {
				timer = time.NewTimer(dl.failDelay)
			} else {
				timer.Reset(dl.failDelay)
			}
			select {
			case <-ctx.Done():
				timer.Stop()
				return false, ctx.Err()
			case <-timer.C:
				// Fall-through when the delay timer completes.
			}
		}

		acquireLock, lockErr = service.AcquireLock(ctx, dl)
		if lockErr != nil {
			continue
		}

		if acquireLock {
			break
		}
	}

	// 确保定时器被停止
	if timer != nil {
		timer.Stop()
	}

	if lockErr != nil {
		return acquireLock, lockErr
	}
	if acquireLock {
		dl.locked = true
		go dl.startWatchdog(ctx, serviceType)
	}

	return acquireLock, nil
}

// 启动Watch Dog自动续期
func (dl *DistributedLockInfo) startWatchdog(ctx context.Context, serviceType string) {
	ticker := time.NewTicker(dl.expiration / 2) // 在过期时间的一半进行续期
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dl.mutex.Lock()
			if !dl.locked {
				dl.mutex.Unlock()
				return
			}

			err := dl.renewLock(ctx, serviceType)
			if err != nil {
				log.Printf("WatchDog: %v failed to renew lock: %v\n", dl.key, err)
				dl.locked = false
				dl.mutex.Unlock()
				return
			}
			log.Printf("WatchDog: successfully renewed lock for key: %s\n", dl.key)
			dl.mutex.Unlock()
		case <-dl.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (dl *DistributedLockInfo) renewLock(ctx context.Context, serviceType string) error {
	service, err := GetService(serviceType)
	if err != nil {
		return err
	}
	return service.RenewLock(ctx, dl)
}

func (dl *DistributedLockInfo) ReleaseLock(ctx context.Context, serviceType string) error {
	service, err := GetService(serviceType)
	if err != nil {
		return err
	}
	dl.mutex.Lock()
	defer dl.mutex.Unlock()
	if !dl.locked {
		return nil
	}
	releaseLock, err := service.ReleaseLock(ctx, dl)
	if err != nil {
		return err
	}
	if releaseLock {
		dl.locked = false
		// 安全地关闭 stopChan，避免 panic
		select {
		case <-dl.stopChan:
			// 已经关闭
		default:
			close(dl.stopChan)
		}
		log.Printf("Lock %s released successfully\n", dl.key)
	} else {
		log.Printf("Lock %s might have been released by others\n", dl.key)
	}
	return nil
}
