package metrics

import "time"

type TimeMetric struct {
	NamedMetric
	start time.Time
	end *time.Time
}

func NewTimeMetric(name string) TimeMetric{
	namedMetric := NamedMetric{
		name: name,
	}
	return TimeMetric{
		NamedMetric: namedMetric,
		start: time.Now(),
		end : nil,
	}
}

func (t *TimeMetric) Stop(){
	var end = time.Now()
	t.end = &end

	Capture(t)
}

func (t *TimeMetric) ElapsedDuration() *time.Duration {
	if t.end == nil {
		return nil
	}

	var elapsed = t.end.Sub(t.start)
	return &elapsed
}