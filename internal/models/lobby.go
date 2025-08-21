package models

import "github.com/google/uuid"

type LobbyType string

const (
	LobbyTypePublic  LobbyType = "Public"
	LobbyTypePrivate LobbyType = "Private"
)

type ListedLobby struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	AvailablePlayers int       `json:"availablePlayers"`
	MaxPlayers       int       `json:"maxPlayers"`
	LobbyType        LobbyType `json:"lobbyType"`
}

type LobbySettings struct {
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	MaxPlayers int       `json:"maxPlayers"`
	LobbyType  LobbyType `json:"lobbyType"`
}

func (l LobbySettings) DefaultSettings() LobbySettings {
	return LobbySettings{
		Name:       "Random Name",
		Password:   "",
		MaxPlayers: 8,
		LobbyType:  LobbyTypePublic,
	}
}
