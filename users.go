package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c * apiConfig) usersCreateHandler(w http.ResponseWriter, req *http.Request){
	type userInput struct {
        Email string `json:"email"`
    }

    decoder := json.NewDecoder(req.Body)
    input := userInput{}
    err := decoder.Decode(&input)
    if err != nil {
		log.Printf("ERROR: usersCreateHandler couldn't decode user input [%s]", err)
		respondWithError(w,http.StatusInternalServerError,"Couldn't decode user input")
		return
    }

	user, err := c.db.CreateUser(req.Context(), input.Email)
	if err != nil {
		log.Printf("ERROR: DB couldn't create user [%s]", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
    	
	respondWithJSON(w, http.StatusCreated, User{
		ID: 		user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email: 		user.Email,
    })
}
