package main

import (
	"context"
	"log"

	"github.com/Nutan-Kum12/Gopherso/internal/db"
	"github.com/Nutan-Kum12/Gopherso/internal/env"
	"github.com/Nutan-Kum12/Gopherso/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			uri:         env.GetString("DB_URI", "mongodb://admin:password@localhost:27017/gopherso?authSource=admin"),
			name:        env.GetString("DB_NAME", "gopherso"),
			maxPoolSize: uint64(env.GetInt("DB_MAX_POOL_SIZE", 25)),
			minPoolSize: uint64(env.GetInt("DB_MIN_POOL_SIZE", 5)),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	client, err := db.New(
		cfg.db.uri,
		cfg.db.maxPoolSize,
		cfg.db.minPoolSize,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer client.Disconnect(context.Background())
	log.Println("MongoDB connection established.")

	// Initialize database and collections
	if err := db.Initialize(client, cfg.db.name); err != nil {
		log.Printf("Warning: Error initializing database: %v", err)
	}

	// Create database indexes for better performance
	if err := db.CreateIndexes(client, cfg.db.name); err != nil {
		log.Printf("Warning: Error creating indexes: %v", err)
	}

	store := store.NewStorage(client, cfg.db.name)
	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
