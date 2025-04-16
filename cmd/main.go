package main

import (
	"fmt"
	"mathly/internal/config"
	"mathly/internal/log"
	"mathly/internal/repository"
	"mathly/internal/service"

	"net/http"
	"time"

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
		fmt.Printf("ERROR AAAA: %v", err)
		return
	}
	fmt.Printf("TO JEST OK")
	_, err = databases.DB().Query("select * from information_schema.tables")
	if err != nil {
		fmt.Printf("ERROR AAAA: %v", err)
	}

	repositories := repository.NewRepositories(databases)
	service := service.NewService(c.Services)
	


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

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}
	log.Log.DPanicln(httpServer.ListenAndServe())
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:2137")

		c.Next()
	}
}
