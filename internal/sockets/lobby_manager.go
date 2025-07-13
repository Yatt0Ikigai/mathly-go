package sockets

import (
	"fmt"
	"maps"
	"mathly/internal/service"
	"mathly/internal/sockets/games"
	"slices"

	"github.com/google/uuid"
)

type LobbyManager interface {
	CreateLobby() Lobby
	FindLobby(uuid.UUID) Lobby
	ListLobbies() []uuid.UUID
}

type lobbyManager struct {
	services    service.Service
	gameLibrary games.GameLibrary

	Lobbies map[uuid.UUID]Lobby
}

func NewLobbyManager(services service.Service, gameLibrary games.GameLibrary) LobbyManager {
	return &lobbyManager{
		services:    services,
		gameLibrary: gameLibrary,

		Lobbies: make(map[uuid.UUID]Lobby),
	}
}

func (l lobbyManager) CreateLobby() Lobby {
	lobby := NewLobby(l.services, l.gameLibrary)
	id := lobby.GetID()
	l.Lobbies[id] = lobby

	fmt.Println(&l)
	return lobby
}

func (l lobbyManager) FindLobby(id uuid.UUID) Lobby {
	return l.Lobbies[id]
}

func (l lobbyManager) ListLobbies() []uuid.UUID {
	return slices.Collect(maps.Keys(l.Lobbies))
}
