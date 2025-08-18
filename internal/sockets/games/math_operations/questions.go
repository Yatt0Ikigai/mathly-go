package math_operations

import (
	"fmt"
	"math/rand"
	"mathly/internal/shared"
	math_operations_events "mathly/internal/sockets/games/math_operations/events"

	"github.com/google/uuid"
)

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

func (m mathOperations) GetRightAnswer(q *int) string {
	qNumber := 0
	if q != nil {
		qNumber = *q
	}
	return m.questions[qNumber].correctAnswer
}

func (m *mathOperations) handleAnswer(pId uuid.UUID, correctAnswer bool) {
	p := m.config.Players[pId]
	if m.playerQuestion[pId] >= 10 {
		return
	}

	m.removePlayerTurnJob(pId)
	if correctAnswer {
		m.scoreBoard[pId]++
	} else {
		m.scoreBoard[pId]--
	}

	m.broadcastScoreboard()
	m.playerQuestion[pId]++

	if m.playerQuestion[pId] >= 10 {
		m.sendGameEnd(pId)
		m.stillPlaying -= 1

		if m.stillPlaying <= 0 {
			m.endGame()
		}

		return
	}

	nextQuestion := m.questions[m.playerQuestion[pId]]
	p.SendMessage(shared.CreateSocketResponse(
		shared.EventGame,
		math_operations_events.MathOperationsEventQuestion,
		nextQuestion.String(),
	))

	m.addPlayerTurnJob(pId)
}
