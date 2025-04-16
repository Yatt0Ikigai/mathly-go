package sockets

import (
	"fmt"
	"maps"
	"mathly/internal/service"
	"slices"

	"github.com/google/uuid"
)

type LobbyManager interface {
	CreateLobby(services service.Service) Lobby
	FindLobby(uuid.UUID) Lobby
	ListLobbies() []uuid.UUID
}

type lobbyManager struct {
	Lobbies map[uuid.UUID]Lobby
}

func NewLobbyManager() LobbyManager {
	return &lobbyManager{
		Lobbies: make(map[uuid.UUID]Lobby),
	}
}

func (l lobbyManager) CreateLobby(services service.Service) Lobby {
	lobby := NewLobby(services)
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
