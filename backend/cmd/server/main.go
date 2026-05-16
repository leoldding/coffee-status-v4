package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/leoldding/coffee/internal/awsclients"
	"github.com/leoldding/coffee/internal/handlers"
	"github.com/leoldding/coffee/internal/middlewares"
)

func main() {
	clients, err := awsclients.New()
	if err != nil {
		log.Fatal("failed to initialize aws clients")
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/coffee/health", healthcheck)

	r.Get("/api/v1/coffee/test", test)

	coffeeHandler := handlers.NewCoffeeHandler(clients)
	r.Route("/api/v1/coffee", func(r chi.Router) {
		r.Use(middlewares.CORS())

		r.Get("/status", coffeeHandler.GetStatus)
		r.With(middlewares.AuthCheck).Post("/status", coffeeHandler.PutStatus)
	})

	log.Println("coffee service running on port 8081")
	http.ListenAndServe(":8081", r)
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("coffee service is healthy"))
	return
}

func test(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("new bucket works"))
	return
}
