package middleware

import (
	"github.com/jdubbwya/go-experiment1/metrics"
	"log"
	"net/http"
)

func StatsMetric(handler http.Handler, aggregator metrics.Aggregator) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var metric = metrics.NewTimeMetric("stats")
		log.Println("Defering metric collection")
		defer aggregator.Collect(metric)
		handler.ServeHTTP(w, r)
	})
}