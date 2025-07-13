package math_operations

import (
	"mathly/internal/models"

	"github.com/google/uuid"
)

func (m mathOperations) findPlayerById(id uuid.UUID) *models.Player {
	for _, p := range m.config.Players {
		if p.ConnectionID == id {
			return &p
		}
	}
	return nil
}
