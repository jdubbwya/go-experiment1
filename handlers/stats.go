package handlers

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/metrics"
	"io"
	"net/http"
)

func StatsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var averageDurationInMicroseconds int64 = 0
	if metrics.TimeReport.TotalRequests > 0 {
		averageDurationInMicroseconds = metrics.TimeReport.TotalMicroseconds / metrics.TimeReport.TotalRequests
	}
	io.WriteString(w, fmt.Sprintf(`{"total": %d, "average": %d}`, metrics.TimeReport.TotalRequests, averageDurationInMicroseconds ) )
}
