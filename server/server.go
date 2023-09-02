package server

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func StartServer() {
	router := chi.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			if request.Method == http.MethodOptions {
				writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Referrer, User-Agent")
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(writer, request)
		})
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
