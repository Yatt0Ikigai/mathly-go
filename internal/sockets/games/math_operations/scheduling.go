package math_operations

import (
	"mathly/internal/log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func (m *mathOperations) addEventGameEvent() {
	task, _ := m.config.Scheduler.NewJob(
		// TODO: add config for that
		gocron.DurationJob(3*time.Minute),
		gocron.NewTask(func() {
			m.endGame()
		}),
		gocron.WithLimitedRuns(1),
	)

	m.events.EndOfGame = task.ID()
}

func (m *mathOperations) removePlayerTurnJob(pId uuid.UUID) {
	previousJob := m.events.PlayerTurns[pId]
	if previousJob == nil {
		return
	}

	_ = m.config.Scheduler.RemoveJob(*previousJob)
	m.events.PlayerTurns[pId] = nil
}

func (m *mathOperations) addPlayerTurnJob(pId uuid.UUID) {
	task, err := m.config.Scheduler.NewJob(
		// TODO: add config for that
		gocron.DurationJob(30*time.Second),
		gocron.NewTask(func() {
			m.handleAnswer(pId, false)
		}),
		gocron.WithLimitedRuns(1),
	)

	if err != nil {
		log.Log.Errorf("Error scheduling turn for player %s: %v", pId, err)
		return
	}

	taskId := task.ID()
	m.events.PlayerTurns[pId] = &taskId
}

func (m *mathOperations) prepareScheduler() {
	m.addEventGameEvent()
	for pId := range m.config.Players {
		m.addPlayerTurnJob(pId)
	}
}
