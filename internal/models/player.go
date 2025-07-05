package models

import (
	"mathly/internal/shared"

	"github.com/google/uuid"
)

type Player struct {
	AccountID    *int64
	ConnectionID uuid.UUID
	Nickname     string
	SendMessage  func(shared.SocketReponse)
}
