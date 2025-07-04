package shared

import (
	"encoding/json"
	"mathly/internal/log"
)

type Event string

var (
	EventGame  Event = "Game"
	EventLobby Event = "Lobby"
)

type eventType interface {
	String() string
}

type socketReponse struct {
	Event Event  `json:"event"`
	Type  string `json:"type"`
	Data  string `json:"data"`
}

func (s socketReponse) Stringify() string {
	m, e := json.Marshal(s)
	if e != nil {
		log.Log.Errorf(`couldn't stringify %w`, s)
		return ""
	}
	return string(m)
}

func CreateSocketResponse(e Event, eT eventType, d string) socketReponse {
	return socketReponse{
		Event: e,
		Type:  eT.String(),
		Data:  d,
	}
}
