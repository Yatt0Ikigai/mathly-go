package math_operations

import (
	"encoding/json"
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

	marshaledScoreboard, _ := json.Marshal(scoreboard)

	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventEndOfGame,
		string(marshaledScoreboard),
	)
}

func (m mathOperations) broadcastScoreboard() {
	scoreboard := map[string]int{}
	for id, score := range m.scoreBoard {
		scoreboard[m.config.Players[id].Nickname] = score
	}

	marshaledScoreboard, _ := json.Marshal(scoreboard)
	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventGame,
		shared.CommonGameEventScoreboard,
		string(marshaledScoreboard),
	)
}

func (m mathOperations) broadcastQuestion() {
	question, _ := json.Marshal(m.questions[0])

	m.config.Broadcast <- shared.CreateSocketResponse(
		shared.EventGame,
		math_operations_events.MathOperationsEventQuestion,
		string(question),
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
