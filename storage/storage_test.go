package storage

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"pelotechfun/messages"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Load many items concurrently then pull them out and check them
func Test_LocalMapStorage(t *testing.T) {
	underTest := NewLocalMapStorage()
	wg := sync.WaitGroup{}
	now := time.Now()
	future := now.AddDate(00, 0, 10000)
	var dateMap = map[time.Time][]messages.ArticleCount{}
	for d := now; d.Before(future) == true; d = d.AddDate(0, 0, 1) {
		wg.Add(1)

		payload := make([]messages.ArticleCount, 1000)
		for index, countobj := range payload {
			countobj.Name = d.String() + strconv.Itoa(index)
			countobj.Views = rand.Intn(100)
		}
		dateMap[d] = payload
		go func(key time.Time) {
			defer wg.Done()
			underTest.Put(key, payload)
		}(d)
	}
	wg.Wait()

	//spin through the map of sent keys and verify there is an object there for it and that the views match
	for datekey := range dateMap {
		_, found := underTest.Get(datekey)
		assert.True(t, found)
		delete(dateMap, datekey)
	}
	assert.Equal(t, 0, len(dateMap))
}
