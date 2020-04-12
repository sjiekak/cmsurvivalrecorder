package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type httpHandler func(http.ResponseWriter, *http.Request)

type handler struct {
	db *sql.DB
}

type timeValue struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

func middleware(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var responseHeaders = map[string]string{
			"Content-Type": "application/json",
		}
		for key, value := range responseHeaders {
			w.Header().Set(key, value)
		}
		logger := log.WithFields(log.Fields{
			"endpoint": r.RequestURI,
			"method":   r.Method,
		})
		logger.Info("handling new request")
		inner(w, r)
		logger.Info("request completed")
	})
}

func (h *handler) lastValue(response http.ResponseWriter, request *http.Request) {
	rows, err := h.db.Query("select * FROM raised order by time DESC limit 1")
	if err != nil {
		http.Error(response, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r timeValue
		if err = rows.Scan(&r.Time, &r.Value); err != nil {
			http.Error(response, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(response).Encode(&r)
		return
	}
	http.Error(response, `{"error": "no data"}`, http.StatusInternalServerError)
}

func (h *handler) allValues(response http.ResponseWriter, request *http.Request) {
	rows, err := h.db.Query("select * FROM raised")
	if err != nil {
		http.Error(response, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	timeseries := map[time.Time]float64{}
	for rows.Next() {
		var r timeValue
		if err = rows.Scan(&r.Time, &r.Value); err != nil {
			continue
		}
		timeseries[r.Time] = r.Value
	}
	json.NewEncoder(response).Encode(timeseries)
}

func lastValueHandler(db *sql.DB) http.HandlerFunc {
	h := handler{db}
	return middleware(h.lastValue)
}

func allValuesHandler(db *sql.DB) http.HandlerFunc {
	h := handler{db}
	return middleware(h.allValues)
}
