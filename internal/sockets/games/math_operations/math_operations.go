package math_operations

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"mathly/internal/models"
	"mathly/internal/sockets/games/utils"
)

func (m mathOperations) StartTheGame() {
	for range 10 {
		m.questions = append(m.questions, m.generateAdditionQuestion())
	}
	m.broadcastStartOfGame()
	m.broadcastScoreboard()
	m.broadcastQuestion()
}

func (m mathOperations) broadcastStartOfGame() {
	message := utils.GameMessage{
		Type:    utils.LobbyEventStartOfGame.String(),
		Message: "",
	}

	msg, _ := message.ToByteArray()
	m.broadcast <- msg
}

func (m mathOperations) broadcastScoreboard() {
	scoreboard := "🏆 Scoreboard:\n"
	for player, score := range m.scoreBoard {
		scoreboard += fmt.Sprintf("%s: %d\n", player.Nickname, score)
	}

	message := utils.GameMessage{
		Type:    utils.LobbyEventScoreboard.String(),
		Message: scoreboard,
	}

	msg, _ := message.ToByteArray()
	m.broadcast <- msg
}

func (m mathOperations) broadcastQuestion() {
	question, _ := json.Marshal(m.questions[0])
	message := utils.GameMessage{
		Type:    "Question",
		Message: string(question),
	}

	msg, _ := message.ToByteArray()
	m.broadcast <- msg
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
		m.handleAnswer(msg)
	}
}

func (m mathOperations) handleAnswer(msg models.Message) {
	p := m.findPlayerById(msg.SenderID)
	if p == nil {
		return // TODO
	}

	var data UserMessageData
	err := json.Unmarshal([]byte(msg.Data), &data)
	if err != nil {
		return // TODO
	}

	question := m.questions[m.playerQuestion[*p]]
	// TODO later check if it wasn't last question
	nextQuestion := m.questions[m.playerQuestion[*p]+1]

	if question.correctAnswer == data.Answer {
		m.scoreBoard[*p]++
		m.broadcastScoreboard()
	} else {
		m.scoreBoard[*p]--
		m.broadcastScoreboard()
	}

	m.playerQuestion[*p]++
	p.SendMessage(nextQuestion.String())
}
