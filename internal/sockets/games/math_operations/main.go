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

	generateAdditionQuestion() MathQuestion
	handleAnswer(msg models.Message)
}

type mathOperations struct {
	config         utils.GameConfig
	broadcast      chan []byte
	questions      []MathQuestion
	players        []models.Player
	scoreBoard     map[models.Player]int
	playerQuestion map[models.Player]int
}

func InitMathOperationsGame(c utils.GameConfig, players []models.Player, broadcast chan []byte) MathOperations {
	m := mathOperations{
		config:         c,
		broadcast:      broadcast,
		players:        players,
		scoreBoard:     make(map[models.Player]int),
		playerQuestion: make(map[models.Player]int),
	}

	for _, player := range players {
		m.scoreBoard[player] = 0
		m.playerQuestion[player] = 0
	}

	return m
}

func InitMockOperationsGame(c utils.GameConfig, players []models.Player, broadcast chan []byte) mathOperations {
	m := mathOperations{
		config:         c,
		broadcast:      broadcast,
		players:        players,
		scoreBoard:     make(map[models.Player]int),
		playerQuestion: make(map[models.Player]int),
	}

	for _, player := range players {
		m.scoreBoard[player] = 0
		m.playerQuestion[player] = 0
	}

	return m
}
