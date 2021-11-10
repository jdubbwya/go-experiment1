package main

import (
	"github.com/jdubbwya/go-experiment1/hasher"
	"github.com/jdubbwya/go-experiment1/middleware"
	"github.com/jdubbwya/go-experiment1/stats"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	http.Handle("/hash/",
		middleware.MonitorRequest(
			middleware.AllowOnly(
				http.HandlerFunc(hasher.HandleEntry),
				[]string{http.MethodGet})))

	http.Handle(
		"/hash",
		middleware.MonitorRequest(
			middleware.AllowOnly(
				http.HandlerFunc(hasher.HandleQueue),
				[]string{http.MethodPost})))

	http.Handle("/stats",
		middleware.MonitorRequest(
			middleware.AllowOnly(
				http.HandlerFunc(stats.HandleStats),
				[]string{http.MethodGet})))

	http.Handle(
		"/shutdown",
		middleware.AllowOnly(
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				log.Println("Shutdown requested")
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

	defer http.ListenAndServe(":8080", nil)
	log.Println("Server listening at http://localhost:8080")
}
