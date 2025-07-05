package math_operations

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mathly/internal/models"
	"mathly/internal/shared"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"

	"github.com/google/uuid"
)

func (m mathOperations) StartTheGame() {
	m.broadcastStartOfGame()
	m.broadcastScoreboard()
	m.broadcastQuestion()
}

func (m mathOperations) broadcastStartOfGame() {
	m.broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventStartOfGame,
		"",
	)
}

func (m mathOperations) broadcastGameEnd() {
	m.broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventEndOfGame,
		"",
	)
}

func (m mathOperations) messagePlayer(playerID uuid.UUID, message shared.SocketReponse) {
	p := m.findPlayerById(playerID)
	if p != nil {
		p.SendMessage(message)
	}
}

func (m mathOperations) broadcastScoreboard() {
	scoreboard := map[string]int{}
	for id, score := range m.scoreBoard {
		scoreboard[m.players[id].Nickname] = score
	}

	marshaledScoreboard, _ := json.Marshal(scoreboard)
	m.broadcast <- shared.CreateSocketResponse(
		shared.EventGame,
		shared.CommonGameEventScoreboard,
		string(marshaledScoreboard),
	)
}

func (m mathOperations) broadcastQuestion() {
	question, _ := json.Marshal(m.questions[0])
	marshaledQuestion, _ := json.Marshal(question)

	m.broadcast <- shared.CreateSocketResponse(
		shared.EventGame,
		math_operations_events.MathOperationsEventQuestion,
		string(marshaledQuestion),
	)
}

func (m mathOperations) sendGameEnd(playerID uuid.UUID) {
	m.messagePlayer(playerID, shared.CreateSocketResponse(
		shared.EventGame,
		shared.CommonGameEventFinishedGame,
		"",
	))
}

func (m mathOperations) generateAdditionQuestion() MathQuestion {
	numberOne := m.config.Random.Intn(1000) - 500
	numberTwo := m.config.Random.Intn(1000) - 500
	result := numberOne + numberTwo

	question := fmt.Sprintf(`What's the sum of %d + %d ?`, numberOne, numberTwo)

	answers := make([]string, 0)
	answers = append(answers, fmt.Sprintf("%d", result))
	for range 3 {
		answers = append(answers, fmt.Sprintf("%d", (result+m.config.Random.Intn(100)-50)))
	}
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	return MathQuestion{
		Question:      question,
		Answers:       answers,
		correctAnswer: fmt.Sprintf("%d", result),
	}
}

// func (m MathOperations) processGame() {
// 	for {
// 		select {
// 		case msg := <-m.config.MessageListener:
// 			msg.Data
// 		}
// 	}
// }

func (m mathOperations) findPlayerById(id uuid.UUID) *models.Player {
	for _, p := range m.players {
		if p.ConnectionID == id {
			return &p
		}
	}
	return nil
}

func (m mathOperations) HandleMessage(msg models.Message) {
	switch msg.Action {
	case models.ActionTypeGuessAnswer:
		m.handleAnswerMessage(msg)
	}
}

func (m mathOperations) handleAnswerMessage(msg models.Message) {
	p := m.findPlayerById(msg.SenderID)
	if p == nil {
		return // TODO
	}

	var data UserMessageData
	err := json.Unmarshal([]byte(msg.Data), &data)
	if err != nil {
		return // TODO
	}

	question := m.questions[m.playerQuestion[p.ConnectionID]]
	if question.correctAnswer == data.Answer {
		m.scoreBoard[p.ConnectionID]++
		m.broadcastScoreboard()
	} else {
		m.scoreBoard[p.ConnectionID]--
		m.broadcastScoreboard()
	}

	m.playerQuestion[p.ConnectionID]++

	// TODO add settings for that
	if m.playerQuestion[p.ConnectionID] >= 10 {
		m.sendGameEnd(p.ConnectionID)
		return
	}

	nextQuestion := m.questions[m.playerQuestion[p.ConnectionID]]
	p.SendMessage(shared.CreateSocketResponse(
		shared.EventGame,
		math_operations_events.MathOperationsEventQuestion,
		nextQuestion.String(),
	))
}

func (m mathOperations) GetRightAnswer(q *int) string {
	qNumber := 0
	if q != nil {
		qNumber = *q
	}
	return m.questions[qNumber].correctAnswer
}
