// Provides a standard interface and implementations for a kv-style store for lists of article counts
package storage

import (
	"gt_mtc_takehome/messages"
	"sync"
	"time"
)

var (
	// Truncate all keys to midnight
	TRUNCATE_TO_DAY time.Duration = (24 * time.Hour)
)

// Wrapper interface for a key-value store to allow for different backends (e.g. distributed cache, db, etc...)
type Storage interface {
	Put(key time.Time, value []messages.ArticleCount)
	Get(key time.Time) ([]messages.ArticleCount, bool)
}

// A very naive (but threadsafe!) ever growing in-memory local cache for non-prod usage.  Implements Storage interface
type LocalMapStorage struct {
	internal map[time.Time][]messages.ArticleCount
	rwMutex  sync.RWMutex
}

// factory for a LocalMapStorage instance
func NewLocalMapStorage() *LocalMapStorage {
	return &LocalMapStorage{
		make(map[time.Time][]messages.ArticleCount),
		sync.RWMutex{},
	}
}

// Add an article day count
func (t *LocalMapStorage) Put(key time.Time, value []messages.ArticleCount) {
	key = key.Truncate(TRUNCATE_TO_DAY)
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()
	t.internal[key] = value
}

// Retrieve a pointer to an article day count. Second return value will be true if the key is present
func (t *LocalMapStorage) Get(key time.Time) ([]messages.ArticleCount, bool) {
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
