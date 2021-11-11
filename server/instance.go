package server

import (
	"github.com/jdubbwya/go-experiment1/stats"
	"net/http"
)

type Instance struct {
	baseUrl string
	monitor *Monitor
	server http.Server
	statsAggregator *stats.Aggregator
}

func (i *Instance) Aggregator() *stats.Aggregator {
	return i.statsAggregator
}

func (i *Instance ) IsAlive() bool  {
	select {
	case <-i.monitor.channel:
		//prevent work from being done
		return true
	default:
	}
	return false
}

func (i *Instance) BaseUrl() string {
	return i.baseUrl
}

func (i *Instance) RequestShutdown(){

}
