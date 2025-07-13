package common_games

import (
	"encoding/json"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/shared"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type GameMessage struct {
	Type    string
	Message string
}

func (g GameMessage) ToByteArray() ([]byte, error) {
	return json.Marshal(g)
}

type GameServices struct {
	Random service.Random
}

type GameSettings map[string]any

type GameConfig struct {
	Services  GameServices
	Settings  GameSettings
	Broadcast chan shared.SocketResponse
	Players   map[uuid.UUID]models.Player
	Scheduler gocron.Scheduler
	EndGame   func()
}
