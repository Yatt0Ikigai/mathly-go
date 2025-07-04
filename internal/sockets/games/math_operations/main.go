package math_operations

import (
	"github.com/google/uuid"
	"mathly/internal/models"
	"mathly/internal/sockets/games/utils"
)

type MathOperations interface {
	utils.Game

	findPlayerById(id uuid.UUID) *models.Player

	broadcastStartOfGame()
	broadcastScoreboard()
	broadcastQuestion()
	broadcastGameEnd()

	sendGameEnd(playerID uuid.UUID)
	messagePlayer(playerID uuid.UUID, message []byte)

	generateAdditionQuestion() MathQuestion
	handleAnswerMessage(msg models.Message)
}

type mathOperations struct {
	config         utils.GameConfig
	broadcast      chan []byte
	questions      []MathQuestion
	players        map[uuid.UUID]models.Player
	scoreBoard     map[uuid.UUID]int
	playerQuestion map[uuid.UUID]int
}

func InitMathOperationsGame(c utils.GameConfig, players map[uuid.UUID]models.Player, broadcast chan []byte) MathOperations {
	m := mathOperations{
		config:         c,
		broadcast:      broadcast,
		players:        players,
		scoreBoard:     make(map[uuid.UUID]int),
		playerQuestion: make(map[uuid.UUID]int),
	}

	for _, player := range players {
		m.scoreBoard[player.ConnectionID] = 0
		m.playerQuestion[player.ConnectionID] = 0
	}

	m.questions = make([]MathQuestion, 0, 10)
	for range 10 {
		m.questions = append(m.questions, m.generateAdditionQuestion())
	}

	return m
}

func InitMockOperationsGame(c utils.GameConfig, players map[uuid.UUID]models.Player, broadcast chan []byte) mathOperations {
	m := mathOperations{
		config:         c,
		broadcast:      broadcast,
		players:        players,
		scoreBoard:     make(map[uuid.UUID]int),
		playerQuestion: make(map[uuid.UUID]int),
	}

	for _, player := range players {
		m.scoreBoard[player.ConnectionID] = 0
		m.playerQuestion[player.ConnectionID] = 0
	}

	return m
}
