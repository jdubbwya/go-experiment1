package main

import (
	"github.com/jdubbwya/go-experiment1/hasher"
	"github.com/jdubbwya/go-experiment1/middleware"
	"github.com/jdubbwya/go-experiment1/stats"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	http.Handle(
		"/hash",
		middleware.MonitorRequest(
			middleware.AllowOnly(
				http.HandlerFunc(hasher.HandleQueue),
				[]string{http.MethodPost, http.MethodGet})))

	http.Handle("/hash/\\d+",
		middleware.MonitorRequest(
			middleware.AllowOnly(
				http.HandlerFunc(hasher.HandleEntry),
				[]string{http.MethodGet})))

	http.Handle("/stats",
		middleware.MonitorRequest(
			middleware.AllowOnly(
				http.HandlerFunc(stats.HandleStats),
				[]string{http.MethodGet})))

	http.Handle(
		"/shutdown",
		middleware.AllowOnly(
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				close(middleware.Monitor.Channel)
				writer.Header().Add("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusOK)
				io.WriteString(writer, "Shutting down")
				go func() {
					middleware.Monitor.WaitGroup.Wait()
					time.Sleep(time.Second)
					os.Exit(0)
				}()
			}),
			[]string{http.MethodGet}))

	http.ListenAndServe(":8080", nil)
}
