package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gt_mtc_takehome/service"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/mostviewed/{startdate}/{enddate}", service.DoGetArticleCountsForDateRange)
	r.Get("/viewcount/{article}/{startdate}/{enddate}", service.DoCalcViewCountForArticle)
	r.Get("/mostviewedday/{article}/{year}/{month}", service.DoCalcMostViewedDayInRange)
	log.Fatal(http.ListenAndServe(":8080", r))
}
