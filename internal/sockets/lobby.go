package sockets

import (
	"fmt"
	"mathly/internal/models"
	"mathly/internal/service"

	"github.com/google/uuid"
)

type Lobby interface {
	GetID() uuid.UUID
	GetOwnerID() uuid.UUID
	GetClientBySocketID(socketID uuid.UUID) Client
	GetPlayersNicknamesWithout(string) []string

	JoinLobby(Client)
	LeaveLobby(Client)

	ForwardMessage(Message)
}

type lobby struct {
	ID uuid.UUID

	Owner Client

	Join      chan Client
	Leave     chan Client
	Forward   chan Message
	Broadcast chan []byte

	Clients map[Client]bool

	Game         *models.Game
	LobbyHandler service.LobbyHandler

	Settings Settings
}

type Settings struct{}

type Message struct {
	SocketID uuid.UUID
	Data     []byte
}

func NewLobby(services service.Service) Lobby {
	id := uuid.New()
	maxMessageAmount := 10

	l := lobby{
		ID: id,

		Forward:   make(chan Message, maxMessageAmount),
		Join:      make(chan Client),
		Leave:     make(chan Client),
		Clients:   make(map[Client]bool),
		Broadcast: make(chan []byte, maxMessageAmount),

		LobbyHandler: services.LobbyHandler(),
	}

	go l.run()

	return &l
}

func (l *lobby) run() {
	for {
		select {
		case client := <-l.Join:
			l.Clients[client] = true
		case client := <-l.Leave:
			delete(l.Clients, client)
			client.Close()
		case msg := <-l.Forward:
			fmt.Println("Action From: ", msg.SocketID, " Data: ", msg.Data)
			l.Broadcast <- msg.Data
		case msg := <-l.Broadcast:
			for c := range l.Clients {
				c.SendMessage(msg)
			}
		}
	}
}

func (l *lobby) GetID() uuid.UUID {
	return l.ID
}

func (l *lobby) ForwardMessage(msg Message) {
	l.Forward <- msg
}

func (l *lobby) JoinLobby(c Client) {
	l.Join <- c

	if l.Owner == nil {
		l.Owner = c
	}
}

func (l *lobby) LeaveLobby(c Client) {
	l.Leave <- c
}

func (l *lobby) GetOwnerID() uuid.UUID {
	return l.Owner.GetID()
}

func (l *lobby) GetPlayersNicknamesWithout(nickname string) []string {
	var opponents []string

	for c := range l.Clients {
		n := c.GetNickname()
		if n != nickname {
			opponents = append(opponents, n)
		}
	}

	return opponents
}

func (l *lobby) GetClientBySocketID(socketID uuid.UUID) Client {
	for c := range l.Clients {
		if c.GetID() == socketID {
			return c
		}
	}
	return nil
}
