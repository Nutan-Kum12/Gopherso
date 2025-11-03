package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, primitive.ObjectID) (*Post, error)
		GetByUserID(context.Context, primitive.ObjectID) ([]Post, error)
		GetWithUser(context.Context, primitive.ObjectID) (*PostWithUser, error)
		GetAllWithUsers(context.Context, int64) ([]PostWithUser, error)
		Update(context.Context, primitive.ObjectID, primitive.ObjectID, bson.M) error
		Delete(context.Context, primitive.ObjectID, primitive.ObjectID) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, primitive.ObjectID) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		GetWithPosts(context.Context, primitive.ObjectID) (*UserWithPosts, error)
		GetPostsCount(context.Context, primitive.ObjectID) (int64, error)
	}
}

func NewStorage(client *mongo.Client, dbName string) Storage {
	db := client.Database(dbName)
	usersCollection := db.Collection("users")
	postsCollection := db.Collection("posts")

	return Storage{
		Users: &UserStore{
			collection:      usersCollection,
			postsCollection: postsCollection,
		},
		Posts: &PostStore{
			collection:      postsCollection,
			usersCollection: usersCollection,
		},
	}
}
