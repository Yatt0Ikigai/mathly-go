package utils

import "github.com/google/uuid"

type Game interface {
	StartTheGame()
	HandleMessage(userID uuid.UUID, msg string)
}
