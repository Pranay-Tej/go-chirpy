package main

import (
	"sync/atomic"
	"time"

	"github.com/Pranay-Tej/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
type SignupInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
