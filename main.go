package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	fmt.Println("Starting up...")
	mux := http.NewServeMux()
	mux.Handle("/",http.FileServer(http.Dir(".")))
	mux.Handle("/assets",http.FileServer(http.Dir("./assets")))
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
