package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// Initialize ensures the database and collections exist
func Initialize(client *mongo.Client, dbName string) error {
	ctx := context.Background()
	db := client.Database(dbName)

	// Create collections if they don't exist
	collections := []string{"users", "posts"}
	for _, collName := range collections {
		err := db.CreateCollection(ctx, collName)
		if err != nil {
			// Collection might already exist, which is fine
			log.Printf("Collection '%s' creation: %v", collName, err)
		} else {
			log.Printf("Created collection: %s", collName)
		}
	}

	log.Printf("Database '%s' initialized with collections: %v", dbName, collections)
	return nil
}

// CreateIndexes creates necessary indexes for better performance
func CreateIndexes(client *mongo.Client, dbName string) error {
	ctx := context.Background()
	db := client.Database(dbName)

	// Create indexes for users collection
	usersCollection := db.Collection("users")
	usersIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"email": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"username": 1,
			},
		},
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, usersIndexes)
	if err != nil {
		log.Printf("Error creating user indexes: %v", err)
		return err
	}

	// Create indexes for posts collection
	postsCollection := db.Collection("posts")
	postsIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"user_id": 1,
			},
		},
		{
			Keys: map[string]interface{}{
				"created_at": -1,
			},
		},
	}

	_, err = postsCollection.Indexes().CreateMany(ctx, postsIndexes)
	if err != nil {
		log.Printf("Error creating post indexes: %v", err)
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}
