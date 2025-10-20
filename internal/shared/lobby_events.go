package shared

type LobbyEvent string

var (
	LobbyEventStartOfGame LobbyEvent = "StartOfGame"
	LobbyEventEndOfGame   LobbyEvent = "EndOfGame"

	LobbyEventPlayerList LobbyEvent = "PlayerList"
	LobbyEventPlayerInfo LobbyEvent = "PlayerInfo"
)

func (l LobbyEvent) String() string {
	switch l {
	case LobbyEventStartOfGame:
		return "StartOfGame"
	case LobbyEventEndOfGame:
		return "EndOfGame"
	case LobbyEventPlayerList:
		return "PlayerList"
	case LobbyEventPlayerInfo:
		return "PlayerInfo"
	default:
		return "Unknown"
	}
}
