package main

import (
	// "fmt"
	"net/http"
	"log"
	"sync/atomic"
	"encoding/json"
)


/* Utils */

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorStruct struct {
		Error string `json:"error"`
	}
	err := errorStruct { Error: msg }
	respondWithJSON(w,code,err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/* Handlers */
func validationHandler(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
        Content string `json:"body"`
    }

    decoder := json.NewDecoder(req.Body)
    msg := chirp{}
    err := decoder.Decode(&msg)
    if err != nil {
		log.Printf("Error decoding chirp body: %s", err)
		respondWithError(w,http.StatusBadRequest,"Error decoding chirp body")
		return
    }

	if len(msg.Content) > 140 {
		respondWithError(w,http.StatusBadRequest,"Chirp is too long")
		return
	}

	// Response
	type respStruct struct {
        Valid bool `json:"valid"`
    }
    response := respStruct{
        Valid: true,
    }
	respondWithJSON(w, http.StatusOK, response)
}
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

	const readinessRoutePath = "GET /api/healthz"
	const metricsRoutePath = "GET /admin/metrics"
	const metricsResetRoutePath = "POST /admin/reset"

	const validationRoutePath = "POST /api/validate_chirp"
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
	
	// Validation Chirp Endpoint
	mux.HandleFunc(validationRoutePath,validationHandler)
	server := &http.Server{
		Handler: mux,
		Addr: ":" +port,
	}
	log.Printf("Serving files from %s on port: %s", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}
