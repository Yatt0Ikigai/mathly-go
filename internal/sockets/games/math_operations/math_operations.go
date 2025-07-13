package math_operations

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mathly/internal/models"
	"mathly/internal/shared"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func (m mathOperations) StartTheGame() {
	m.config.Scheduler.NewJob(
		gocron.DurationJob(5*time.Second),
		gocron.NewTask(func() {
			m.endGame()
		}),
		gocron.WithLimitedRuns(1),
	)

	m.broadcastStartOfGame()
	m.broadcastScoreboard()
	m.broadcastQuestion()
}

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

func (m mathOperations) endGame() {
	m.config.EndGame()
	m.broadcastGameEnd()
}

func (m mathOperations) messagePlayer(playerID uuid.UUID, message shared.SocketResponse) {
	p := m.findPlayerById(playerID)
	if p != nil {
		p.SendMessage(message)
	}
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

func (m mathOperations) sendGameEnd(playerID uuid.UUID) {
	m.messagePlayer(playerID, shared.CreateSocketResponse(
		shared.EventGame,
		shared.CommonGameEventFinishedGame,
		"",
	))
}

func (m mathOperations) generateAdditionQuestion() MathQuestion {
	random := m.config.Services.Random

	numberOne, _ := random.GenerateRandomNumber(1000)
	numberOne -= 500
	numberTwo, _ := random.GenerateRandomNumber(1000)
	numberTwo -= 500

	result := numberOne + numberTwo

	question := fmt.Sprintf(`What's the sum of %d + %d ?`, numberOne, numberTwo)

	answers := make([]string, 0)
	answers = append(answers, fmt.Sprintf("%d", result))
	for range 3 {
		diff, _ := random.GenerateRandomNumber(100)
		diff -= 50
		answers = append(answers, fmt.Sprintf("%d", (result+diff)))
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

func (m mathOperations) findPlayerById(id uuid.UUID) *models.Player {
	for _, p := range m.config.Players {
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

	if m.playerQuestion[p.ConnectionID] >= 10 {
		return
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
