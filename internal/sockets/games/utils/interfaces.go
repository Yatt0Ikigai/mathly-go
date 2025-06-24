package utils

import (
	"mathly/internal/models"
)

type Game interface {
	StartTheGame()
	HandleMessage(models.Message)
	GetRightAnswer(*int) string
}
