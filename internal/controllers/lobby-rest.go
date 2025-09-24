package controllers

import (
	"net/http"

	"mathly/internal/controllers/auth"
	"mathly/internal/models"
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
	var data models.CreateLobbyRequest
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, "Not valid data")
		return
	}

	lobby := s.lobbyManager.CreateLobby(data)

	c.JSON(http.StatusOK, gin.H{
		"lobbyID": lobby.GetID(),
	})
}

func (s *lobbyController) obtainLobbies(c *gin.Context) {
	lobbies := s.lobbyManager.ObtainLobbies()

	c.JSON(http.StatusOK, lobbies)
}

func (s *lobbyController) RegisterLobbyRestHandlers(router gin.IRouter) {
	protected := router.Group("/lobby")
	protected.Use(auth.AuthMiddleware(s.service.JWT()))
	protected.POST("", s.createLobby);
	protected.GET("", s.obtainLobbies);
}
