package main

import (
	"encoding/json"
	"net/http"
	"log"
	"strings"
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

	// Check message length
	const maxChirpLength = 140
	if len(msg.Content) > maxChirpLength {
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
	cleanMsg := replaceProfaneWords(msg.Content,censoredProfanity,profaneWords)

	// Response
	type respStruct struct {
        CleanMsg string `json:"cleaned_body"`
    }
    respondWithJSON(w, http.StatusOK, respStruct {
        CleanMsg: cleanMsg,
    })
}
