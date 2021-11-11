package server

import (
	"github.com/jdubbwya/go-experiment1/handlers"
	"github.com/jdubbwya/go-experiment1/handlers/middleware"
	"io"
	"log"
	"net/http"
)

func Start() {
	http.Handle("/hash/",
		MonitorHandler(
			middleware.SecurityOnlyAllowMethods(
				http.HandlerFunc(handlers.HashTransactionDetailHandler),
				[]string{http.MethodGet})))

	http.Handle(
		"/hash",
		MonitorHandler(
			middleware.SecurityOnlyAllowMethods(
				http.HandlerFunc(handlers.HashAddTransactionHandler),
				[]string{http.MethodPost})))

	http.Handle("/stats",
		MonitorHandler(
			middleware.SecurityOnlyAllowMethods(
				http.HandlerFunc(handlers.StatsHandler),
				[]string{http.MethodGet})))

	http.Handle(
		"/shutdown",
		middleware.SecurityOnlyAllowMethods(
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				log.Println("Shutdown requested")
				notifyShutdown()
				writer.Header().Add("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusOK)
				io.WriteString(writer, "Shutting down")
				go drainConnectionsThenShutdown()
			}),
			[]string{http.MethodGet}))

	defer http.ListenAndServe(":8080", nil)
	log.Println("Server listening at http://localhost:8080")
}