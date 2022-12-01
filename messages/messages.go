package messages

import "time"

// Type ArticleCount captures the counts for an article
type ArticleCount struct {
	Name  string `json:"name"`
	Views int    `json:"views"`
}

// Type ArticleCountsForDateRange wrappers a set of article counts aggregated for the days between StartDate and EndDate (inclusive of both)
type ArticleCountsForDateRange struct {
	StartDate     time.Time      ` json:"startdate"`
	EndDate       time.Time      `json:"enddate"`
	ArticleCounts []ArticleCount `json:"articles"`
}

// Type WPPageViewsPayload models the response payload of the Wikipedia Pageviews API
type WPPageViewsPayload struct {
	Items []struct {
		Project  string `json:"project"`
		Access   string `json:"access"`
		Year     string `json:"year"`
		Month    string `json:"month"`
		Day      string `json:"day"`
		Articles []struct {
			Article string `json:"article"`
			Views   int    `json:"views"`
			Rank    int    `json:"rank"`
		} `json:"articles"`
	} `json:"items"`
}
