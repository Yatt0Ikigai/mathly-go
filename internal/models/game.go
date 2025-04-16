package models

import "github.com/google/uuid"

type Game struct {
	ID         uuid.UUID
	Players    []*Player
	Scoreboard map[*Player]int
}

func (g Game) NewGame() Game {
	return Game{
		ID:         uuid.New(),
		Players:    make([]*Player, 0),
		Scoreboard: make(map[*Player]int),
	}
}
