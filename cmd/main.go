package main

import(
	"book-collection/api"
)
func main(){
	server := api.NewAPIServer(":8080")
	server.Run()
}