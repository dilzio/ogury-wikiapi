package service

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

/*
*
Hook to monkey patch the chi library in tests
*/
var (
	ChiUrlParam = chi.URLParam
)

const date_layout = "20060102"

/*
Service encapsulates all the API business logic
*/
type Service struct {
}

func NewInstance() Service {
	return Service{}
}

func (s Service) DoCalcMostViewedDayInRange(w http.ResponseWriter, r *http.Request) {

}

func (s Service) DoCalcViewCountForArticle(writer http.ResponseWriter, request *http.Request) {

}

func (s Service) DoCalcMostViewedArticles(w http.ResponseWriter, r *http.Request) {

	startdate := ChiUrlParam(r, "startdate")
	enddate := ChiUrlParam(r, "enddate")

	start, err := time.Parse(date_layout, startdate)
	if err != nil {
		log.Error("Bad startdate value; ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	end, err := time.Parse(date_layout, enddate)
	if err != nil {
		log.Error("Bad enddate value: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !end.After(start) {
		log.Error("End date not after start date")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//TODO: set maximum days spread
	//DayCounts[], err = repo.GetDayCountsForDateRange(start, end)

}
