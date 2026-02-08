package main

import (
	"net/http"
	"fmt"
)

// Metrics endpoint handler. It shows the number of requests made to the server.
func (c *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type","text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, c.fileserverHits.Load())))
}

// This is a metric middleware. It keeps tracks of the number of requests received
func (c *apiConfig) middlewareMetricsInc(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, req *http.Request) {
		c.fileserverHits.Add(1)
		nextHandler.ServeHTTP(w,req)
	})
}

