package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Nutan-Kum12/Gopherso/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}
type config struct {
	addr string
	db   dbConfig
}
type dbConfig struct {
	addr         string // Address of your DB (host + port + dbname + user + password)
	maxOpenConns int    // Maximum number of open (in-use + idle) connections
	maxIdleConns int    // Maximum number of idle (ready-to-use) connections
	maxIdleTime  string // Duration string, e.g. "15m" meaning 15 minutes
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
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
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
