package server

import (
	"context"
	"fmt"
	"github.com/jdubbwya/go-experiment1/handlers"
	"github.com/jdubbwya/go-experiment1/handlers/middleware"
	"github.com/jdubbwya/go-experiment1/stats"
	"io"
	"log"
	"net/http"
)

func Start(addr *string) Instance {

	if addr == nil {
		var defaultAddr = "localhost:8080"
		addr = &defaultAddr
	}

	var handler = http.NewServeMux()

	instance := Instance{
		baseUrl: fmt.Sprintf("http://%s", *addr),
		statsAggregator: stats.NewStatsAggregator(),
		monitor: newMonitor(),
		server: http.Server{
			Addr : *addr,
			Handler: handler,
		},
	}

	handler.Handle("/hash/",
		monitorHandler(
			middleware.SecurityOnlyAllowMethods(
				http.HandlerFunc(handlers.HashTransactionDetailHandler),
				[]string{http.MethodGet}),
			instance.monitor))

	handler.Handle(
		"/hash",
		monitorHandler(
			middleware.SecurityOnlyAllowMethods(
				middleware.StatsMetric(
					http.HandlerFunc(handlers.HashAddTransactionHandler),
					instance.statsAggregator),
				[]string{http.MethodPost}),
			instance.monitor))

	handler.Handle("/stats",
		monitorHandler(
			middleware.SecurityOnlyAllowMethods(
				handlers.StatsHandler(instance.statsAggregator),
				[]string{http.MethodGet}),
			instance.monitor))

	handler.Handle(
		"/shutdown",
		middleware.SecurityOnlyAllowMethods(
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				close(instance.monitor.channel)
				defer func() {
					go func() {
						instance.monitor.waitGroup.Wait()
						if err := instance.server.Shutdown(context.Background()); err != nil {
							// Error from closing listeners, or context timeout:
							log.Printf("HTTP server Shutdown: %v", err)
						}
					}()
				}()
				log.Println("Shutdown requested")
				writer.Header().Add("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusOK)
				io.WriteString(writer, "Shutting down")
			}),
			[]string{http.MethodGet}))


	defer instance.server.ListenAndServe()
	log.Println("Server listening at http://localhost:8080")

	return instance
}
