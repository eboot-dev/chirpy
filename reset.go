package main

import (
	"net/http"
	"log"
)


// Reset Metrics endpoint handler. It resets the counter of number of requests made to the server.
func (c *apiConfig) metricsResetHandler(w http.ResponseWriter, req *http.Request) {
	if c.platform != "dev" {
		log.Println("ERROR: Accessing reset endpoint outside `dev` environment")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in `dev` environment."))
		return
	}

	log.Printf("Reset hits [last value %d] and `users` table",c.fileserverHits.Load())
	err := c.db.DeleteUsers(req.Context())
	if err != nil {
		log.Println("ERROR: Can't reset `users` table %s",err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset `users` table: " + err.Error()))
	}
	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	// This might be inefficient, see: https://stackoverflow.com/questions/37863374/whats-the-difference-between-responsewriter-write-and-io-writestring
	w.Write([]byte("Hits reset to 0 and `users` table reset to initial state."))
}
