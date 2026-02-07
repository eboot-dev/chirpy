package main

import (
	"encoding/json"
	"net/http"
	"log"
	"strings"
)

func replaceProfaneWords(msg, replacement string,profaneWords map[string]struct{}) string {
	words := strings.Split(msg," ")
	clean := make([]string,len(words))
	copy(clean,words)
	for i,word := range words {
		_,ok := profaneWords[strings.ToLower(word)]
		if ok {
			clean[i] = replacement
		}
	}
	return strings.Join(clean," ")
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
	if len(msg.Content) > 140 {
		respondWithError(w,http.StatusBadRequest,"Chirp is too long")
		return
	}

	// Check profanity
	profaneWords := map[string]struct{}{
		"kerfuffle":struct{}{},
		"sharbert":struct{}{},
		"fornax":struct{}{},
	}
	cleanMsg := replaceProfaneWords(msg.Content,"****",profaneWords)

	// Response
	type respStruct struct {
        CleanMsg string `json:"cleaned_body"`
    }
    response := respStruct{
        CleanMsg: cleanMsg,
    }
	respondWithJSON(w, http.StatusOK, response)
}
