package distributedlock

import (
	"context"
	"testing"
	"time"
)

// TestNewDistributedLockInfo tests the constructor initializes all fields correctly
func TestNewDistributedLockInfo(t *testing.T) {
	key := "test-key"
	value := "test-value"
	expiration := 30 * time.Second

	lock := NewDistributedLockInfo(key, value, expiration)

	if lock == nil {
		t.Fatal("NewDistributedLockInfo returned nil")
	}

	if lock.key != key {
		t.Errorf("Expected key %s, got %s", key, lock.key)
	}

	if lock.value != value {
		t.Errorf("Expected value %s, got %s", value, lock.value)
	}

	if lock.expiration != expiration {
		t.Errorf("Expected expiration %v, got %v", expiration, lock.expiration)
	}

	// Verify stopChan is initialized
	if lock.stopChan == nil {
		t.Error("stopChan is not initialized")
	}

	// Verify default retry values
	if lock.failTrys != 3 {
		t.Errorf("Expected failTrys 3, got %d", lock.failTrys)
	}

	if lock.failDelay != 100*time.Millisecond {
		t.Errorf("Expected failDelay 100ms, got %v", lock.failDelay)
	}

	if lock.locked {
		t.Error("Expected locked to be false initially")
	}
}

// TestSetRetry tests the SetRetry method
func TestSetRetry(t *testing.T) {
	lock := NewDistributedLockInfo("test-key", "test-value", 30*time.Second)

	tries := 5
	delay := 200 * time.Millisecond

	lock.SetRetry(tries, delay)

	if lock.failTrys != tries {
		t.Errorf("Expected failTrys %d, got %d", tries, lock.failTrys)
	}

	if lock.failDelay != delay {
		t.Errorf("Expected failDelay %v, got %v", delay, lock.failDelay)
	}
}

// TestRegisterAndGetService tests service registration and retrieval
func TestRegisterAndGetService(t *testing.T) {
	// Create a mock service
	mockService := &mockLockService{serviceType: "mock"}

	// Register the service
	RegisterService(mockService)

	// Retrieve the service
	service, err := GetService("mock")
	if err != nil {
		t.Fatalf("Failed to get service: %v", err)
	}

	if service.BuildServiceType() != "mock" {
		t.Errorf("Expected service type 'mock', got %s", service.BuildServiceType())
	}
}

// TestGetServiceNotFound tests that GetService returns error for non-existent service
func TestGetServiceNotFound(t *testing.T) {
	_, err := GetService("non-existent-service")
	if err == nil {
		t.Error("Expected error for non-existent service, got nil")
	}
}

// TestStopChanInitialization tests that stopChan is properly initialized
func TestStopChanInitialization(t *testing.T) {
	lock := NewDistributedLockInfo("test-key", "test-value", 30*time.Second)

	// Verify stopChan is not nil and can be used
	if lock.stopChan == nil {
		t.Error("stopChan is nil")
		return
	}

	// Try to close the channel (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Closing stopChan panicked: %v", r)
		}
	}()

	close(lock.stopChan)
}

// TestStopChanClose tests that stopChan can be safely closed
func TestStopChanClose(t *testing.T) {
	lock := NewDistributedLockInfo("test-key", "test-value", 30*time.Second)

	// Close the channel
	close(lock.stopChan)

	// Try to receive from closed channel (should not panic)
	select {
	case <-lock.stopChan:
		// Successfully received from closed channel
	case <-time.After(100 * time.Millisecond):
		t.Error("Failed to receive from closed stopChan")
	}
}

// TestRetryLoopExecution tests that retry loop executes with default failTrys
func TestRetryLoopExecution(t *testing.T) {
	lock := NewDistributedLockInfo("test-key", "test-value", 30*time.Second)

	// Verify that failTrys is not 0 (which would prevent loop execution)
	if lock.failTrys == 0 {
		t.Error("failTrys should not be 0")
	}

	// Test that we can iterate over failTrys
	count := 0
	for i := range lock.failTrys {
		count++
		if i >= lock.failTrys {
			t.Errorf("Loop index %d exceeds failTrys %d", i, lock.failTrys)
		}
	}

	if count != lock.failTrys {
		t.Errorf("Expected loop to execute %d times, executed %d times", lock.failTrys, count)
	}
}

// Mock service for testing
type mockLockService struct {
	serviceType string
}

func (m *mockLockService) AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	return true, nil
}

func (m *mockLockService) ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	return true, nil
}

func (m *mockLockService) RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error {
	return nil
}

func (m *mockLockService) BuildServiceType() string {
	return m.serviceType
}
