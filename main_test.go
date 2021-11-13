package main

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/server"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

var testUrl string = "localhost:8080"
var runningInstance *server.Instance

func makeUrl( path string ) string {
	return fmt.Sprintf("http://%s%s", testUrl, path)
}

func beforeEach(){
	if runningInstance != nil && runningInstance.IsAlive() {
		runningInstance.Kill()
	}
	instance := server.NewInstance(&testUrl)
	runningInstance = &instance
	go func() {
		runningInstance.Start()
	}()

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
		{
			name: "Server: /hash/{id} responds after 5 seconds",
			test: serverTestCaseHashTransactionAfter5seconds,
		},
		{
			name: "Server: can handle parallel requests",
			test: serverTestCanHandleParallelRequests,
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


	type endpointTestCase struct {
		name string
		path string
		err *error
		statusCode int
	}

	testCases := [] endpointTestCase{
		endpointTestCase{
			name: "Stats returns data",
			path: "/stats",
			err: nil,
			statusCode: http.StatusInternalServerError,
		},
		endpointTestCase{
			name: "Shutdown returns data",
			path: "/shutdown",
			err: nil,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, c := range testCases {
		t.Run( c.name, func(t *testing.T) {
			t.Parallel()
			res, err := http.Get(makeUrl(c.path))
			c.statusCode = res.StatusCode
			if err != nil {
				c.err = &err
			}

			if c.err != nil {
				t.Fatal(c.err)
			}

			if c.statusCode > 299 {
				t.Fatalf("Response failed with status code: %d\n", c.statusCode)
			}
		})
	}
}

func serverTestCaseHashTransactionAfter5seconds(t *testing.T) {
	formData := url.Values{
		"password": { "angryMonkey231569745269" },
	}

	http.Post(makeUrl("/hash"), "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	time.Sleep(5 * time.Second)
	res, err := http.Get(makeUrl("/hash/1"))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Response failed with status code: %d\n", res.StatusCode)
	}

	if res.ContentLength == 0 {
		log.Fatalf("Response contained no data\n")
	}
}

func serverTestCanHandleParallelRequests(t *testing.T) {
	var i int
	for i = 0; i < 50; i ++ {
		t.Run(fmt.Sprintf("Parallel Request to /stats #%d", i), func(t *testing.T) {
			t.Parallel()
			res, err := http.Get(makeUrl("/stats"))
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != http.StatusOK {
				t.Fatalf("Response failed with status code: %d\n", res.StatusCode)
			}

			if res.ContentLength <= 0 {
				t.Fatalf("Response content length was 0\n")
			}
		})
	}
}