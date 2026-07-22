package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authclient "github.com/leoldding/auth-middleware/client"
	authmiddleware "github.com/leoldding/auth-middleware/middleware"

	"github.com/leoldding/coffee/internal/awsclients"
	"github.com/leoldding/coffee/internal/config"
	"github.com/leoldding/coffee/internal/handlers"
	"github.com/leoldding/coffee/internal/middlewares"
)

func main() {
	clients, err := awsclients.New()
	if err != nil {
		log.Fatalf("failed to initialize AWS clients: %v", err)
	}

	cfg, err := config.Load(context.Background(), clients.SSM)
	if err != nil {
		log.Fatal(err)
	}

	authClient := authclient.NewAuthClient(cfg.AuthServiceURL)

	coffeeHandler := handlers.NewCoffeeHandler(
		clients,
		cfg.DynamoDBTable,
	)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/coffee/health", healthcheck)

	r.Route("/api/v1/coffee", func(r chi.Router) {
		r.Use(middlewares.CORS())

		r.Get("/status", coffeeHandler.GetStatus)

		r.With(authmiddleware.AuthCheck(authClient)).
			Post("/status", coffeeHandler.PutStatus)
	})

	log.Println("coffee service running on :8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("coffee service is healthy"))
}
