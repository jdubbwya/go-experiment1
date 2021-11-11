package metrics

import "time"

type TimeMetric struct {
	start time.Time
	end *time.Time
}

func TimeStartCapture() *TimeMetric{
	return &TimeMetric{
		start: time.Now(),
		end : nil,
	}
}

func TimeStopCapture( m *TimeMetric ){
	var end = time.Now()
	m.end = &end

	go captureTime(m)
}

func TimeElapsedDuration( m *TimeMetric ) *time.Duration {
	if m.end == nil {
		return nil
	}

	var elapsed = m.end.Sub(m.start)
	return &elapsed
}