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

}
