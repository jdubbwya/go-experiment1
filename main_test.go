package main

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/server"
	"log"
	"net/http"
	"sync"
	"testing"
)

var testUrl string = "localhost:8080"
var runningInstance *server.Instance

func makeUrl( path string ) string {
	return fmt.Sprintf("http://%s%s", testUrl, path)
}

func beforeEach(){
	if runningInstance == nil || !runningInstance.IsAlive() {
		instance := server.NewInstance(&testUrl)
		runningInstance = &instance
		go func() {
			runningInstance.Start()
		}()
	}

}

type testCase struct {
	name string
	test func(t *testing.T)
}

func TestServerSuite(t *testing.T) {

	cases := []testCase{
		testCase{
			"Server: Verify Graceful shutdown",
			serverTestCaseGracefulShutdown,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			beforeEach()
			c.test(t)
		})
	}

}

func serverTestCaseGracefulShutdown(t *testing.T) {

	var parallelGets = sync.WaitGroup{}

	type httpResponse struct {
		err *error
		statusCode int
	}

	statsResponse := httpResponse{
		err : nil,
		statusCode: http.StatusInternalServerError,
	}
	shutdownResponse := httpResponse{
		err : nil,
		statusCode: http.StatusInternalServerError,
	}

	parallelGets.Add(2)

	go func(wg *sync.WaitGroup, response *httpResponse) {
		res, err := http.Get(makeUrl("/stats"))
		(*response).statusCode = res.StatusCode
		if err != nil {
			(*response).err = &err
		}
		wg.Done()
	}( &parallelGets, &statsResponse )

	go func(wg *sync.WaitGroup, response *httpResponse) {
		res, err := http.Get(makeUrl("/shutdown"))
		(*response).statusCode = res.StatusCode
		if err != nil {
			(*response).err = &err
		}
		wg.Done()
	}( &parallelGets, &shutdownResponse )

	parallelGets.Wait()

	if statsResponse.err != nil {
		t.Fatal(statsResponse.err)
	}

	if statsResponse.statusCode > 299 {
		log.Fatalf("Status response failed with status code: %d\n", statsResponse.statusCode)
	}

	if shutdownResponse.err != nil {
		t.Fatal(shutdownResponse.err)
	}

	if shutdownResponse.statusCode > 299 {
		log.Fatalf("Shutdown response failed with status code: %d\n", shutdownResponse.statusCode)
	}
}