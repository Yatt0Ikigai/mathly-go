package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Nickname string    `json:"nickname"`
	jwt.StandardClaims
}
