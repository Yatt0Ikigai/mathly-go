package main

import (
	"mathly/internal/config"
	"mathly/internal/log"

	"net/http"
	"time"

	"mathly/internal/controllers/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.InitLogger()

	c := config.AppConfig
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

	auth.InitGoogleOAuth(r, c.OAuth)

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
