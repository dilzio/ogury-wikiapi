package main

import (
	"net/http"
	"testing"
)

func Test_DoCalcMostViewedArticles_HappyPath(t *testing.T) {
	go main()
	r, err := http.Get("http://localhost:8080/mostviewed/20220101/20220120")
	print(r, err)
	r, err = http.Get("http://localhost:8080/mostviewed/20220101/")
	print(r, err)
	r, err = http.Get("http://localhost:8080/mostviewed/20220101/20210120")
	print(r, err)
}
