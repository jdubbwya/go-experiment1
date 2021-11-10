package benchmark

import "time"

type Benchmark struct {
	start time.Time
	end *time.Time
}

func StartCapture() *Benchmark{
	return &Benchmark{
		start: time.Now(),
		end : nil,
	}
}

func StopCapture( b *Benchmark ){
	var end = time.Now()
	b.end = &end
}

func TimeElapsed( b Benchmark ) *time.Duration {
	if b.end == nil {
		return nil
	}

	var elapsed = b.end.Sub(b.start)
	return &elapsed
}