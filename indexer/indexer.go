// Package indexer contains functions for sourcing and aggregating article counts. It makes heavy use of the
// github.com/zavitax/sortedset-go implementation for ranking
package indexer

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zavitax/sortedset-go"
	"gt_mtc_takehome/constants"
	"gt_mtc_takehome/messages"
	"gt_mtc_takehome/storage"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Type fetcher is an internal type that describes a standard function for fetching day counts from an external source
type fetcher = func(date time.Time) ([]messages.ArticleCount, error)

var (
	//Var Fetcher holds an instance of a fetcher function. It is exported to enable  stubbing for tests
	Fetcher fetcher = wikipediafetcher
	//Var DB is a cache for article day counts.  It is exported to enable stubbing for tests
	DB = &storage.StorageImpl
)

// wikipediafetcher is a wrapper fetcher function for the Wikipedia Pageviews API.
func wikipediafetcher(date time.Time) ([]messages.ArticleCount, error) {
	counts := []messages.ArticleCount{}
	year := strconv.Itoa(date.Year())
	month := date.Format(constants.TWODAYMONTH)
	day := date.Format(constants.TWODAYDAYOFWEEK)
	url := fmt.Sprintf(constants.PAGEVIEWS_URL, year, month, day)

	//Call the API and return an error if any
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error("Error Calling Wikipedia API. Status: ", resp.StatusCode, err)
		return counts, err
	}
	//Map body into struct representation, return an error if it fails
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	responseStruct := messages.WPPageViewsPayload{}
	err2 := json.Unmarshal(body, &responseStruct)
	if err2 != nil {
		return []messages.ArticleCount{}, err2
	}

	//Finally map into our internal entity representation
	articles := responseStruct.Items[0].Articles
	for _, article := range articles {
		counts = append(counts, messages.ArticleCount{
			Name:  article.Article,
			Views: article.Views,
		})
	}
	return counts, nil
}

// Function GetArticleCountsForDateRange concurrently fetches and assembles article counts for a date range
func GetArticleCountsForDateRange(startdate time.Time, enddate time.Time) (messages.ArticleCountsForDateRange, error) {

	//declare a threadsafe sorted set
	//spawn a go function for each day
	//retrieve the slice of objects for the date first from storage then from wikipedia if necessesary
	//cache in storage if necessary
	//update itmes in sorted set from cached itmes
	//end goroutines
	//construct the return struct

	wg := sync.WaitGroup{}
	index := sortedset.New[string, int, messages.ArticleCount]()
	ssUpdateMutex := sync.Mutex{}
	for d := startdate; !d.After(enddate) == true; d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()
			//defer ssUpdateMutex.Unlock()
			countsForDay := getArticleCountsForDay(date)
			for _, countobject := range countsForDay {
				ssUpdateMutex.Lock()
				node := index.GetByKey(countobject.Name)
				if node == nil {
					index.AddOrUpdate(countobject.Name, countobject.Views, countobject)
				} else {
					aggregateCountObj := node.Value
					aggregateCountObj.Views = aggregateCountObj.Views + countobject.Views
					//ssUpdateMutex.Lock()
					index.AddOrUpdate(countobject.Name, aggregateCountObj.Views, aggregateCountObj)
				}
				ssUpdateMutex.Unlock()
			}
		}(d)
	}

	wg.Wait()
	allTheRankedNodes := index.GetRangeByRank(-1, 1, false)
	payload := messages.ArticleCountsForDateRange{}
	payload.StartDate = startdate
	payload.EndDate = enddate
	for _, node := range allTheRankedNodes {
		payload.ArticleCounts = append(payload.ArticleCounts, node.Value)
	}

	return payload, nil
}

// Function getArticleCountsForDay will check the db cache for the slice of article counts and if not found will
// pull from the Wikipedia api
func getArticleCountsForDay(day time.Time) []messages.ArticleCount {
	cachedcounts, ok := DB.Get(day)
	if !ok {
		fetchedCounts, err := Fetcher(day)
		if err != nil {
			log.Errorf("Unable to retrieve counts for date: %s", day.String())
			return nil
		}
		DB.Put(day, fetchedCounts)
		return fetchedCounts
	}
	return cachedcounts
}
