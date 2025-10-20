package models

type CreateLobbyRequest struct {
	Name       string    `json:"lobbyName"`
	Type       LobbyType `json:"lobbyType"`
	MaxPlayers int8      `json:"maxPlayers"`
	Password   *string   `json:"password"`
}
