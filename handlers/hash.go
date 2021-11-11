package handlers

import (
	"fmt"
	"github.com/jdubbwya/go-experiment1/hasher"
	"github.com/jdubbwya/go-experiment1/server/responses"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HashAddTransactionHandler(w http.ResponseWriter, r *http.Request){

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		responses.BadRequest(w)
		return
	}

	var rawPassword = r.Form.Get("password")

	if len(rawPassword) == 0 {
		responses.BadRequest(w)
		return
	}

	var id, err = hasher.Enqueue(rawPassword)

	if err != nil {
		log.Fatal(err)
		responses.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("%d", id))

}

func HashTransactionDetailHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain")

	var pathParts = strings.Split(r.URL.Path, "/")

	if len(pathParts) > 3 {
		responses.NotFound(w)
		return
	}

	var entryId = pathParts[2]
	var id, err = strconv.ParseInt(entryId, 0, 64)
	if err != nil {
		responses.BadRequest(w)
		return
	}

	var hashedPassword, transactionState = hasher.TransactionResult(id)

	switch transactionState {
		case hasher.TransactionUnknown:
			responses.BadRequest(w)
			return
		case hasher.TransactionInProgress:
			w.WriteHeader(http.StatusTooEarly)
			io.WriteString(w, "Request is still processing")
			return
		case hasher.TransactionComplete:
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, *hashedPassword)
			return
	}

	// respond with any unknown states
	responses.InternalServerError(w)
}