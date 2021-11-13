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
	"sync"
)

const StateStop = 0
const StateUp = StateStop+1
const StateQuit = StateUp+1

type instanceState struct {
	value int
	mu sync.Mutex
}

func (is *instanceState) nextState( next int ){
	is.mu.Lock()
	is.value = next
	is.mu.Unlock()
}

type Instance struct {
	baseUrl string
	state *instanceState
	waitGroup sync.WaitGroup
	server *http.Server
	statsAggregator *stats.Aggregator
}

func (i *Instance) Aggregator() *stats.Aggregator {
	return i.statsAggregator
}

func (i *Instance ) IsAlive() bool  {
	return i.state.value == StateUp
}

func (i *Instance) BaseUrl() string {
	return i.baseUrl
}

func (i *Instance) Kill() error {
	err := i.server.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
		return err
	}

	i.state.nextState(StateStop)
	i.statsAggregator.Stop()
	return nil
}
func (i *Instance) Start(){
	i.state.nextState(StateUp)
	if err := i.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (i *Instance) Quit(){
	i.state.nextState(StateQuit)
	i.waitGroup.Add(1)
	defer i.waitGroup.Done()
	go func() {
		i.waitGroup.Wait()
		i.Kill()
	}()
}

func onlyAcceptRequestsWhenAlive(handler http.Handler, i Instance) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !i.IsAlive() {
			log.Println("Request not accepted")
			return
		}

		i.waitGroup.Add(1)
		defer i.waitGroup.Done()
		handler.ServeHTTP(w, r)
	})
}

func NewInstance(addr *string) Instance{
	if addr == nil {
		var defaultAddr = "localhost:8080"
		addr = &defaultAddr
	}

	var handler = http.NewServeMux()
	var server = http.Server{
		Addr : *addr,
		Handler: handler,
	}
	var state = instanceState{
		value: StateStop,
	}

	instance := Instance{
		baseUrl: fmt.Sprintf("http://%s", *addr),
		statsAggregator: stats.NewStatsAggregator(),
		state: &state,
		server: &server,
	}

	handler.Handle("/hash/",
		onlyAcceptRequestsWhenAlive(
			middleware.SecurityOnlyAllowMethods(
				http.HandlerFunc(handlers.HashTransactionDetailHandler),
				[]string{http.MethodGet}),
			instance))

	handler.Handle(
		"/hash",
		onlyAcceptRequestsWhenAlive(
			middleware.SecurityOnlyAllowMethods(
				middleware.StatsMetric(
					http.HandlerFunc(handlers.HashAddTransactionHandler),
					instance.statsAggregator),
				[]string{http.MethodPost}),
			instance))

	handler.Handle("/stats",
		onlyAcceptRequestsWhenAlive(
			middleware.SecurityOnlyAllowMethods(
				handlers.StatsHandler(instance.statsAggregator),
				[]string{http.MethodGet}),
			instance))

	handler.Handle(
		"/shutdown",
		middleware.SecurityOnlyAllowMethods(
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Add("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusOK)
				io.WriteString(writer, "Shutting down")
				instance.Quit()
			}),
			[]string{http.MethodGet}))

	return instance
}