package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {

	mux := chi.NewRouter()

	// Enable Cross-Origin Resource Sharing middleware to specify who is allowed to connect.
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},                                  // Allow requests from any origin with HTTP or HTTPS.
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                // Allow specified HTTP methods.
		AllowedHeaders:   []string{"Accept", "Authorizaton", "Content-Type", "X-CSRG-Token"}, // Allow specific headers.
		ExposedHeaders:   []string{"Link"},                                                   // Expose the 'Link' header in responses.
		AllowCredentials: true,                                                               // Allow sending credentials (e.g., cookies) with requests.
		MaxAge:           300,                                                                // Cache preflight (OPTIONS) request results for 300 seconds.
	}))

	// Add a middleware that responds to a '/ping' endpoint with a heartbeat message.
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/send", app.SendMail)

	// Return the configured router as an HTTP handler.
	return mux
}
