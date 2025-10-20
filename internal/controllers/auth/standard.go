package auth

import "github.com/gin-gonic/gin"

type StandardController interface {
	InitStandard(s *gin.Engine)
}

type standardController struct{}

func NewStandard() StandardController {
	return &standardController{}
}

func (s standardController) InitStandard(r *gin.Engine) {
	r.POST("/logout", s.logout)
}

func (s standardController) logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "successfully logged out"})
}
