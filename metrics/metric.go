package metrics

var capture = make(chan Metric)

type Metric interface {
	Name() string
}

func Capture( m Metric ){
	go func() {
		capture <- m
	}()
}