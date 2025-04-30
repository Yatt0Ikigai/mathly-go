package models

import "github.com/google/uuid"

type MessageType string

const (
	MessageTypeLobby MessageType = "Lobby"
	MessageTypeGame  MessageType = "Game"
)

type Message struct {
	SenderID uuid.UUID
	MessageDetails
}

type MessageDetails struct {
	Type MessageType
	Data string
}
