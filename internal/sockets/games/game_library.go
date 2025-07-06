//go:generate mockgen -source=game_library.go -package games -destination=game_library_mock.go
package games

import (
	common_games "mathly/internal/sockets/games/common"
	"mathly/internal/sockets/games/math_operations"
)

type AvailableGames string

const (
	AvailableGamesMathOperations AvailableGames = "MathOperations"
)

type GameLibrary interface {
	StartNewGame(game AvailableGames, c common_games.GameConfig) common_games.Game
}

type gameLibrary struct{}

func NewGameLibrary() GameLibrary {
	return &gameLibrary{}
}

func (g gameLibrary) StartNewGame(game AvailableGames, c common_games.GameConfig) common_games.Game {
	switch game {
	case AvailableGamesMathOperations:
		return math_operations.InitMathOperationsGame(c)
	default:
		return nil
	}
}
