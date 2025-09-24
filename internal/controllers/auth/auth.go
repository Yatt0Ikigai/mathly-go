package auth

import (
	"mathly/internal/config"
	"mathly/internal/repository"
	"mathly/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	RegisterAuthEndpoints(s *gin.Engine)
}

type authController struct {
	oAuthController    OAuthController
	standardController StandardController
}

func NewAuthController(uR repository.User, jwtS service.JWT, c config.AuthOAuth) AuthController {
	oAuthController := NewOAuthController(uR, jwtS, c)
	standardController := NewStandard()

	return &authController{
		oAuthController:    oAuthController,
		standardController: standardController,
	}
}

func (a authController) RegisterAuthEndpoints(s *gin.Engine) {
	a.oAuthController.InitGoogleOAuth(s)
	a.standardController.InitStandard(s)
}
