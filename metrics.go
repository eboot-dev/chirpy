package main

import (
	"net/http"
	"sync/atomic"
	"fmt"
	"log"
)

type apiConfig struct {
	fileserverHits atomic.Int32 // counts how many times a request is made
}

// Metrics endpoint handler. It shows the number of requests made to the server.
func (c *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", c.fileserverHits.Load())))
}

// Reset Metrics endpoint handler. It resets the counter of number of requests made to the server.
func (c *apiConfig) metricsResetHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Resetting fileserverHits... (last value %d)",c.fileserverHits.Load())
	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	// This might be inefficient, see: https://stackoverflow.com/questions/37863374/whats-the-difference-between-responsewriter-write-and-io-writestring
	w.Write([]byte("Hits reset to 0"))
}

// This is a metric middleware. It keeps tracks of the number of requests received
func (c *apiConfig) middlewareMetricsInc(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, req *http.Request) {
		// log.Printf("Old hits: %d",c.fileserverHits.Load())
		c.fileserverHits.Add(1)
		nextHandler.ServeHTTP(w,req)
	})
}

