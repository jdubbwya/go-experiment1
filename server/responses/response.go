package responses

import (
	"io"
	"net/http"
)


func BadRequest(w http.ResponseWriter){
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "Request malformed and/or invalid")
}

func NotFound(w http.ResponseWriter){
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Resource not found.")
}

func InternalServerError(w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "An unexpected error occurred.")
}