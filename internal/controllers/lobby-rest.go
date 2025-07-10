package controllers

import (
	"net/http"

	"mathly/internal/repository"
	"mathly/internal/service"
	"mathly/internal/sockets"
	"mathly/internal/sockets/games"

	"github.com/gin-gonic/gin"
)

type lobbyController struct {
	gameLibrary  games.GameLibrary
	service      service.Service
	databases    repository.Databases
	lobbyManager sockets.LobbyManager
}

type LobbyControllerParameters struct {
	Service      service.Service
	Databases    repository.Databases
	LobbyManager sockets.LobbyManager
	GameLibrary  games.GameLibrary
}

func NewLobbyController(p LobbyControllerParameters) *lobbyController {
	return &lobbyController{
		gameLibrary:  p.GameLibrary,
		service:      p.Service,
		databases:    p.Databases,
		lobbyManager: p.LobbyManager,
	}
}

func (s *lobbyController) createLobby(c *gin.Context) {
	lobby := s.lobbyManager.CreateLobby(s.service, s.gameLibrary)

	c.JSON(http.StatusOK, gin.H{
		"lobbyID": lobby.GetID(),
	})
}

func (s *lobbyController) RegisterLobbyRestHandlers(router gin.IRouter) {
	router.POST("/create-lobby", s.createLobby)
}
