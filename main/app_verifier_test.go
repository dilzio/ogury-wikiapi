package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

//Test_E2E_API is quick and dirty full E2E blackbox integration test program that can be run against
//either live or stubbed out dependencies.  It can be used in a few different ways:
//1. During development by engineers as a smoke test
//2. As part of a CD/CI pipeline with stubbed data
//3. As a cronned health check against a live production instance

func Test_E2E_API(t *testing.T) {
	/*
		Example code to show how Wikipedia fetcher can be stubbed out
		indexer.Fetcher = func(date time.Time) ([]messages.ArticleCount, error) {
			// returns stubbed data
		}
	*/
	go main()
	//TestDoCalcMostViewedArticles HappyPath
	r, _ := http.Get("http://localhost:8080/mostviewed/20220101/20220102")
	defer r.Body.Close()
	bytes, _ := io.ReadAll(r.Body)
	payloadString := string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "\"Main_Page\",\"views\":10226718"))

	//TestDoCalcMostViewedArticles bad dates
	r, _ = http.Get("http://localhost:8080/mostviewed/20220101/")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "404 page not found\n"))

	//TestDoCalcMostViewedArticles data not found
	r, _ = http.Get("http://localhost:8080/mostviewed/20010101/20010102")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "Unable to retrieve page count data from Wikipedia: 20010101"))
	assert.True(t, strings.Contains(payloadString, "Unable to retrieve page count data from Wikipedia: 20010102"))

	//TestDoCalcMostViewedArticles more than maximum duration
	r, _ = http.Get("http://localhost:8080/mostviewed/20210101/20220101")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "Maximum interval between dates is: 100 days"))

	//TestDoCalcMostViewedArticles end date before after startdate
	r, _ = http.Get("http://localhost:8080/mostviewed/20210101/20200101")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "End date cannot be before start date"))

	//Test DoCalcMostViewedDayInMonthForArticle - happy path
	r, _ = http.Get("http://localhost:8080/mostviewedday/Albert_Einstein/2015/07")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "{\"startdate\":\"2015-07-01T00:00:00Z\",\"enddate\":\"2015-08-01T00:00:00Z\",\"articles\":[{\"name\":\"Albert_Einstein\",\"views\":17269,\"time\":\"2015-07-23T00:00:00Z\"}]}"))

	//Test DoCalcMostViewedDayInMonthForArticle - bad month
	r, _ = http.Get("http://localhost:8080/mostviewedday/Albert_Einstein/2015/14")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "Bad date params.  Format should be 4-digit year and 2 digit month eg: /mostviewedday/myarticle/2022/01"))

	//Test DoCalcMostViewedDayInMonthForArticle - bad year
	r, _ = http.Get("http://localhost:8080/mostviewedday/Albert_Einstein/201/01")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "Bad date params.  Format should be 4-digit year and 2 digit month eg: /mostviewedday/myarticle/2022/01"))

	//Test viewcount - happy path
	r, _ = http.Get("http://localhost:8080/viewcount/Dua_Lipa/20210101/20210103")
	defer r.Body.Close()
	bytes, _ = io.ReadAll(r.Body)
	payloadString = string(bytes[:])
	assert.True(t, strings.Contains(payloadString, "{\"startdate\":\"2021-01-01T00:00:00Z\",\"enddate\":\"2021-01-03T00:00:00Z\",\"articles\":[{\"name\":\"Dua_Lipa\",\"views\":96635,\"time\":\"0001-01-01T00:00:00Z\"}]}"))

}
