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

type SocketReponse struct {
	Event Event  `json:"event"`
	Type  string `json:"type"`
	Data  string `json:"data"`
}

func (s SocketReponse) Stringify() string {
	m, e := json.Marshal(s)
	if e != nil {
		log.Log.Errorf(`couldn't stringify %w`, s)
		return ""
	}
	return string(m)
}

func (s SocketReponse) ToByte() []byte {
	m, e := json.Marshal(s)
	if e != nil {
		log.Log.Errorf(`couldn't convert to byte %w`, s)
		return nil
	}
	return m
}

func CreateSocketResponse(e Event, eT eventType, d string) SocketReponse {
	return SocketReponse{
		Event: e,
		Type:  eT.String(),
		Data:  d,
	}
}
