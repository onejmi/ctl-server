package util

import (
	"encoding/json"
	"net/http"

	"../../data"
)

//JSONResponse - Sets the response type of the endpoint to be JSON
func JSONResponse(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

//OnlyMethod - Restricts access to endpoint to one HTTP method
func OnlyMethod(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Avoid CORS issue when working locally:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method != method {
			jsonError, _ := json.Marshal(data.Error{Message: "Method not supported."})
			w.Write(jsonError)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
