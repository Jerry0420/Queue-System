package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Store struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	SessionID   string    `json:"session_id"`
}

type TokenClaims struct {
	StoreID        int    `json:"store_id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	StoreCreatedAt int64  `json:"store_created_at"`
	SignKeyID      int    `json:"signkey_id"`
	jwt.StandardClaims
}
