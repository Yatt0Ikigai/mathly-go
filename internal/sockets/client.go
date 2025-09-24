//go:generate mockgen -source=client.go -package sockets -destination=client_mock.go
package sockets

import (
	"encoding/json"
	"mathly/internal/log"
	"mathly/internal/models"
	"mathly/internal/shared"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client interface {
	SendMessage(shared.SocketResponse)

	GetNickname() string
	GetID() uuid.UUID
	GetReceiver() chan []byte

	Close()
}

type client struct {
	ID       uuid.UUID
	Nickname string
	Socket   *websocket.Conn

	Receive chan []byte

	Lobby Lobby
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewClient(conn *websocket.Conn, l Lobby, nickname string) Client {
	c := &client{
		ID:       uuid.New(),
		Nickname: nickname,
		Socket:   conn,
		Receive:  make(chan []byte, Upgrader.ReadBufferSize),
		Lobby:    l,
	}

	l.JoinLobby(c)

	go c.Write()
	c.Read() // blocks until client disconnects

	// cleanup after read loop exits
	l.LeaveLobby(c)
	c.Socket.Close()

	return c
}

func (c *client) Read() {
	for {
		_, msg, err := c.Socket.ReadMessage()
		log.Log.Infof("Received message: ", string(msg))
		if err != nil {
			log.Log.Errorf("there was a error when reading message for Client: %+v", err)
			return // stop reading on error;
		}
		var data models.MessageDetails
		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Log.Infof("couldn't unmarshal message details %s", msg)
			continue
		}
		c.Lobby.ForwardMessage(models.Message{
			SenderID:       c.ID,
			MessageDetails: data,
		})
	}
}

func (c *client) Write() {
	for msg := range c.Receive {
		log.Log.Infof("Writing message: ", string(msg))
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Log.Errorf("write error for client %s: %v", c.Nickname, err)
			return
		}
	}
}

func (c *client) SendMessage(msg shared.SocketResponse) {
	c.Receive <- msg.ToByte()
}

func (c *client) Close() {
	close(c.Receive)
}

func (c *client) GetID() uuid.UUID {
	return c.ID
}

func (c *client) GetNickname() string {
	return c.Nickname
}

func (c *client) GetReceiver() chan []byte { return c.Receive }
