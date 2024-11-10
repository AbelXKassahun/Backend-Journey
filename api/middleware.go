package api

import(
	"net/http"
	"log"
)

type middleWare func(http.Handler) http.HandlerFunc

var MiddleWareChain = middleWareChain(
	requireAuthentication,
	requestLoggerMiddleWare,
)

func middleWareChain(middlwares ...middleWare) middleWare {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlwares) - 1; i >= 0; i-- {
			next = middlwares[i](next)
		}

		return next.ServeHTTP
	}
}


func requestLoggerMiddleWare(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func requireAuthentication(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}