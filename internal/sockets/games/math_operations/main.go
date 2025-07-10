package math_operations

import (
	"mathly/internal/models"
	"mathly/internal/shared"
	common_games "mathly/internal/sockets/games/common"

	"github.com/google/uuid"
)

type MathOperations interface {
	common_games.Game

	findPlayerById(id uuid.UUID) *models.Player

	broadcastStartOfGame()
	broadcastScoreboard()
	broadcastQuestion()
	broadcastGameEnd()

	sendGameEnd(playerID uuid.UUID)
	messagePlayer(playerID uuid.UUID, message shared.SocketResponse)

	generateAdditionQuestion() MathQuestion
	handleAnswerMessage(msg models.Message)
}

type mathOperations struct {
	config         common_games.GameConfig
	questions      []MathQuestion
	scoreBoard     map[uuid.UUID]int
	playerQuestion map[uuid.UUID]int
}

func (m mathOperations) InitGame(c common_games.GameConfig) MathOperations {
	m.config = c
	m.scoreBoard = make(map[uuid.UUID]int)
	m.playerQuestion = make(map[uuid.UUID]int)

	for _, player := range c.Players {
		m.scoreBoard[player.ConnectionID] = 0
		m.playerQuestion[player.ConnectionID] = 0
	}

	m.questions = make([]MathQuestion, 0, 10)
	for range 10 {
		m.questions = append(m.questions, m.generateAdditionQuestion())
	}

	return m
}

func InitMathOperationsGame(c common_games.GameConfig) MathOperations {
	return mathOperations{}.InitGame(c)
}
