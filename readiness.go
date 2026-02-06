package main

import (
	"net/http"
)

// Readines endpoint handler. It indicate that the server is up and running
func readinessHandler(w http.ResponseWriter, req *http.Request) {
	/*
	The endpoint should simply return a '200 OK' status code 
	The endpoint should return a 'Content-Type: text/plain; charset=utf-8' header, 
	The body will contain the message 'OK' (the text associated with the 200 status code).
	*/
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

