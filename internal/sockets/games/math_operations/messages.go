package math_operations

import (
	"mathly/internal/shared"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"

	"github.com/google/uuid"
)

func (m mathOperations) broadcastStartOfGame() {
	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventStartOfGame,
		"",
	)
}

func (m mathOperations) broadcastGameEnd() {
	scoreboard := map[string]int{}
	for id, score := range m.scoreBoard {
		scoreboard[m.config.Players[id].Nickname] = score
	}

	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventEndOfGame,
		scoreboard,
	)
}

func (m mathOperations) broadcastScoreboard() {
	scoreboard := map[string]int{}
	for id, score := range m.scoreBoard {
		scoreboard[m.config.Players[id].Nickname] = score
	}

	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventGame,
		shared.CommonGameEventScoreboard,
		scoreboard,
	)
}

func (m mathOperations) broadcastQuestion() {
	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventGame,
		math_operations_events.MathOperationsEventQuestion,
		m.questions[0],
	)
}

func (m mathOperations) messagePlayer(playerID uuid.UUID, message shared.SocketResponse) {
	p := m.findPlayerById(playerID)
	if p != nil {
		p.SendMessage(message)
	}
}

func (m mathOperations) sendGameEnd(playerID uuid.UUID) {
	m.messagePlayer(playerID, shared.CreateSocketResponse(
		shared.EventGame,
		shared.CommonGameEventFinishedGame,
		"",
	))
}
