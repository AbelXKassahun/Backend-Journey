package api

import (
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	router *http.ServeMux
	subrouter *http.ServeMux
	// db connection
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
		router: http.NewServeMux(),
		subrouter: http.NewServeMux(),
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request){
		userID := r.PathValue("userID")
		w.Write([]byte("User ID: "+ userID))
	})
	
	// subrouting
	v1 := http.NewServeMux()
	v1.Handle("api/v1/", http.StripPrefix("api/v1/", router))

	middleWareChain := aMiddleWareChain(
		RequestLoggerMiddleWare,
		RequireAuthentication,
	)

	server := http.Server{
		Addr: s.addr,
		Handler: middleWareChain(router),
	}

	log.Printf("Server has started %s", s.addr)
	return server.ListenAndServe()
}

func RequestLoggerMiddleWare(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func RequireAuthentication(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

type MiddleWare func(http.Handler) http.HandlerFunc

func aMiddleWareChain(middlwares ...MiddleWare) MiddleWare {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlwares) - 1; i >= 0; i-- {
			next = middlwares[i](next)
		}

		return next.ServeHTTP
	}
}