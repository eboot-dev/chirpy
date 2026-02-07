package main

import (
	// "fmt"
	"net/http"
	"log"
	"sync/atomic"
	"encoding/json"
	"database/sql"
	"os"
	"github.com/joho/godotenv"
	"time"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/eboot-dev/chirpy/internal/database"
)

/* Utils */

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	err := errorResponse { Error: msg }
	respondWithJSON(w,code,err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

/* Handlers */
func usersHandler(w http.ResponseWriter, req *http.Request){
	type userInput struct {
        Email string `json:"email"`
    }

    decoder := json.NewDecoder(req.Body)
    input := userInput{}
    err := decoder.Decode(&input)
    if err != nil {
		log.Printf("Error decoding user input: %s", err)
		respondWithError(w,http.StatusBadRequest,"Error decoding user input")
		return
    }

	user, err := apiCfg.db.CreateUser(req.Context(), input.Email)
	// Response
    type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}	
	respondWithJSON(w, http.StatusCreated, User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
    })
}
/* Middlewares */

// This is a looging middleware. It logs the Method and URL.Path of a request and pass it to the next handler to be processed
func middlewareLogger(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.URL.Path)
		nextHandler.ServeHTTP(w, req)
	})
}

var apiCfg apiConfig

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
	const usersRoutePath = "POST /api/users"	
	
	// Load the `.env` file
	godotenv.Load()
	// Connect to db
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	environment := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("ERROR: Can't establish connection with DB %s",err)
	}
	dbQueries := database.New(db)
	// init apiConfig struct
	apiCfg = apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		env: environment,
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

	mux.HandleFunc(usersRoutePath,usersHandler)
	server := &http.Server{
		Handler: mux,
		Addr: ":" +port,
	}
	log.Printf("Serving files from %s on port: %s", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}
