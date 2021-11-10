package stats

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/benchmark"
	"io"
	"net/http"
	"sync"
	"time"
)

var capturedData sync.Map

var countChannel  = make(chan bool)
var count int64 = 0

var durationChannel = make(chan bool)
var averageDuration int64 = 0


func Capture( b *benchmark.Benchmark ){
	var timeElapsed *time.Duration = benchmark.TimeElapsed(*b)

	if timeElapsed != nil {
		countChannel <- true
		count++
		<- countChannel

		capturedData.Store(count, timeElapsed.Microseconds())

		durationChannel <- true
		var totalDuration int64 = 0
		var i int64
		for i = 0; i < count; i++ {
			var datum, _ = capturedData.Load(i)
			totalDuration = totalDuration + datum.(int64)
		}

		averageDuration = totalDuration / count

		<- durationChannel
	}
}

func HandleStats(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`{"total": %d, "average": %d}`, count, averageDuration) )
}
