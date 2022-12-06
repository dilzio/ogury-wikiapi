package indexer

import (
	"github.com/stretchr/testify/assert"
	"github.com/zavitax/sortedset-go"
	"gt_mtc_takehome/constants"
	"gt_mtc_takehome/messages"
	"gt_mtc_takehome/storage"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_GetArticleCountsForDateRange_1year(t *testing.T) {
	const NUM_DAILY_ARTICLES = 1000
	mu := sync.Mutex{}
	//keeps a running count of views per article. Will use to compare with api results
	verificationMap := make(map[string]messages.ArticleCount)
	//set a clean storage impl

	DB = storage.NewLocalMapStorage()
	//set a stub fetcher which will generate some fake data
	Fetcher = func(date time.Time) ([]messages.ArticleCount, error) {
		countsSlice := make([]messages.ArticleCount, NUM_DAILY_ARTICLES)
		for i := 0; i < NUM_DAILY_ARTICLES; i++ {
			countObject := messages.ArticleCount{
				Name:  "article " + strconv.Itoa(i),
				Views: rand.Intn(10000),
			}
			countsSlice[i] = countObject
			mu.Lock()
			aggregateCount, ok := verificationMap[countObject.Name]
			if !ok {
				verificationMap[countObject.Name] = countObject
			} else {
				aggregateCount.Views = aggregateCount.Views + countObject.Views
				verificationMap[aggregateCount.Name] = aggregateCount
			}
			mu.Unlock()

		}
		return countsSlice, nil
	}

	//call the indexer and check values
	start, _ := time.Parse(constants.DATELAYOUT, "20210101")
	end, _ := time.Parse(constants.DATELAYOUT, "20220101")
	result, _ := GetArticleCountsForDateRange(start, end)
	assert.NotNil(t, result)
	assert.Equal(t, start.Year(), result.StartDate.Year())
	assert.Equal(t, start.Month(), result.StartDate.Month())
	assert.Equal(t, start.Day(), result.StartDate.Day())
	assert.Equal(t, end.Year(), result.EndDate.Year())
	assert.Equal(t, end.Month(), result.EndDate.Month())
	assert.Equal(t, end.Day(), result.EndDate.Day())
	assert.Equal(t, NUM_DAILY_ARTICLES, len(result.ArticleCounts))

	//verify counts for each article
	for _, countObject := range result.ArticleCounts {
		verificationCount, ok := verificationMap[countObject.Name]
		if !ok {
			t.Errorf("object: %s not found", countObject.Name)
		}
		assert.Equal(t, verificationCount.Views, countObject.Views)
	}
}

func Test_GetCountsForArticleInRange_1year(t *testing.T) {
	const NUM_DAILY_ARTICLES = 1000
	const TARGET_ARTICLE = "article 0"
	mu := sync.Mutex{}
	//keeps a running count of views per article. Will use to compare with api results
	verificationMap := make(map[string]messages.ArticleCount)
	//set a clean storage impl
	DB = storage.NewLocalMapStorage()
	//set a stub fetcher which will generate some fake data
	Fetcher = func(date time.Time) ([]messages.ArticleCount, error) {
		countsSlice := make([]messages.ArticleCount, NUM_DAILY_ARTICLES)
		for i := 0; i < NUM_DAILY_ARTICLES; i++ {
			countObject := messages.ArticleCount{
				Name:  "article " + strconv.Itoa(i),
				Views: rand.Intn(10000),
			}
			countsSlice[i] = countObject
			mu.Lock()
			aggregateCount, ok := verificationMap[countObject.Name]
			if !ok {
				verificationMap[countObject.Name] = countObject
			} else {
				aggregateCount.Views = aggregateCount.Views + countObject.Views
				verificationMap[aggregateCount.Name] = aggregateCount
			}
			mu.Unlock()

		}
		return countsSlice, nil
	}
	start, _ := time.Parse(constants.DATELAYOUT, "20210101")
	end, _ := time.Parse(constants.DATELAYOUT, "20220101")
	result, err := GetCountsForArticleInRange(TARGET_ARTICLE, start, end)
	if err != nil {
		print(err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, start.Year(), result.StartDate.Year())
	assert.Equal(t, start.Month(), result.StartDate.Month())
	assert.Equal(t, start.Day(), result.StartDate.Day())
	assert.Equal(t, end.Year(), result.EndDate.Year())
	assert.Equal(t, end.Month(), result.EndDate.Month())
	assert.Equal(t, end.Day(), result.EndDate.Day())
	assert.Equal(t, 1, len(result.ArticleCounts))

	//verify counts for each article
	for _, countObject := range result.ArticleCounts {
		verificationCount, ok := verificationMap[countObject.Name]
		if !ok {
			t.Errorf("object: %s not found", countObject.Name)
		}
		assert.Equal(t, verificationCount.Views, countObject.Views)
	}
}

func xTest_ssplayground(t *testing.T) {
	index := sortedset.New[string, int, messages.ArticleCount]()
	index.AddOrUpdate("article1", 900, messages.ArticleCount{
		Name:  "article1",
		Views: 100,
	})
	index.AddOrUpdate("article2", 200, messages.ArticleCount{
		Name:  "article2",
		Views: 200,
	})
	index.AddOrUpdate("article3", 300, messages.ArticleCount{
		Name:  "article1",
		Views: 400,
	})
	index.AddOrUpdate("article4", 100, messages.ArticleCount{
		Name:  "article4",
		Views: 400,
	})

	println(index.PeekMax())
	println(index.GetCount())
	println(index.GetByKey("article1").Score())
	println(index.PeekMax().Score())
	println(index.PeekMax().Key())
}
