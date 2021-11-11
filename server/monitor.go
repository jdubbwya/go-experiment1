package server

import (
	"net/http"
	"sync"
)

// As long as this channel is open shutdown has not been called
type monitor struct  {
	Channel chan struct{}
	WaitGroup sync.WaitGroup
}

var Monitor = monitor {
	Channel: make(chan struct{}),
}

func MonitorHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-Monitor.Channel:
			//prevent work from being done
			return
		default:
		}

		Monitor.WaitGroup.Add(1)
		handler.ServeHTTP(w, r)
		defer Monitor.WaitGroup.Done()
	})
}