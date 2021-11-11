package server

import (
	"net/http"
	"sync"
)

type Monitor struct  {
	channel chan struct{}
	waitGroup sync.WaitGroup
}

func newMonitor() *Monitor{
	return &Monitor {
		channel: make(chan struct{}),
	}
}

func monitorHandler(handler http.Handler, monitor *Monitor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-monitor.channel:
			//prevent work from being done
			return
		default:
		}

		monitor.waitGroup.Add(1)
		handler.ServeHTTP(w, r)
		defer monitor.waitGroup.Done()
	})
}