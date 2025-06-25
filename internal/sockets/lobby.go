package sockets

import (
	"fmt"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/sockets/games/math_operations"
	gameUtils "mathly/internal/sockets/games/utils"
	"mathly/internal/utils"

	"github.com/google/uuid"
)

type Lobby interface {
	GetID() uuid.UUID
	GetOwnerID() uuid.UUID
	GetClientBySocketID(socketID uuid.UUID) Client
	GetPlayersNicknamesWithout(string) []string
	GetPlayers() map[uuid.UUID]models.Player
	GetGame() gameUtils.Game

	JoinLobby(Client)
	LeaveLobby(Client)

	ForwardMessage(models.Message)
	BroadcastMessage([]byte)

	handleJoin(c Client)
	handleLeave(c Client)
	handleMessage(msg models.Message)
	handleLobbyMessage(msg models.Message)
}

type lobby struct {
	ID uuid.UUID

	Owner Client

	Join      chan Client
	Leave     chan Client
	Forward   chan models.Message
	Broadcast chan []byte

	Clients map[Client]bool

	Game         gameUtils.Game
	LobbyHandler service.LobbyHandler

	Settings Settings
}

type Settings struct{}

func NewLobby(services service.Service) Lobby {
	id := uuid.New()
	maxMessageAmount := 10

	l := lobby{
		ID: id,

		Forward:   make(chan models.Message, maxMessageAmount),
		Join:      make(chan Client),
		Leave:     make(chan Client),
		Clients:   make(map[Client]bool),
		Broadcast: make(chan []byte, maxMessageAmount),
	}

	go l.run()

	return &l
}

func (l *lobby) run() {
	for {
		select {
		case client := <-l.Join:
			l.handleJoin(client)
		case client := <-l.Leave:
			l.handleLeave(client)
		case msg := <-l.Forward:
			l.handleMessage(msg)
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

func (l *lobby) ForwardMessage(msg models.Message) {
	l.Forward <- msg
}

func (l *lobby) BroadcastMessage(msg []byte) {
	l.Broadcast <- msg
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

func (l *lobby) GetPlayers() map[uuid.UUID]models.Player {
	players := make(map[uuid.UUID]models.Player)

	for c := range l.Clients {
		players[c.GetID()] = models.Player{
			ConnectionID: c.GetID(),
			Nickname:     c.GetNickname(),
			SendMessage:  c.SendMessage,
		}
	}

	return players
}

func (l *lobby) GetGame() gameUtils.Game {
	return l.Game
}

func (l *lobby) GetClientBySocketID(socketID uuid.UUID) Client {
	for c := range l.Clients {
		if c.GetID() == socketID {
			return c
		}
	}
	return nil
}

func (l *lobby) handleJoin(c Client) {
	if l.Owner == nil {
		l.Owner = c
	}

	var playerJoinMessage, returnPlayerIDMessage []byte
	l.Broadcast <- fmt.Appendf(playerJoinMessage, "New Player %s Joined", c.GetNickname())
	c.SendMessage(fmt.Appendf(returnPlayerIDMessage, "%s", c.GetID().String()))
	l.Clients[c] = true
}

func (l *lobby) handleLeave(c Client) {
	delete(l.Clients, c)
	c.Close()
	var playerLeftMessage []byte
	l.Broadcast <- fmt.Appendf(playerLeftMessage, "Player %s Left", c.GetNickname())
}

func (l *lobby) handleMessage(msg models.Message) {
	if msg.Type == models.MessageTypeLobby {
		l.handleLobbyMessage(msg)
	}
	if msg.Type == models.MessageTypeGame && l.Game != nil {
		l.Game.HandleMessage(msg)
	}
	if msg.Type == models.MessageTypeLobby {
		fmt.Println("do sth")
	}
}

func (l *lobby) handleLobbyMessage(msg models.Message) {
	if msg.SenderID == l.GetOwnerID() {
		if msg.Action == models.ActionTypeStartGame {
			l.Game = math_operations.InitMathOperationsGame(gameUtils.GameConfig{
				Random:          utils.NewRandom(),
				MessageListener: l.Forward,
			}, l.GetPlayers(), l.Broadcast)
			l.Game.StartTheGame()
		}
	}
}
