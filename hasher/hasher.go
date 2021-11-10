package hasher

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/jdubbwya/go-experiment1/benchmark"
	"github.com/jdubbwya/go-experiment1/stats"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const initialWaitDuration = time.Duration( 5 * (10^9) ) // 5 nanoseconds

// Tracks the previous allocated transaction ids
var maxId uint64 = 0

// stores the values of all hashed password with a concurrency safe datastore
var hashedPasswords sync.Map


func nextId(  ) uint64 {
	atomic.AddUint64(&maxId, 1)

	return maxId
}

func HandleQueue(w http.ResponseWriter, r *http.Request){

	var id = nextId()
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing and/or empty password.")
		return
	}

	var mark *benchmark.Benchmark = benchmark.StartCapture()

	var rawPassword = r.Form.Get("password")

	if len(rawPassword) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing and/or empty password.")
		return
	}


	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("%d", id))

	go func() {
		time.Sleep(initialWaitDuration)
		h512 := sha512.New()
		io.WriteString(h512, rawPassword)
		hashedPasswords.Store( id, base64.StdEncoding.EncodeToString(h512.Sum(nil)) )
		benchmark.StopCapture(mark)
		stats.Capture(mark)
	}()
}

func HandleEntry(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain")

	var pathParts = strings.Split(r.URL.Path, "/")

	if len(pathParts) > 3 {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Resource not found.")
		return
	}

	var entryId = pathParts[2]
	var id, err = strconv.ParseUint(entryId, 0, 64)
	if err != nil || id > maxId {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Request malformed and/or invalid")
		return
	}

	var mapEntry, ok = hashedPasswords.Load(id)
	if ! ok {
		w.WriteHeader(http.StatusTooEarly)
		io.WriteString(w, "Request is still processing")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, mapEntry.(string))
}