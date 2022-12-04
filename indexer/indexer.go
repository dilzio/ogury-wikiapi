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

// Function GetArticleCountsForDateRange concurrently fetches and assembles a view ranking of all articles in a date range
func GetArticleCountsForDateRange(startdate time.Time, enddate time.Time) (messages.ArticleCountsForDateRange, error) {
	wg := sync.WaitGroup{}
	index := sortedset.New[string, int, messages.ArticleCount]()
	ssUpdateMutex := sync.Mutex{}
	for d := startdate; !d.After(enddate) == true; d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()
			countsForDay := getArticleCountsForDay(date)
			for _, countobject := range countsForDay {
				ssUpdateMutex.Lock()
				node := index.GetByKey(countobject.Name)
				if node == nil {
					index.AddOrUpdate(countobject.Name, countobject.Views, countobject)
				} else {
					aggregateCountObj := node.Value
					aggregateCountObj.Views = aggregateCountObj.Views + countobject.Views
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

// Function GetCountsForArticleInRange assembles a total view count for q specific article in a date range
func GetCountsForArticleInRange(article string, startdate time.Time, enddate time.Time) (messages.ArticleCountsForDateRange, error) {
	wg := sync.WaitGroup{}
	index := sortedset.New[string, int, messages.ArticleCount]()
	ssUpdateMutex := sync.Mutex{}
	for d := startdate; !d.After(enddate) == true; d = d.AddDate(0, 0, 1) {
		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()
			countsForDay := getArticleCountsForDay(date)
			for _, countobject := range countsForDay {
				if countobject.Name == article {
					ssUpdateMutex.Lock()
					node := index.GetByKey(countobject.Name)
					if node == nil {
						index.AddOrUpdate(countobject.Name, countobject.Views, countobject)
					} else {
						aggregateCountObj := node.Value
						aggregateCountObj.Views = aggregateCountObj.Views + countobject.Views
						index.AddOrUpdate(countobject.Name, aggregateCountObj.Views, aggregateCountObj)
					}
					ssUpdateMutex.Unlock()
				}

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

// Function GetTopDayForArticle returns the most viewed day for an article in the time range
func GetTopDayForArticle(article string, startdate time.Time, enddate time.Time) (messages.ArticleCountsForDateRange, error) {
	wg := sync.WaitGroup{}
	index := sortedset.New[string, int, messages.ArticleCount]()
	ssUpdateMutex := sync.Mutex{}
	for d := startdate; d.Before(enddate) == true; d = d.AddDate(0, 0, 1) {

		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()
			countsForDay := getArticleCountsForDay(date)
			for _, countobject := range countsForDay {
				if countobject.Name == article {
					log.Debugf("count object date: %s  views: %d", date.String(), countobject.Views)
					ssUpdateMutex.Lock()
					node := index.GetByKey(countobject.Name)
					if node == nil {
						countobject.Date = date
						index.AddOrUpdate(countobject.Name, countobject.Views, countobject)
					} else {
						aggregateCountObj := node.Value
						if countobject.Views > aggregateCountObj.Views {
							countobject.Date = date
							//println("count object d is:", date.String())
							index.AddOrUpdate(countobject.Name, countobject.Views, countobject)
						}
					}
					ssUpdateMutex.Unlock()
				}

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
