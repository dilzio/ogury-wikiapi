package service

import (
	"net/http"
	"testing"
)

type stubRequest struct{}
type stubResponseWriter struct{}

func (stubResponseWriter) Header() http.Header {
	return nil
}
func (stubResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (stubResponseWriter) WriteHeader(statusCode int) {

}

func Test_DoCalcMostViewedArticles_HappyPath(t *testing.T) {
	ChiUrlParam = func(r *http.Request, key string) string {
		if key == "startdate" {
			return "20020101"
		}
		if key == "enddate" {
			return "20020102"
		} else {
			return ""
		}
	}
	underTest := Service{}
	r := http.Request{}
	w := stubResponseWriter{}
	underTest.DoCalcMostViewedArticles(w, &r)

}
