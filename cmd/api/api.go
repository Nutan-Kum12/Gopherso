package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}
type config struct {
	addr string
}

// chi.Mux implements http.Handler
// ⚙️ Returning http.Handler keeps your code generic (loose coupling)
// ⚠️ Returning chi.Mux ties you to that specific library (tight coupling) in the mount function(below)
// http.Handler   Interface (generic)   You can swap routers (Chi, Gin, Echo, etc.) easily.
func (app application) mount() http.Handler {
	// mux := http.NewServeMux()                                //create a new router
	// mux.HandleFunc("GET /v1/health", app.healthCheckHandler) // Register routes (URL → handler)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //this will log the start and end of each request with the elapsed processing time
	// r.Get("/v1/health", app.healthCheckHandler)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		//posts
		//auth
		//users
	})
	return r
}

func (app application) run(mux http.Handler) error {
	// mux := http.NewServeMux()
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux, // Use the custom mux as the handler
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Starting server on %s", app.config.addr)
	return srv.ListenAndServe()
}
