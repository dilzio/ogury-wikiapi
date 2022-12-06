package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_E2E_API(t *testing.T) {
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

}
