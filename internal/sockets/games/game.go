package games

import "mathly/internal/models"

type Game interface {
	StartTheGame()
	HandleMessage(msg models.Message)
}

type GameConfig struct {
	MessageListener chan models.Message
}
