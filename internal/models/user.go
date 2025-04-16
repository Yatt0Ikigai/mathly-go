package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Nickname  string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
