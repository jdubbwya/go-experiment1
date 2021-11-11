package middleware

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

func SecurityOnlyAllowMethods(handler http.Handler, allowMethods []string) http.Handler {
	var allowList = []byte(strings.Join(allowMethods, ","))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var matched, err = regexp.Match(r.Method, allowList)
		if !matched || err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method not allowed for requested resource")
			return
		}

		handler.ServeHTTP(w, r)
	})
}