package metrics

import (
	"sync/atomic"
	"time"
)


type timingReport struct {
	TotalRequests int64
	TotalMicroseconds int64
}

var TimeReport = timingReport{
	TotalRequests : 0,
	TotalMicroseconds : 0,
}

func captureTime( m *TimeMetric ){
	var timeElapsed *time.Duration = TimeElapsedDuration(m)

	if timeElapsed != nil {
		atomic.AddInt64(&TimeReport.TotalRequests, 1)
		atomic.AddInt64(&TimeReport.TotalMicroseconds, (*timeElapsed).Microseconds())
	}
}