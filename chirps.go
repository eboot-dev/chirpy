package main

import (
	"encoding/json"
	"net/http"
	"log"
	"errors"
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

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	// Check profanity
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	const censoredProfanity = "****"
	
	cleaned := replaceProfaneWords(body, censoredProfanity, profaneWords)
	return cleaned, nil
}


type Chirp struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	UserID		uuid.UUID 	`json:"user_id"`
	Body 		string 		`json:"body"`
	// CleanedBody string 		`json:"cleaned_body"`
}
 
func (c *apiConfig) getChirpHandler(w http.ResponseWriter, req *http.Request){
	log.Printf("PATHVALUE: %v\n",req.PathValue("chirpID"))
	id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("ERROR: Invalid UUID string: %v\n", err)
		respondWithError(w,http.StatusNotFound,"Coudn't get chirp")
		return
	}

	chirp, err := c.db.Chirp(req.Context(), id)
	if err != nil {
		log.Printf("ERROR: coudn't get chirp [%s]", err)
		respondWithError(w,http.StatusNotFound,"Coudn't get chirp")
		return
	}
	
    respondWithJSON(w, http.StatusOK, Chirp{
        ID:				chirp.ID,
		CreatedAt:		chirp.CreatedAt,	
		UpdatedAt:		chirp.UpdatedAt,
		Body:      		chirp.Body,
		UserID:			chirp.UserID,
    })
    return
}

func (c *apiConfig) chirpsHandler(w http.ResponseWriter, req *http.Request){

	chirps, err := c.db.Chirps(req.Context())
	if err != nil {
		log.Printf("ERROR: coudn't get chirps [%s]", err)
		respondWithError(w,http.StatusInternalServerError,"Coudn't get chirps")
		return
	}
	nChirps := len(chirps)
	log.Printf("#chirps %v", nChirps)
	log.Printf("chirps: %v", chirps)
	var res []Chirp
	for _,c := range chirps {
		res = append(res,Chirp{
        ID:				c.ID,
		CreatedAt:		c.CreatedAt,	
		UpdatedAt:		c.UpdatedAt,
		Body:      		c.Body,
		UserID:			c.UserID,
    })
	}
    respondWithJSON(w, http.StatusOK, res)
    return
}



func (c *apiConfig) chirpsCreateHandler(w http.ResponseWriter, req *http.Request) {
	type userInput struct {
        Body 	string 		`json:"body"`
        UserID	uuid.UUID 	`json:"user_id"`
    }

    decoder := json.NewDecoder(req.Body)
    input := userInput{}
    err := decoder.Decode(&input)
    if err != nil {
    	log.Printf("ERROR: chirpsCreateHandler couldn't decode user input [%s]", err)
		respondWithError(w,http.StatusInternalServerError,"Couldn't decode user input")
		return
    }

	cleanBody, err := validateChirp(input.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	
	chirp, err := c.db.CreateChirp(req.Context(), database.CreateChirpParams{
						Body: 	cleanBody,
						UserID: input.UserID,
					})
	if err != nil {
		log.Printf("ERROR: coudn't create chirp [%s]", err)
		respondWithError(w,http.StatusInternalServerError,"Coudn't create chirp")
		return
	}

    respondWithJSON(w, http.StatusCreated, Chirp{
        ID:				chirp.ID,
		CreatedAt:		chirp.CreatedAt,	
		UpdatedAt:		chirp.UpdatedAt,
		Body:      		chirp.Body,
		UserID:			chirp.UserID,
    })
}
