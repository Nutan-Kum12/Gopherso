package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Nutan-Kum12/Gopherso/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserRequest represents the JSON payload for creating a user
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse represents the JSON response for user data
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// createUserHandler handles POST /v1/users
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Basic validation
	if req.Username == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Username is required")
		return
	}
	if req.Email == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "Password is required")
		return
	}

	// Check if user already exists
	existingUser, err := app.store.Users.GetByEmail(r.Context(), req.Email)
	if err == nil && existingUser != nil {
		app.writeErrorResponse(w, http.StatusConflict, "User with this email already exists")
		return
	}

	// Create user
	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // In production, hash this password!
	}

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		app.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Return user response
	response := UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	app.writeJSONResponse(w, http.StatusCreated, response)
}

// getUserHandler handles GET /v1/users/{id}
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameter (we'll implement this with chi URL params)
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	response := UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	app.writeJSONResponse(w, http.StatusOK, response)
}

// getUserWithPostsHandler handles GET /v1/users/{id}/posts
func (app *application) getUserWithPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		app.writeErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		app.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	userWithPosts, err := app.store.Users.GetWithPosts(r.Context(), userID)
	if err != nil {
		app.writeErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	app.writeJSONResponse(w, http.StatusOK, userWithPosts)
}
