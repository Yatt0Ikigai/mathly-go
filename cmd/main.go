package main

import (
	"mathly/internal/config"
	"mathly/internal/log"
	"mathly/internal/repository"
	"mathly/internal/service"
	"mathly/internal/sockets"
	"mathly/internal/sockets/games"

	"net/http"
	"time"

	"mathly/internal/controllers"
	"mathly/internal/controllers/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.InitLogger()

	c := config.AppConfig

	databases, err := repository.NewDatabases(&c.Databases)
	if err != nil {
		log.Log.Fatalf("Couldn't create databases - reason: %s", err.Error())
	}

	err = databases.DB().Health()
	if err != nil {
		return
	}

	repositories := repository.NewRepositories(databases)
	service := service.NewService(c.Services)

	gameLib := games.NewGameLibrary()

	lobbyManager := sockets.NewLobbyManager(service, gameLib)
	lobbySockets := controllers.NewLobbySockets(controllers.LobbySocketsControllerParameters{
		Service:        service,
		LobbyManager:   lobbyManager,
		UserRepository: repositories.User(),
	})

	lobbyRest := controllers.NewLobbyController(controllers.LobbyControllerParameters{
		Service:      service,
		Databases:    databases,
		LobbyManager: lobbyManager,
	})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowWebSockets:  true,
	}))

	r.Use(CORSMiddleware())

	oAuthController := auth.NewOAuthController(repositories.User(), service.JWT(), c.OAuth)
	oAuthController.InitGoogleOAuth(r)
	lobbySockets.RegisterLobbyHandlers(r)
	lobbyRest.RegisterLobbyRestHandlers(r)

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}
	log.Log.DPanicln(httpServer.ListenAndServe())
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow requests from localhost frontend
		if origin == "http://localhost:5173" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}

		// Handle preflight (OPTIONS) requests early
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
