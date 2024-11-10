package main

import (
	"book-collection/api"
	"log"
	"net/http"
)

func main(){
	router := http.NewServeMux()
	router.HandleFunc("GET users/{userID}", api.GetUsers)


	log.Printf("Server has started %s", ":8080")
	log.Fatal(http.ListenAndServe(":8080", api.MiddleWareChain(router)))
}