package main

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func startServer() {
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

	router.Route("/subscription/{topic}", func(router chi.Router) {
		router.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			handleSubscription(writer, request, "POST", messagingClient.SubscribeToTopic)
		})

		router.Delete("/", func(writer http.ResponseWriter, request *http.Request) {
			handleSubscription(writer, request, "DELETE", messagingClient.UnsubscribeFromTopic)
		})
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func handleSubscription(writer http.ResponseWriter, request *http.Request, method string, firebaseFunction func(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)) {
	topic := chi.URLParam(request, "topic")

	token := request.Header.Get("Authorization")
	if token == "" {
		log.Println(method, "/subscription/"+topic+" - no token")
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte("no token"))
		return
	}

	toTopic, err := firebaseFunction(context.Background(), []string{token}, topic)
	if err != nil {
		log.Println(method, "/subscription/"+topic+" - error: "+err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	if toTopic.FailureCount > 0 {
		log.Println(method, "/subscription/"+topic+" - error: "+toTopic.Errors[0].Reason)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(toTopic.Errors[0].Reason))
		return
	}

	log.Println(method, "/subscription/"+topic+" - success")
}
