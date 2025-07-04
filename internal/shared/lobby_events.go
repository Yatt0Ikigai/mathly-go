package shared

type LobbyEvent string

var (
	LobbyEventStartOfGame LobbyEvent = "StartOfGame"
	LobbyEventEndOfGame   LobbyEvent = "EndOfGame"
	LobbyEventScoreboard  LobbyEvent = "Scoreboard"

	LobbyEventPlayerJoined LobbyEvent = "PlayerJoined"
	LobbyEventPlayerID     LobbyEvent = "PlayerID"
)

func (l LobbyEvent) String() string {
	switch l {
	case LobbyEventStartOfGame:
		return "StartOfGame"
	case LobbyEventEndOfGame:
		return "EndOfGame"
	case LobbyEventScoreboard:
		return "Scoreboard"
	case LobbyEventPlayerJoined:
		return "LobbyPlayerJoined"
	case LobbyEventPlayerID:
		return "PlayerID"
	default:
		return "Unknown"
	}
}
