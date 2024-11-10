package api

import "net/http"

func GetUsers(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	w.Write([]byte("User ID: "+ userID))
}