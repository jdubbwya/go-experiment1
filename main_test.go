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

func makeUrl( path string ) string {
	return fmt.Sprintf("http://%s%s", testUrl, path)
}

func startServer(){
	go func() {
		server.Start(&testUrl)
	}()
}

func TestGracefulShutdown(t *testing.T) {

	startServer()
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