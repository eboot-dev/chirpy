package main

import (
	"fmt"
	"net/http"
	// "errors"
)

func main() {
	fmt.Println("Starting up...")
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("ERROR: can't start server, %w",err)
	}
}
