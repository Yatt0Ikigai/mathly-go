package math_operations

import (
	"encoding/json"
	"mathly/internal/models"
)

func (m mathOperations) StartTheGame() {
	m.prepareScheduler()

	m.broadcastStartOfGame()
	m.broadcastScoreboard()
	m.broadcastQuestion()

	m.config.Scheduler.Start()
}

func (m mathOperations) endGame() {
	m.config.EndGame()
	_ = m.config.Scheduler.Shutdown()
	m.broadcastGameEnd()
}

func (m *mathOperations) HandleMessage(msg models.Message) {
	switch msg.Action {
	case models.ActionTypeGuessAnswer:
		m.handleAnswerMessage(msg)
	}
}

func (m *mathOperations) handleAnswerMessage(msg models.Message) {
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
	correctAnswer := question.correctAnswer == data.Answer
	m.handleAnswer(msg.SenderID, correctAnswer)
}
