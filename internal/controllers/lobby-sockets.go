package controllers

import (
	"mathly/internal/log"
	"mathly/internal/repository"
	"mathly/internal/service"
	"mathly/internal/sockets"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type lobbySocketsController struct {
	service        service.Service
	userRepository repository.User
	lobbyManager   sockets.LobbyManager
}

type LobbySocketsControllerParameters struct {
	Service        service.Service
	LobbyManager   sockets.LobbyManager
	UserRepository repository.User
}

// Dodaj pobieranie username z DB po context'cie
func NewLobbySockets(p LobbySocketsControllerParameters) *lobbySocketsController {
	return &lobbySocketsController{
		service:        p.Service,
		userRepository: p.UserRepository,
		lobbyManager:   p.LobbyManager,
	}
}

func (s lobbySocketsController) joinLobby(c *gin.Context) {
	var conn *websocket.Conn
	var parsedID uuid.UUID
	var err error

	tokenCookie, err := c.Cookie("access_token")
	if tokenCookie == "" || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := s.service.JWT().ValidateToken(tokenCookie, service.Access)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "expired token"})
		return
	}

	user, err := s.userRepository.GetByID(claims.UserID)
	if err != nil {
		log.Log.Errorln("There was an error during parsing %s context user ID to UUID", claims.UserID)
		return
	}

	if user == nil {
		if err = conn.WriteJSON(map[string]string{
			"message": "Couldn't retrieve session user",
		}); err != nil {
			log.Log.Errorln(err)
		}
		return
	}

	conn, err = sockets.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Log.Errorln(err)
		return
	}
	defer conn.Close()
	id := c.Param("id")

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
			sockets.NewClient(conn, l, user.Nickname)
			if err = conn.WriteJSON(map[string]string{
				"message": "User connected",
			}); err != nil {
				log.Log.Errorln(err)
			}
		}
	}
}

func (s lobbySocketsController) RegisterLobbyHandlers(router gin.IRouter) {
	router.GET("/:id/join", s.joinLobby)
}
