package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gt_mtc_takehome/service"
	"log"
	"net/http"
)

func main() {
	s := service.NewInstance()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/mostviewed/{startdate}/{enddate}", s.DoGetArticleCountsForDateRange)
	r.Get("/viewcount/{article}/{startdate}/{enddate}", s.DoCalcViewCountForArticle)
	r.Get("/mostviewedday/{article}/{year}/{month}", s.DoCalcMostViewedDayInRange)
	log.Fatal(http.ListenAndServe(":8080", r))
}
