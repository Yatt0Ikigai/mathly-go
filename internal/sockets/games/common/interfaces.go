//go:generate mockgen -source=interfaces.go -package common_games -destination=interfaces_mock.go

package common_games

import (
	"mathly/internal/models"
)

type Game interface {
	StartTheGame()
	HandleMessage(models.Message)
	GetRightAnswer(*int) string
}
