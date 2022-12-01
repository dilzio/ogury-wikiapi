// Package service is the top-level API container and handles inbound and outbound param validation and (un)marshalling
package service

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"gt_mtc_takehome/constants"
	"gt_mtc_takehome/indexer"
	"net/http"
	"time"
)

type Service struct {
}

func NewInstance() Service {
	return Service{}
}

func (s Service) DoCalcMostViewedDayInRange(w http.ResponseWriter, r *http.Request) {

}

func (s Service) DoCalcViewCountForArticle(writer http.ResponseWriter, request *http.Request) {

}

func (s Service) DoGetArticleCountsForDateRange(w http.ResponseWriter, r *http.Request) {

	startdate := chi.URLParam(r, "startdate")
	enddate := chi.URLParam(r, "enddate")

	start, err := time.Parse(constants.DATELAYOUT, startdate)
	if err != nil {
		message := "Bad startdate value: " + err.Error()
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	}

	end, err := time.Parse(constants.DATELAYOUT, enddate)
	if err != nil {
		message := "Bad enddate value: " + err.Error()
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	}

	if !end.After(start) {
		log.Error("End date not after start date")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := indexer.GetArticleCountsForDateRange(start, end)

	var bytes []byte
	if bytes, err = json.Marshal(&result); err != nil {
		log.Println("Failed to marshal reply:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(bytes)

	//TODO: set maximum days spread
	//DayCounts[], err = repo.GetDayCountsForDateRange(start, end)

}
