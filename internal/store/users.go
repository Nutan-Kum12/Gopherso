package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserWithPosts represents a user with their posts
type UserWithPosts struct {
	User  `bson:",inline"`
	Posts []Post `json:"posts" bson:"posts,omitempty"`
}

type UserStore struct {
	collection      *mongo.Collection
	postsCollection *mongo.Collection
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a user by their ID
func (s *UserStore) GetByID(ctx context.Context, userID primitive.ObjectID) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetWithPosts retrieves a user with all their posts
func (s *UserStore) GetWithPosts(ctx context.Context, userID primitive.ObjectID) (*UserWithPosts, error) {
	// First get the user
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Then get all posts by this user
	cursor, err := s.postsCollection.Find(ctx, bson.M{"user_id": userID},
		options.Find().SetSort(bson.D{{"created_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return &UserWithPosts{
		User:  *user,
		Posts: posts,
	}, nil
}

// GetPostsCount returns the number of posts for a user
func (s *UserStore) GetPostsCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := s.postsCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	return count, err
}
