package messages

import "time"

// Type ArticleCount captures the counts for an article
type ArticleCount struct {
	Name  string `json:"name"`
	Views int64  `json:"views"`
}

// Type ArticleCountsForDateRange wrappers a set of article counts aggregated for the days between StartDate and EndDate (inclusive of both)
type ArticleCountsForDateRange struct {
	StartDate     time.Time      ` json:"startdate"`
	EndDate       time.Time      `json:"enddate"`
	ArticleCounts []ArticleCount `json:"articles"`
}
