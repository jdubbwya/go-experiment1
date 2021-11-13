package stats

import (
	"context"
	"github.com/jdubbwya/go-experiment1/metrics"
	"sync"
	"time"
)

type monitor struct {
	ctx  context.Context
	stop context.CancelFunc
}

func newMonitor(a *Aggregator) monitor {
	ctx, done := context.WithCancel(context.Background())

	// start a monitor routine to capture the metrics
	go func(aggregator *Aggregator) {
	monitorLoop:
		for {
			select {
			case metric, ok := <-aggregator.metrics:
				if !ok {
					return
				}
				aggregator.totalRequests++

				metricElapsedDuration := metric.ElapsedDuration()
				if metricElapsedDuration != nil {
					aggregator.totalProcessingTime = aggregator.totalProcessingTime + metricElapsedDuration.Nanoseconds()
				}
			case <-ctx.Done():
				break monitorLoop
			}
		}
	}(a)

	return monitor{
		ctx:  ctx,
		stop: done,
	}
}

type Aggregator struct {
	metrics.Aggregator

	metrics chan metrics.TimeMetric
	monitorMu sync.Mutex
	monitor *monitor

	totalRequests       int64
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

func (a Aggregator) Collect(m metrics.Metric) {
	if m.Name() == "stats" {
		t := m.(metrics.TimeMetric)
		t.Stop()
		// publish the metrics to the channel without blocking current routine
		go func() {
			a.metrics <- t
		}()
	}
}

func (a Aggregator) Stop() {
	a.monitorMu.Lock()
	if a.monitor != nil {
		a.monitor.stop()
		a.monitor = nil
	}
	a.monitorMu.Unlock()
}

func NewStatsAggregator() *Aggregator {
	aggregator := Aggregator{
		metrics: make(chan metrics.TimeMetric),
	}

	mon := newMonitor(&aggregator)
	aggregator.monitor = &mon

	return &aggregator
}
