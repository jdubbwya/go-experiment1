package metrics

type NamedMetric struct {
	Metric
	name string
}

func (m NamedMetric) Name() string{
	return m.name
}
func NewMetric(name string) *NamedMetric{
	return &NamedMetric{
		name: name,
	}
}