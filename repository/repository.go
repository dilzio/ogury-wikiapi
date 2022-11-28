package repository

import (
	"fmt"
	"gt_mtc_takehome/messages"
	"gt_mtc_takehome/storage"
	"time"
)

var (
	//Fetches batches from Wikipedia
	Fetcher = wikipediafetcher
	//The underlying storage implementation
	DB = storage.StorageImpl
)

func wikipediafetcher(date time.Time) ([]messages.ArticleDayCount, error) {
	return nil, nil
}

func GetArticleCountsForDateRange(startdate time.Time, enddate time.Time) (map[time.Time]*[]messages.ArticleDayCount, error) {

	for d := startdate; d.After(enddate) == false; d = d.AddDate(0, 0, 1) {
		go func(dd time.Time) { fmt.Println(dd.Format("2006-01-02")) }(d)
	}

	time.Sleep(1000000)
	return nil, nil
}
