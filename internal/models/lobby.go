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

type LobbyPlayer struct {
	ConnectionID uuid.UUID `json:"connectionId"`
	Nickname     string    `json:"nickname"`
	AvatarUrl    string    `json:"avatarUrl"`
	Permission   int       `json:"permission"` // 0 - normal player, 1 - host // later change it
}
