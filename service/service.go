// Package service is the top-level API container and handles inbound and outbound param validation and (un)marshalling
package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"gt_mtc_takehome/constants"
	"gt_mtc_takehome/indexer"
	"net/http"
	"time"
)

func DoCalcMostViewedDayInRange(w http.ResponseWriter, r *http.Request) {

}

func DoCalcViewCountForArticle(w http.ResponseWriter, r *http.Request) {
	start, end, ok := validateDates(w, r)
	if !ok {
		return
	}
	articleName := chi.URLParam(r, "article")
	if len(articleName) == 0 {
		message := "Article name param not found: "
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	}
	result, err := indexer.GetCountsForArticleInRange(articleName, start, end)
	var bytes []byte
	if bytes, err = json.Marshal(&result); err != nil {
		log.Println("Failed to marshal reply:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// Function validateDates does basic date parsing and validation. Will return parsed start
// and end dates if successful with a true boolean or placeholders with a false boolean value if unsuccessfulX
func validateDates(w http.ResponseWriter, r *http.Request) (time.Time, time.Time, bool) {
	start, err := time.Parse(constants.DATELAYOUT, chi.URLParam(r, "startdate"))
	if err != nil {
		message := "Bad startdate value: " + err.Error()
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return time.Now(), time.Now(), false
	}

	end, err := time.Parse(constants.DATELAYOUT, chi.URLParam(r, "enddate"))
	if err != nil {
		message := "Bad enddate value: " + err.Error()
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return time.Now(), time.Now(), false
	}

	if end.Before(start) {
		message := "End date cannot be before start date"
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return time.Now(), time.Now(), false
	}

	if end.Sub(start).Hours()/24 >= constants.MAXDAYINTERVAL {
		message := fmt.Sprintf("Maximum interval between dates is: %d days ", constants.MAXDAYINTERVAL)
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return time.Now(), time.Now(), false
	}
	return start, end, true
}
func DoGetArticleCountsForDateRange(w http.ResponseWriter, r *http.Request) {
	start, end, ok := validateDates(w, r)
	if !ok {
		return
	}
	result, err := indexer.GetArticleCountsForDateRange(start, end)
	var bytes []byte
	if bytes, err = json.Marshal(&result); err != nil {
		log.Println("Failed to marshal reply:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
