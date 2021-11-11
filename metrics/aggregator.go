package metrics

type Aggregator interface {
	Collect(m Metric)
	Stop()
}
