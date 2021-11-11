package handlers

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/stats"
	"io"
	"net/http"
)

func StatsHandler(aggregator *stats.Aggregator) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, fmt.Sprintf(`{"total": %d, "average": %d}`, aggregator.TotalRequests(), aggregator.AverageRequestDuration().Microseconds()))
	})
}