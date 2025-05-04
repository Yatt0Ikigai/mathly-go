package models

import "github.com/google/uuid"

type MessageType string

const (
	MessageTypeLobby MessageType = "Lobby"
	MessageTypeGame  MessageType = "Game"
	MessageTypeChat  MessageType = "Chat"
)

type Message struct {
	SenderID uuid.UUID
	MessageDetails
}

type ActionType string

const (
	ActionTypeStartGame   ActionType = "StartGame"
	ActionTypeGuessAnswer ActionType = "GuessAnswer"
)

type MessageDetails struct {
	Type   MessageType `json:"type"`
	Action ActionType  `json:"action"`
	Data   string      `json:"data"`
}
