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

//Function DoCalcMostViewedDayInMonthForArticle returns the day in a specified month
//when the article had the most views

func DoCalcMostViewedDayInMonthForArticle(w http.ResponseWriter, r *http.Request) {
	articleName, articleok := validateArticleParam(w, r)
	if !articleok {
		return
	}
	yearstr := chi.URLParam(r, "year")
	monthstr := chi.URLParam(r, "month")
	firstOfTheMonth, err := time.Parse("20060102", yearstr+monthstr+"01")
	if err != nil {
		message := "Bad date params.  Format should be 4-digit year and 2 digit month eg: /mostviewedday/myarticle/2022/01"
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	}

	onemonthlater := firstOfTheMonth.AddDate(0, 1, 0)
	firstOfNextMonth := time.Date(onemonthlater.Year(), onemonthlater.Month(), 1, 0, 0, 0, 0, onemonthlater.Location())
	result, err := indexer.GetTopDayForArticle(articleName, firstOfTheMonth, firstOfNextMonth)
	var bytes []byte
	if bytes, err = json.Marshal(&result); err != nil {
		log.Error("Failed to marshal reply:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// Function DoCalcViewCountForArticle will return the view count for an article in a date range
func DoCalcViewCountForArticle(w http.ResponseWriter, r *http.Request) {
	start, end, ok := validateDates(w, r)
	if !ok {
		return
	}
	articleName, articleok := validateArticleParam(w, r)
	if !articleok {
		return
	}
	result, err := indexer.GetCountsForArticleInRange(articleName, start, end)
	var bytes []byte
	if bytes, err = json.Marshal(&result); err != nil {
		log.Error("Failed to marshal reply:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// Function validateArticleParam checks for the presence of an article.  Strictly speaking it isn't needed with the current
// rounting setup as if the argument is missing the middleware will catch it, but it's here for completeness if routing were to change.
func validateArticleParam(w http.ResponseWriter, r *http.Request) (string, bool) {
	articleName := chi.URLParam(r, "article")
	if len(articleName) == 0 {
		message := "Article name param not found: "
		log.Error(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return "", false
	}
	return articleName, true
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
		log.Error("Failed to marshal reply:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
