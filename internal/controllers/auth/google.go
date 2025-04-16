package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mathly/internal/config"
	"mathly/internal/models"
	"mathly/internal/repository"
	"mathly/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var AuthConfig oauth2.Config

type OAuthController interface {
	InitGoogleOAuth(s *gin.Engine)
}

type oAuthController struct {
	authConfig oauth2.Config

	jwtService     service.JWT
	userRepository repository.User
}

func NewOAuthController(uR repository.User, jwtS service.JWT, c config.AuthOAuth) OAuthController {
	authConfig := googleConfig(c)

	return &oAuthController{userRepository: uR, jwtService: jwtS, authConfig: authConfig}
}

func (o *oAuthController) InitGoogleOAuth(s *gin.Engine) {
	s.GET("/google_login", o.googleLogin)
	s.GET("/google_callback", o.googleCallback)
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

func (o *oAuthController) googleLogin(c *gin.Context) {
	url := AuthConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusMovedPermanently, url)
}

func (o *oAuthController) googleCallback(c *gin.Context) {
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusNotFound, "JSON Parsing Failed")
		return
	}

	var userData models.OAuthResponse
	json.Unmarshal(data, &userData)

	_, err = o.userRepository.Insert(&models.User{
		ID:        uuid.New(),
		Email:     userData.Email,
		Nickname:  userData.Name,
		Hash:      "123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		fmt.Printf("User creation failed due to %s", err.Error())
		c.String(http.StatusNotFound, "User creation failed")
		return
	}

	c.String(http.StatusOK, `ok`)
}
