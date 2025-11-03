package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Content   string             `json:"content" bson:"content"`
	Title     string             `json:"title" bson:"title"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Tags      []string           `json:"tags" bson:"tags"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// PostWithUser represents a post with user information
type PostWithUser struct {
	Post `bson:",inline"`
	User User `json:"user" bson:"user,omitempty"`
}

type PostStore struct {
	collection      *mongo.Collection
	usersCollection *mongo.Collection
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, post)
	if err != nil {
		return err
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a post by its ID
func (s *PostStore) GetByID(ctx context.Context, postID primitive.ObjectID) (*Post, error) {
	var post Post
	err := s.collection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByUserID retrieves all posts by a specific user
func (s *PostStore) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]Post, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID},
		options.Find().SetSort(bson.D{{"created_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetWithUser retrieves a post with user information
func (s *PostStore) GetWithUser(ctx context.Context, postID primitive.ObjectID) (*PostWithUser, error) {
	// First get the post
	post, err := s.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Then get the user information
	var user User
	err = s.usersCollection.FindOne(ctx, bson.M{"_id": post.UserID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &PostWithUser{
		Post: *post,
		User: user,
	}, nil
}

// GetAllWithUsers retrieves all posts with user information
func (s *PostStore) GetAllWithUsers(ctx context.Context, limit int64) ([]PostWithUser, error) {
	// MongoDB aggregation pipeline to join posts with users
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$sort": bson.D{{"created_at", -1}},
		},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []PostWithUser
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// Update updates a post (only by the owner)
func (s *PostStore) Update(ctx context.Context, postID, userID primitive.ObjectID, updateData bson.M) error {
	updateData["updated_at"] = time.Now()

	filter := bson.M{
		"_id":     postID,
		"user_id": userID, // Ensure only the owner can update
	}

	update := bson.M{"$set": updateData}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete deletes a post (only by the owner)
func (s *PostStore) Delete(ctx context.Context, postID, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":     postID,
		"user_id": userID, // Ensure only the owner can delete
	}

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
