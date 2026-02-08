package main

import (
	"encoding/json"
	"net/http"
	"log"
	"strings"
	"github.com/google/uuid"
	"time"
	"github.com/eboot-dev/chirpy/internal/database"
)

func replaceProfaneWords(msg, replacement string,profaneWords map[string]struct{}) string {
	words := strings.Split(msg," ")
	for i,word := range words {
		_,ok := profaneWords[strings.ToLower(word)]
		if ok {
			words[i] = replacement
		}
	}
	return strings.Join(words," ")
}

func chirpHandler(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
        Body string `json:"body"`
        UserID	uuid.UUID `json:"user_id"`
    }

    decoder := json.NewDecoder(req.Body)
    chirpMsg := chirp{}
    err := decoder.Decode(&chirpMsg)
    if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w,http.StatusBadRequest,"Error decoding chirp")
		return
    }

	// Check message length
	const maxChirpLength = 140
	if len(chirpMsg.Body) > maxChirpLength {
		respondWithError(w,http.StatusBadRequest,"Chirp is too long")
		return
	}

	// Check profanity
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	const censoredProfanity = "****"
	cleanMsg := replaceProfaneWords(chirpMsg.Body,censoredProfanity,profaneWords)
	
	params := database.CreateChirpParams{
						Body: chirpMsg.Body,
						UserID: chirpMsg.UserID,
					}
	newchirp, err := apiCfg.db.CreateChirp(req.Context(), params)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w,http.StatusInternalServerError,"Error decoding chirp")
		return
	}
	// Response
	type respStruct struct {
		ID 			uuid.UUID 	`json:"id"`
		CreatedAt 	time.Time 	`json:"created_at"`
		UpdatedAt 	time.Time 	`json:"updated_at"`
		Body 		string 		`json:"body"`
        CleanedBody string 		`json:"cleaned_body"`
        UserID		uuid.UUID 	`json:"user_id"`
    }
    respondWithJSON(w, http.StatusCreated, respStruct {
        ID:				newchirp.ID,
		CreatedAt:		newchirp.CreatedAt,	
		UpdatedAt:		newchirp.UpdatedAt,
		Body:			chirpMsg.Body,
		CleanedBody: 	cleanMsg,
		UserID:			chirpMsg.UserID,
    })
}
