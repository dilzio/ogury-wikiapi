package indexer

import (
	"github.com/stretchr/testify/assert"
	"github.com/zavitax/sortedset-go"
	"gt_mtc_takehome/messages"
	"testing"
	"time"
)

func Test_GetArticleCountsForDateRange(t *testing.T) {
	start, _ := time.Parse("20060102", "20210101")
	end, _ := time.Parse("20060102", "20210110")
	result, _ := GetArticleCountsForDateRange(start, end)
	assert.NotNil(t, result)
}

func Test_ssplayground(t *testing.T) {
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
