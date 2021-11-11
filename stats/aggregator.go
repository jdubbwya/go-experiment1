package stats

import (
	"context"
	"github.com/jdubbwya/go-experiment1/metrics"
	"time"
)

type Aggregator struct {
	metrics.Aggregator

	metrics chan metrics.TimeMetric
	monitorCancel context.CancelFunc

	totalRequests int64
	totalProcessingTime int64

}

func (a *Aggregator) TotalRequests() int64 {
	return a.totalRequests
}

func (a *Aggregator) AverageRequestDuration() time.Duration {
	if a.totalRequests > 1 && a.totalProcessingTime > 1 {
		return time.Duration(a.totalProcessingTime / a.totalRequests)
	}
	return 0
}

func (a Aggregator) Collect(m metrics.Metric)  {
	if m.Name() == "stats" {
		t := m.(metrics.TimeMetric)
		t.Stop()
		// publish the metrics to the channel without blocking current routine
		go func() {
			a.metrics <- t
		}()
	}
}

func (a Aggregator) Stop()  {
	a.monitorCancel()
}

func NewStatsAggregator() *Aggregator {
	ctx, ctxCancel := context.WithCancel(context.Background())

	aggregator := Aggregator{
		metrics: make(chan metrics.TimeMetric),
		monitorCancel: ctxCancel,
	}

	// start a monitor routine to capture the metrics
	go func(aggregator *Aggregator) {
		for {
			select {
			case metric, ok := <-aggregator.metrics:
				if !ok {
					return
				}
				aggregator.totalRequests++

				metricElapsedDuration := metric.ElapsedDuration()
				if metricElapsedDuration != nil {
					aggregator.totalProcessingTime = aggregator.totalProcessingTime+ metricElapsedDuration.Nanoseconds()
				}
			case <-ctx.Done():
				break
			}
		}
	}(&aggregator)

	return &aggregator
}