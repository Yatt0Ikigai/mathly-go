package sockets

import (
	"fmt"
	"mathly/internal/models"
	"mathly/internal/service"
	"mathly/internal/shared"
	"mathly/internal/sockets/games"
	common_games "mathly/internal/sockets/games/common"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Lobby interface {
	GetID() uuid.UUID
	GetSettings() models.LobbySettings
	GetOwnerID() uuid.UUID
	GetClientBySocketID(socketID uuid.UUID) Client
	GetPlayersNicknamesWithout(string) []string
	GetPlayers() map[uuid.UUID]models.Player
	GetGame() common_games.Game

	JoinLobby(Client)
	LeaveLobby(Client)

	ForwardMessage(models.Message)
	BroadcastMessage(shared.SocketResponse)

	handleJoin(c Client)
	handleLeave(c Client)
	handleMessage(msg models.Message)
	handleLobbyMessage(msg models.Message)
}

type LobbyServices struct {
	LobbyHandler service.LobbyHandler
	Random       service.Random
}

type lobby struct {
	ID       uuid.UUID
	Settings models.LobbySettings

	Owner Client

	Join      chan Client
	Leave     chan Client
	Forward   chan models.Message
	Broadcast chan shared.SocketResponse

	Clients map[Client]bool

	Game        common_games.Game
	Services    LobbyServices
	GameLibrary games.GameLibrary
}

func NewLobby(services service.Service, gameLib games.GameLibrary, settings models.LobbySettings) Lobby {
	id := uuid.New()
	maxMessageAmount := 10

	l := lobby{
		ID:          id,
		Settings:    settings,
		GameLibrary: gameLib,
		Services: LobbyServices{
			LobbyHandler: services.LobbyHandler(),
			Random:       services.Random(),
		},

		Forward:   make(chan models.Message, maxMessageAmount),
		Join:      make(chan Client),
		Leave:     make(chan Client),
		Clients:   make(map[Client]bool),
		Broadcast: make(chan shared.SocketResponse, maxMessageAmount),
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

func (l *lobby) BroadcastMessage(sR shared.SocketResponse) {
	l.Broadcast <- sR
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

func (l *lobby) GetGame() common_games.Game {
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

	l.Broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventPlayerJoined,
		c.GetNickname(),
	)
	l.Clients[c] = true

	c.SendMessage(shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventPlayerID,
		c.GetID().String(),
	))
}

func (l *lobby) handleLeave(c Client) {
	delete(l.Clients, c)
	c.Close()

	l.Broadcast <- shared.CreateSocketResponse(
		shared.EventLobby,
		shared.LobbyEventPlayerLeft,
		c.GetNickname(),
	)
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
		scheduler, _ := gocron.NewScheduler()

		if msg.Action == models.ActionTypeStartGame {
			l.Game = l.GameLibrary.StartNewGame(games.AvailableGamesMathOperations, common_games.GameConfig{
				Services: common_games.GameServices{
					Random: l.Services.Random,
				},
				Settings:  common_games.GameSettings{},
				Broadcast: l.Broadcast,
				Players:   l.GetPlayers(),
				Scheduler: scheduler,
				EndGame: func() {
					l.Game = nil
				},
			})

			l.Game.StartTheGame()
		}
	}
}

func (l *lobby) GetSettings() models.LobbySettings {
	return l.Settings
}
