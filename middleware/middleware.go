package middleware

import (
	"io"
	"net/http"
	"regexp"
	"strings"
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

func MonitorRequest(handler http.Handler) http.Handler {
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

func AllowOnly(handler http.Handler, allowMethods []string) http.Handler {
	var allowList = []byte(strings.Join(allowMethods, ","))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var matched, err = regexp.Match(r.Method, allowList)
		if !matched || err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method not allowed for requested resource")
			return
		}

		handler.ServeHTTP(w, r)
	})
}
