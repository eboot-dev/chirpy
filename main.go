package main

import (
	// "fmt"
	"net/http"
	"log"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	/*
	The endpoint should simply return a '200 OK' status code 
	The endpoint should return a 'Content-Type: text/plain; charset=utf-8' header, 
	The body will contain the message 'OK' (the text associated with the 200 status code).
	*/
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	// This might be inefficient, see: https://stackoverflow.com/questions/37863374/whats-the-difference-between-responsewriter-write-and-io-writestring
	w.Write([]byte("OK\n"))
}

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

	log.Println("Starting up...\n")
	mux := http.NewServeMux()

	// FileServer: root (index.html)
	rootHandler := http.FileServer(http.Dir(rootFilePath))
	mux.Handle(rootRoutePath,middlewareLogger(http.StripPrefix(rootRoutePrefix,rootHandler)))
	
	//FileServer: assets
	assetsHandler := http.FileServer(http.Dir(rootFilePath + assetsFilePath))
	mux.Handle(assetsRoutePath,assetsHandler)

	// Readiness Endpoint
	mux.HandleFunc(readinessRoutePath,readinessHandler)

	server := &http.Server{
		Handler: mux,
		Addr: ":" +port,
	}
	log.Printf("Serving files from %s on port: %s\n", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}
