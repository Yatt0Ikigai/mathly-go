package auth

import (
	"context"
	"io"
	"log"
	"mathly/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var AuthConfig oauth2.Config

func InitGoogleOAuth(s *gin.Engine, c config.AuthOAuth) {
	AuthConfig = googleConfig(c)

	s.POST("/google_login", googleLogin)
	s.GET("/google_callback", googleCallback)
}

func googleConfig(googleConfig config.AuthOAuth) oauth2.Config {
	AuthConfig = oauth2.Config{
		RedirectURL:  "http://localhost:8080/google_callback",
		ClientID:     googleConfig.ClientID,
		ClientSecret: googleConfig.ClientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return AuthConfig
}

func googleLogin(c *gin.Context) {
	url := AuthConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusMovedPermanently, url)
}

func googleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "randomstate" {
		c.String(http.StatusNotFound, "States don't Match!!")
		return
	}
	
	code := c.Query("code")
	
	token, err := AuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusNotFound, "Code-Token Exchange Failed")
		return
	}
	
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(http.StatusNotFound, "User Data Fetch Failed")
		return
	}
	
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusNotFound, "JSON Parsing Failed")
		return
	}
	
	log.Printf("User data: %+v", string(userData))
	c.String(http.StatusOK, string(userData))
}
