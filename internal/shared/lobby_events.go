package shared

type LobbyEvent string

var (
	LobbyEventStartOfGame LobbyEvent = "StartOfGame"
	LobbyEventEndOfGame   LobbyEvent = "EndOfGame"

	LobbyEventPlayerJoined LobbyEvent = "PlayerJoined"
	LobbyEventPlayerLeft   LobbyEvent = "PlayerLeft"
	LobbyEventPlayerID     LobbyEvent = "PlayerID"
)

func (l LobbyEvent) String() string {
	switch l {
	case LobbyEventStartOfGame:
		return "StartOfGame"
	case LobbyEventEndOfGame:
		return "EndOfGame"
	case LobbyEventPlayerJoined:
		return "LobbyPlayerJoined"
	case LobbyEventPlayerLeft:
		return "LobbyEventPlayerLeft"
	case LobbyEventPlayerID:
		return "PlayerID"
	default:
		return "Unknown"
	}
}
