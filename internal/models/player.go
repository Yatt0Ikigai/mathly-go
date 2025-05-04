package models

import "github.com/google/uuid"

type Player struct {
	AccountID    *int64
	ConnectionID uuid.UUID
	Nickname     string
	Receiver     chan []byte
}

func (p Player) SendMessage(m string) {
	msg := []byte(m)
	p.Receiver <- msg
}
