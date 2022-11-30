// Provides a standard interface and implementations for a kv-style store
package storage

import (
	"gt_mtc_takehome/messages"
	"sync"
	"time"
)

// Hook for external builders to set storage implementation
var (
	// Hook for external builders to set storage implementation
	StorageImpl = LocalMapStorage{
		make(map[time.Time][]messages.ArticleDayCount),
		sync.RWMutex{},
	}
	// Truncate all keys to midnight
	TRUNCATE_TO_DAY time.Duration = (24 * time.Hour)
)

// Wrapper interface for a key-value store to allow for different backends (e.g. distributed cache, db, etc...)
type Storage interface {
	Put(key string, value []messages.ArticleDayCount)
	Get(key time.Time) ([]messages.ArticleDayCount, bool)
}

// A very naive (but threadsafe!) ever growing in-memory local cache for non-prod usage.  Implements Storage interface
type LocalMapStorage struct {
	internal map[time.Time][]messages.ArticleDayCount
	rwMutex  sync.RWMutex
}

// Add an article day count
func (t *LocalMapStorage) Put(key time.Time, value []messages.ArticleDayCount) {
	key = key.Truncate(TRUNCATE_TO_DAY)
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()
	t.internal[key] = value
}

// Retrieve a pointer to an article day count. Second return value will be true if the key is present
func (t *LocalMapStorage) Get(key time.Time) ([]messages.ArticleDayCount, bool) {
	key = key.Truncate(TRUNCATE_TO_DAY)
	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()
	obj, ok := t.internal[key]
	return obj, ok
}

// non-exported helper function for testing
func (t *LocalMapStorage) size() int {
	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()
	return len(t.internal)
}
