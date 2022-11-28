package repository

import (
	"testing"
	"time"
)

func Test_Basic(t *testing.T) {
	start, _ := time.Parse("20060102", "20210101")
	end, _ := time.Parse("20060102", "20210110")
	result, _ := GetArticleCountsForDateRange(start, end)
	print(result)
}
