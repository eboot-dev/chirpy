package main

import (
	// "fmt"
	"net/http"
	"log"
	"sync/atomic"
)


/* Handlers */

/* Middlewares */

// This is a looging middleware. It logs the Method and URL.Path of a request and pass it to the next handler to be processed
func middlewareLogger(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.URL.Path)
		nextHandler.ServeHTTP(w, req)
	})
}

func main() {
	const port = "8080"

	const rootRoutePath = "/app/"
	const rootRoutePrefix = "/app"
	const rootFilePath = "."
	
	const assetsRoutePath = "/assets"
	const assetsFilePath = "/assets"

	const readinessRoutePath = "/healthz"
	const metricsRoutePath = "/metrics"
	const metricsResetRoutePath = "/reset"

	// init apiConfig struct
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	log.Println("Starting up...")
	mux := http.NewServeMux()
	/* Registering Handlers */
	// FileServer: root (index.html)
	rootHandler := http.FileServer(http.Dir(rootFilePath))
	rootHandler = apiCfg.middlewareMetricsInc(rootHandler)
	rootHandler = http.StripPrefix(rootRoutePrefix,rootHandler)
	mux.Handle(rootRoutePath,middlewareLogger(rootHandler))
	
	//FileServer: assets
	assetsHandler := http.FileServer(http.Dir(rootFilePath + assetsFilePath))
	mux.Handle(assetsRoutePath,assetsHandler)

	// Readiness Endpoint
	mux.HandleFunc(readinessRoutePath,readinessHandler)

	// Metrics Endpoint
	mux.HandleFunc(metricsRoutePath,apiCfg.metricsHandler)
	// Metrics Reset Endpoint
	mux.HandleFunc(metricsResetRoutePath,apiCfg.metricsResetHandler)
	
	server := &http.Server{
		Handler: mux,
		Addr: ":" +port,
	}
	log.Printf("Serving files from %s on port: %s", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}
