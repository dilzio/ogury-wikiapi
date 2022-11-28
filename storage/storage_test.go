package storage

import (
	"github.com/stretchr/testify/assert"
	"gt_mtc_takehome/messages"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Load many items concurrently then pull them out and check them
func Test_LocalMapStorage(t *testing.T) {
	wg := sync.WaitGroup{}
	now := time.Now()
	future := now.AddDate(00, 0, 20000)
	var dateMap = map[time.Time]messages.ArticleDayCount{}
	for d := now; d.Before(future) == true; d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		payload := messages.ArticleDayCount{d.String(), int64(rand.Intn(100))}
		dateMap[d] = payload
		go func(key time.Time) {
			defer wg.Done()
			StorageImpl.Put(key, payload)
		}(d)
	}
	wg.Wait()

	//spin through the map of sent keys and verify there is an object there for it and that the views match
	for datekey := range dateMap {
		obj, found := StorageImpl.Get(datekey)
		assert.True(t, found)
		assert.Equal(t, dateMap[datekey].Views, obj.Views)
		delete(dateMap, datekey)
	}
	assert.Equal(t, 0, len(dateMap))
}
