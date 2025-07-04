package utils

import (
	"encoding/json"
	"mathly/internal/models"
	"mathly/internal/utils"
)



type GameMessage struct {
	Type    string
	Message string
}

func (g GameMessage) ToByteArray() ([]byte, error) {
	return json.Marshal(g)
}

type GameConfig struct {
	Random          utils.Random
	MessageListener chan models.Message
}
