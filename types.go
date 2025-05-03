package main

import (
	"sync/atomic"

	"github.com/Pranay-Tej/go-chirpy/internal/database"
)

type ApiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}
