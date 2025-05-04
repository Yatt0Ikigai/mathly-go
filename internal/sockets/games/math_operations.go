package games

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mathly/internal/models"

	gameUtils "mathly/internal/sockets/games/utils"
)

type LobbyEvent string

var (
	LobbyEventStartOfGame LobbyEvent = "StartOfGame"
	LobbyEventEndOfGame LobbyEvent = "EndOfGame"
)

type MathOperations struct {
	config           gameUtils.GameConfig
	forwardMessage   func(models.Message)
	broadcastMessage func([]byte)
	questions        []MathQuestion
	scoreBoard       map[models.Player]int
	playerQuestion   map[models.Player]int
}

func InitMathOperationsGame(c gameUtils.GameConfig, players []models.Player, fM func(models.Message), bM func([]byte)) MathOperations {
	m := MathOperations{
		config:           c,
		broadcastMessage: bM,
		forwardMessage:   fM,
		scoreBoard:       make(map[models.Player]int),
		playerQuestion:   make(map[models.Player]int),
	}

	for _, player := range players {
		m.scoreBoard[player] = 0
		m.playerQuestion[player] = 0
	}

	return m
}

type GameMessage struct {
	Type    string
	Message string
}

func (g GameMessage) ToByteArray() ([]byte, error) {
	return json.Marshal(g)
}

func (m MathOperations) StartTheGame() {
	for range 10 {
		m.questions = append(m.questions, m.generateAdditionQuestion())
	}
	m.broadcastStartOfGame()
	m.broadcastScoreboard()
	m.broadcastQuestion()
}

func (m *MathOperations) broadcastStartOfGame() {
	message := GameMessage{
		Type:    gameUtils.LobbyEventStartOfGame.String(),
		Message: "",
	}

	msg, _ := message.ToByteArray()
	m.broadcastMessage(msg)
}

func (m *MathOperations) broadcastScoreboard() {
	scoreboard := "🏆 Scoreboard:\n"
	for player, score := range m.scoreBoard {
		scoreboard += fmt.Sprintf("%s: %d\n", player.Nickname, score)
	}

	message := GameMessage{
		Type:    "Scoreboard",
		Message: scoreboard,
	}

	msg, _ := message.ToByteArray()
	m.broadcastMessage(msg)
}

func (m *MathOperations) broadcastQuestion() {
	question, _ := json.Marshal(m.questions[0])
	message := GameMessage{
		Type:    "Question",
		Message: string(question),
	}

	msg, _ := message.ToByteArray()
	m.broadcastMessage(msg)
}

type MathQuestion struct {
	Question      string
	Answers       []int
	correctAnswer int
}

func (m MathOperations) generateAdditionQuestion() MathQuestion {
	numberOne := m.config.Random.Intn(1000) - 500
	numberTwo := m.config.Random.Intn(1000) - 500
	result := numberOne + numberTwo

	question := fmt.Sprintf(`What's the sum of %d + %d ?`, numberOne, numberTwo)

	answers := make([]int, 0)
	answers = append(answers, result)
	for range 3 {
		answers = append(answers, result+m.config.Random.Intn(100)-50)
	}
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	return MathQuestion{
		Question:      question,
		Answers:       answers,
		correctAnswer: result,
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

func (m MathOperations) HandleMessage(msg models.Message) {

}
