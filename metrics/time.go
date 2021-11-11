package metrics

import "time"

type TimeMetric struct {
	Metric
	start time.Time
	end *time.Time
}

func NewTimeMetric() *TimeMetric{
	return &TimeMetric{
		start: time.Now(),
		end : nil,
	}
}

func (t *TimeMetric) Stop(){
	var end = time.Now()
	t.end = &end

	go captureTime(t)
}

func (t *TimeMetric) ElapsedDuration() *time.Duration {
	if t.end == nil {
		return nil
	}

	var elapsed = t.end.Sub(t.start)
	return &elapsed
}