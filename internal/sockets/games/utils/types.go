package utils

import (
	"encoding/json"
	"mathly/internal/models"
	"mathly/internal/utils"
)

type LobbyEvent string

var (
	LobbyEventStartOfGame LobbyEvent = "StartOfGame"
	LobbyEventEndOfGame   LobbyEvent = "EndOfGame"
	LobbyEventScoreboard  LobbyEvent = "Scoreboard"
)

func (l LobbyEvent) String() string {
	switch l {
	case LobbyEventStartOfGame:
		return "StartOfGame"
	case LobbyEventEndOfGame:
		return "EndOfGame"
	case LobbyEventScoreboard:
		return "Scoreboard"
	default:
		return "Unknown"
	}
}

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
