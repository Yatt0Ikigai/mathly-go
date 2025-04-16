package service

import (
	"mathly/internal/models"

	"github.com/google/uuid"
)

type StartGameParams struct {
	Broadcast chan []byte
	OwnerID   uuid.UUID
	Game      *models.Game
}

type LobbyHandler interface {
	StartNewGame(params StartGameParams) (*models.Game, error)
}

type lobbyHandler struct {
}

func newLobbyHandler() LobbyHandler {
	return &lobbyHandler{}
}

func (l lobbyHandler) StartNewGame(params StartGameParams) (*models.Game, error) {
	newGame := models.Game{}.NewGame()

	return &newGame, nil
}
