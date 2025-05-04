package utils

import "github.com/google/uuid"

type Game interface {
	StartTheGame()
	HandleMessage(id, m uuid.UUID, msg string)
}
