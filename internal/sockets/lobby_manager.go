package sockets

import (
	"fmt"
	"maps"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/sockets/games"
	"slices"

	"github.com/google/uuid"
)

type LobbyManager interface {
	CreateLobby(models.CreateLobbyRequest) Lobby
	ObtainLobbies() []models.ListedLobby
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

func (l lobbyManager) CreateLobby(data models.CreateLobbyRequest) Lobby {
	var password string
	if data.Type == models.LobbyTypePrivate {
		password = *data.Password
	}

	lobby := NewLobby(l.services, l.gameLibrary, models.LobbySettings{
		Name:       data.Name,
		MaxPlayers: int(data.MaxPlayers),
		LobbyType:  data.Type,
		Password:   password,
	})
	lobby.Start()
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

func (l lobbyManager) ObtainLobbies() []models.ListedLobby {
	var res []models.ListedLobby
	for lobbyId := range l.Lobbies {
		currentLobby := l.Lobbies[lobbyId]
		currentPlayers := currentLobby.GetPlayers()

		settings := currentLobby.GetSettings()

		res = append(res, models.ListedLobby{
			ID:               lobbyId,
			Name:             settings.Name,
			AvailablePlayers: len(currentPlayers),
			MaxPlayers:       settings.MaxPlayers,
			LobbyType:        settings.LobbyType,
		})
	}

	return res
}
