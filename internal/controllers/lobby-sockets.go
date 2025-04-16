package controllers

import (
	"mathly/internal/log"
	"mathly/internal/repository"
	"mathly/internal/service"
	"mathly/internal/sockets"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type lobbySocketsController struct {
	service      service.Service
	databases    repository.Databases
	lobbyManager sockets.LobbyManager
}

type LobbySocketsControllerParameters struct {
	Service      service.Service
	Databases    repository.Databases
	LobbyManager sockets.LobbyManager
}

func NewLobbySockets(p LobbySocketsControllerParameters) *lobbySocketsController {
	return &lobbySocketsController{
		service:      p.Service,
		databases:    p.Databases,
		lobbyManager: p.LobbyManager,
	}
}

func (s lobbySocketsController) joinLobby(c *gin.Context) {
	var conn *websocket.Conn
	var parsedID uuid.UUID
	var err error

	conn, err = sockets.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Log.Errorln(err)
		return
	}
	defer conn.Close()
	id := c.Param("id")
	nickname := c.Query("nickname")
	if nickname == "" {
		if err = conn.WriteJSON(map[string]string{
			"message": "User not passed nickname",
		}); err != nil {
			log.Log.Errorln(err)
		}
		return
	}

	parsedID, err = uuid.Parse(id)
	if err != nil {
		if err = conn.WriteJSON(map[string]string{
			"message": "User not passed lobby UUID",
		}); err != nil {
			log.Log.Errorln(err)
		}
	} else {
		l := s.lobbyManager.FindLobby(parsedID)
		if l == nil {
			if err = conn.WriteJSON(map[string]string{
				"message": "Lobby not found",
			}); err != nil {
				log.Log.Errorln(err)
			}
		} else {
			sockets.NewClient(conn, l, nickname)
			if err = conn.WriteJSON(map[string]string{
				"message": "User connected",
			}); err != nil {
				log.Log.Errorln(err)
			}
		}
	}
}

func (s lobbySocketsController) RegisterLobbyHandlers(router gin.IRouter) {
	router.GET("/join-lobby/:id", s.joinLobby)
}
