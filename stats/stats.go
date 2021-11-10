package stats

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/benchmark"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

var totalRequests int64 = 0
var totalDuration int64 = 0

var accumulateDurationChannel = make(chan int64)


func Capture( b *benchmark.Benchmark ){
	var timeElapsed *time.Duration = benchmark.TimeElapsed(*b)

	if timeElapsed != nil {
		atomic.AddInt64(&totalRequests, 1)
		atomic.AddInt64(&totalDuration, (*timeElapsed).Microseconds())
	}
}

func HandleStats(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`{"total": %d, "average": %d}`, totalRequests, totalDuration / totalRequests ) )
}
