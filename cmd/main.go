package main

import (
	"book-collection/api"
	"log"
	"net/http"
)

func main(){
	router := http.NewServeMux()
	router.HandleFunc("GET /users/{userID}", api.GetUsers)


	log.Printf("Server has started %s", ":8080")
	server := http.Server{
		Addr: ":8080",
		Handler: api.MiddleWareChain(router),
	}
	log.Fatal(server.ListenAndServe())
}