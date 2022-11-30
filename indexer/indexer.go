// Package indexer contains functions for sourcing and aggregating article counts
package indexer

import (
	"gt_mtc_takehome/messages"
	"gt_mtc_takehome/storage"
	"time"
)

// Type fetcher is an internal type that describes a standard function for fetching day counts from an external source
type fetcher = func(date time.Time) ([]messages.ArticleCount, error)

var (
	//Var Fetcher holds an instance of a fetcher function. It is exported to enable  stubbing for tests
	Fetcher fetcher = wikipediafetcher
	//Var DB is a cache for article day counts.  It is exported to enable stubbing for tests
	DB = storage.StorageImpl
)

// wikipediafetcher is a wrapper fetcher function for the Wikipedia Pageviews API
func wikipediafetcher(date time.Time) ([]messages.ArticleCount, error) {
	return nil, nil
}

// Get
func GetArticleCountsForDateRange(startdate time.Time, enddate time.Time) (messages.ArticleCountsForDateRange, error) {

	//declare a threadsafe sorted set
	//spawn a go function for each day
	//retrieve the slice of objects for the date first from storage then from wikipedia if necessesary
	//cache in storage if necessary
	//update itmes in sorted set from cached itmes
	//end goroutines
	//construct the return struct

	//wg := sync.WaitGroup{}
	//index := sortedset.New[string, int, messages.ArticleCount]
	for d := startdate; d.Before(enddate) == true; d = d.AddDate(0, 0, 1) {

	}

	time.Sleep(1000000)
	return messages.ArticleCountsForDateRange{}, nil
}
