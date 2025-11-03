package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nutan-Kum12/Gopherso/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePostRequest represents the JSON payload for creating a post
type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required,max=5000"`
	UserID  string   `json:"user_id" validate:"required"`
	Tags    []string `json:"tags,omitempty"`
}

// PostResponse represents the JSON response for post data
type PostResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// createPostHandler handles POST /v1/posts
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Basic validation
	if req.Title == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.Content == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Content is required")
		return
	}
	if req.UserID == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Convert user ID to ObjectID
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Verify user exists
	_, err = app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "User not found")
		return
	}

	// Create post
	post := &store.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
		Tags:    req.Tags,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Return post response
	response := PostResponse{
		ID:        post.ID.Hex(),
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID.Hex(),
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	app.writeJSONResponse(w, http.StatusCreated, response)
}

// getPostHandler handles GET /v1/posts/{id}
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), postID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusNotFound, "Post not found")
		return
	}

	response := PostResponse{
		ID:        post.ID.Hex(),
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID.Hex(),
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	app.writeJSONResponse(w, http.StatusOK, response)
}

// getPostWithUserHandler handles GET /v1/posts/{id}/user
func (app *application) getPostWithUserHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	if postIDStr == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid post ID format")
		return
	}

	postWithUser, err := app.store.Posts.GetWithUser(r.Context(), postID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusNotFound, "Post not found")
		return
	}

	app.writeJSONResponse(w, http.StatusOK, postWithUser)
}

// getPostsHandler handles GET /v1/posts (get all posts with user info)
func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Get limit from query parameter (default to 20)
	limit := int64(20)

	posts, err := app.store.Posts.GetAllWithUsers(r.Context(), limit)
	if err != nil {
		app.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	app.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"posts": posts,
		"count": len(posts),
	})
}

// getPostsByUserHandler handles GET /v1/users/{id}/posts
func (app *application) getPostsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	posts, err := app.store.Posts.GetByUserID(r.Context(), userID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve user posts")
		return
	}

	app.writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"posts": posts,
		"count": len(posts),
	})
}
