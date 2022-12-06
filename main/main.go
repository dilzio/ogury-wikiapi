package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"gt_mtc_takehome/service"
	"net/http"
)

func main() {
	log.SetLevel(log.InfoLevel)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/mostviewed/{startdate}/{enddate}", service.DoGetArticleCountsForDateRange)
	r.Get("/viewcount/{article}/{startdate}/{enddate}", service.DoCalcViewCountForArticle)
	r.Get("/mostviewedday/{article}/{year}/{month}", service.DoCalcMostViewedDayInMonthForArticle)
	log.Infof("Hi! listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
